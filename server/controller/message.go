package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hunderaweke/sma-go/domain"
	"github.com/hunderaweke/sma-go/options"
)

type MessageController struct {
	messageUsecase domain.MessageUsecase
	roomUsecase    domain.RoomUsecase
	userUsecase    domain.UserUsecase
}

func NewMessageController(messageUC domain.MessageUsecase, roomUC domain.RoomUsecase, userUC domain.UserUsecase) *MessageController {
	return &MessageController{messageUsecase: messageUC, roomUsecase: roomUC, userUsecase: userUC}
}

type createMessage struct {
	FromUnique string `json:"from_unique"`
	Text       string `json:"text"`
}

func (mc *MessageController) CreateInRoom(c *fiber.Ctx) error {
	room, err := mc.roomUsecase.GetByUniqueString(c.Params("uniqueString"))
	if err != nil {
		return writeDomainError(c, err)
	}
	var req createMessage
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	msg, err := mc.messageUsecase.Create(domain.Message{
		RoomId:     room.ID,
		FromUnique: req.FromUnique,
		Text:       req.Text,
	})
	if err != nil {
		return writeDomainError(c, err)
	}
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

	res, err := mc.messageUsecase.GetByRoomUniqueString(room.UniqueString, opts)
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

func (mc *MessageController) authorizedRoom(c *fiber.Ctx) (*domain.Room, error) {
	room, err := mc.roomUsecase.GetByUniqueString(c.Params("uniqueString"))
	if err != nil {
		return nil, err
	}
	user, err := mc.currentUser(c)
	if err != nil {
		return nil, err
	}
	if room.OwnerID != user.ID {
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
