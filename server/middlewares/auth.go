package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hunderaweke/sma-go/services"
)

func JWTMiddleware(c *fiber.Ctx) error {
	header := c.Get("Authorization", "")
	if header == "" {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Authorization header"})
		return nil
	}
	tokenString := header[len("Bearer "):]
	claims, err := services.ValidateToken(tokenString)
	if err != nil {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		return nil
	}
	c.Locals("user_id", claims.UserID)
	c.Locals("email", claims.Email)
	return c.Next()
}
