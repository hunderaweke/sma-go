package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hunderaweke/sma-go/domain"
)

type AnalyticsController struct {
	usecase domain.AnalyticsUsecase
}

func NewAnalyticsController(uc domain.AnalyticsUsecase) *AnalyticsController {
	return &AnalyticsController{usecase: uc}
}

func (ac *AnalyticsController) Get(c *fiber.Ctx) error {
	res, err := ac.usecase.Get()
	if err != nil {
		return writeDomainError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(res)
}
