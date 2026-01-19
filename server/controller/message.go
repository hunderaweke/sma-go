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
	usecase         domain.MessageUsecase
}

func NewMessageController(uc domain.MessageUsecase) *MessageController {
	return &MessageController{usecase: uc, messageChannels: sync.Map{}}
}

func (mc *MessageController) addClient(uniqueString string) chan string {
	ch := make(chan string)
	mc.messageChannels.Store(uniqueString, ch)
	return ch
}

func (mc *MessageController) removeClient(uniqueString string) {
	mc.messageChannels.Delete(uniqueString)
}

func (mc *MessageController) sendToClient(uniqueString, msg string) {
	value, ok := mc.messageChannels.Load(uniqueString)
	if !ok {
		return
	}
	ch := value.(chan string)
	select {
	case ch <- msg:
	default:
	}
}

func (mc *MessageController) Create(c *fiber.Ctx) error {
	var req struct {
		FromUnique string `json:"from_unique"`
		ToUnique   string `json:"to_unique"`
		Text       string `json:"text"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	msg, err := mc.usecase.Create(domain.Message{FromUnique: req.FromUnique, ToUnique: req.ToUnique, Text: req.Text})
	if err != nil {
		return writeDomainError(c, err)
	}
	mc.sendToClient(req.ToUnique, msg.Text)
	return c.Status(fiber.StatusCreated).JSON(msg)
}

func (mc *MessageController) List(c *fiber.Ctx) error {
	var opts options.MessageFetchOptions
	if err := c.QueryParser(&opts); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, derr := mc.usecase.GetAll(opts)
	if domainErr := domainErrorFromValue(derr); domainErr != nil {
		return writeDomainError(c, domainErr)
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

func (mc *MessageController) GetByReceiver(c *fiber.Ctx) error {
	receiver := c.Params("unique")
	res, derr := mc.usecase.GetByReceiverIdentity(receiver)
	if domainErr := domainErrorFromValue(derr); domainErr != nil {
		return writeDomainError(c, domainErr)
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

func (mc *MessageController) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	msg, err := mc.usecase.GetByID(id)
	if err != nil {
		return writeDomainError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(msg)
}

func (mc *MessageController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := mc.usecase.Delete(id); err != nil {
		return writeDomainError(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (mc *MessageController) ReceiveMessages(c *fiber.Ctx) error {
	uniqueString := c.Params("uniqueString", "")
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Status(fiber.StatusOK).Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		ch := mc.addClient(uniqueString)
		defer mc.removeClient(uniqueString)
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
