package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hunderaweke/sma-go/domain"
)

func writeDomainError(c *fiber.Ctx, err error) error {
	if err == nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if de, ok := err.(*domain.Error); ok {
		status := mapKindToStatus(de.Kind)
		return c.Status(status).JSON(fiber.Map{
			"error":  de.Msg,
			"kind":   de.Kind,
			"field":  de.Field,
			"entity": de.Entity,
		})
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
}

func mapKindToStatus(k domain.Kind) int {
	switch k {
	case domain.Required, domain.Invalid:
		return fiber.StatusBadRequest
	case domain.NotFound:
		return fiber.StatusNotFound
	case domain.AlreadyExists, domain.Conflict:
		return fiber.StatusConflict
	case domain.Unauthorized:
		return fiber.StatusUnauthorized
	case domain.Forbidden:
		return fiber.StatusForbidden
	default:
		return fiber.StatusInternalServerError
	}
}

func domainErrorFromValue(err domain.Error) *domain.Error {
	if err.Kind == "" && err.Msg == "" && err.Err == nil {
		return nil
	}
	return &err
}
