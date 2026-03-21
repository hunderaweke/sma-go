package controller

import (
	"bufio"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hunderaweke/sma-go/domain"
	"github.com/hunderaweke/sma-go/options"
	"github.com/valyala/fasthttp"
)

type MessageController struct {
	messageChannels sync.Map
	messageUsecase  domain.MessageUsecase
	roomUsecase     domain.RoomUsecase
	userUsecase     domain.UserUsecase
}

func NewMessageController(messageUC domain.MessageUsecase, roomUC domain.RoomUsecase, userUC domain.UserUsecase) *MessageController {
	return &MessageController{messageUsecase: messageUC, roomUsecase: roomUC, userUsecase: userUC, messageChannels: sync.Map{}}
}

type createMessage struct {
	FromUnique string `json:"from_unique"`
	Text       string `json:"text"`
}

func (mc *MessageController) addClient(roomUniqueString string) chan string {
	ch := make(chan string)
	mc.messageChannels.Store(roomUniqueString, ch)
	return ch
}

func (mc *MessageController) removeClient(roomUniqueString string) {
	mc.messageChannels.Delete(roomUniqueString)
}

func (mc *MessageController) sendToClient(roomUniqueString, msg string) {
	value, ok := mc.messageChannels.Load(roomUniqueString)
	if !ok {
		return
	}
	ch := value.(chan string)
	select {
	case ch <- msg:
	default:
	}
}

func (mc *MessageController) CreateInRoom(c *fiber.Ctx) error {
	room, err := mc.authorizedRoom(c)
	if err != nil {
		return writeDomainError(c, err)
	}

	var req createMessage
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	msg, err := mc.messageUsecase.Create(domain.Message{RoomId: room.ID, FromUnique: req.FromUnique, Text: req.Text})
	if err != nil {
		return writeDomainError(c, err)
	}
	mc.sendToClient(room.UniqueString, msg.Text)
	return c.Status(fiber.StatusCreated).JSON(msg)
}

func (mc *MessageController) ListInRoom(c *fiber.Ctx) error {
	room, err := mc.authorizedRoom(c)
	if err != nil {
		return writeDomainError(c, err)
	}

	var opts options.MessageFetchOptions
	if err := c.QueryParser(&opts); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	opts.RoomUniqueString = room.UniqueString

	res, err := mc.messageUsecase.GetAll(opts)
	if err != nil {
		return writeDomainError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

func (mc *MessageController) GetByIDInRoom(c *fiber.Ctx) error {
	room, err := mc.authorizedRoom(c)
	if err != nil {
		return writeDomainError(c, err)
	}
	msg, err := mc.messageUsecase.GetByID(c.Params("id"))
	if err != nil {
		return writeDomainError(c, err)
	}
	if msg.RoomId != room.ID {
		return writeDomainError(c, domain.New(domain.Forbidden, "you do not own this room"))
	}
	return c.Status(fiber.StatusOK).JSON(msg)
}

func (mc *MessageController) DeleteInRoom(c *fiber.Ctx) error {
	room, err := mc.authorizedRoom(c)
	if err != nil {
		return writeDomainError(c, err)
	}
	msg, err := mc.messageUsecase.GetByID(c.Params("id"))
	if err != nil {
		return writeDomainError(c, err)
	}
	if msg.RoomId != room.ID {
		return writeDomainError(c, domain.New(domain.Forbidden, "you do not own this room"))
	}
	if err := mc.messageUsecase.Delete(c.Params("id")); err != nil {
		return writeDomainError(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (mc *MessageController) ReceiveMessages(c *fiber.Ctx) error {
	room, err := mc.authorizedRoom(c)
	if err != nil {
		return writeDomainError(c, err)
	}

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Status(fiber.StatusOK).Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		ch := mc.addClient(room.UniqueString)
		defer mc.removeClient(room.UniqueString)
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case msg := <-ch:
				fmt.Fprintf(w, "data: %s\n\n", msg)
				if err := w.Flush(); err != nil {
					if err == fasthttp.ErrConnectionClosed {
						return
					}
					log.Printf("error while flushing: %v", err)
					return
				}
			case <-ticker.C:
				fmt.Fprint(w, "heartbeat\n\n")
				if err := w.Flush(); err != nil {
					if err == fasthttp.ErrConnectionClosed {
						return
					}
					log.Printf("error while flushing heartbeat: %v", err)
					return
				}
			}
		}
	}))
	return nil
}

func (mc *MessageController) authorizedRoom(c *fiber.Ctx) (*domain.Room, error) {
	room, err := mc.roomUsecase.GetByUniqueString(c.Params("uniqueString"))
	if err != nil {
		return nil, err
	}
	owner, err := mc.currentUser(c)
	if err != nil {
		return nil, err
	}
	if room.OwnerID != owner.ID {
		return nil, domain.New(domain.Forbidden, "you do not own this room")
	}
	return room, nil
}

func (mc *MessageController) currentUser(c *fiber.Ctx) (*domain.User, error) {
	userID, _ := c.Locals("user_id").(string)
	if userID == "" {
		return nil, domain.New(domain.Unauthorized, "missing authenticated user")
	}
	return mc.userUsecase.GetById(userID)
}
