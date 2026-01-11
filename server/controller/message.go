package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hunderaweke/sma-go/domain"
	"github.com/hunderaweke/sma-go/options"
)

type MessageController struct {
	usecase domain.MessageUsecase
}

func NewMessageController(uc domain.MessageUsecase) *MessageController {
	return &MessageController{usecase: uc}
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
