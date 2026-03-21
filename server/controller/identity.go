package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hunderaweke/sma-go/domain"
	"github.com/hunderaweke/sma-go/utils"
)

type IdentityController struct {
	usecase domain.IdentityUsecase
}

func NewIdentityController(uc domain.IdentityUsecase, handler *utils.PGPHandler) *IdentityController {
	if handler == nil {
		handler = utils.NewPGPHandler()
	}
	return &IdentityController{usecase: uc}
}

func (ic *IdentityController) Create(c *fiber.Ctx) error {
	identity, err := ic.usecase.Create()
	if err != nil {
		return writeDomainError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(identity)
}

func (ic *IdentityController) GetByUnique(c *fiber.Ctx) error {
	unique := c.Params("unique")
	identity, err := ic.usecase.GetByUniqueString(unique)
	if err != nil {
		return writeDomainError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(identity)
}
