package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/hunderaweke/sma-go/services"
)

func JWTMiddleware(c *fiber.Ctx) error {
	header := c.Get("Authorization", "")
	tokenString := ""
	if header != "" && strings.HasPrefix(header, "Bearer ") {
		tokenString = header[len("Bearer "):]
	} else if header != "" {
		tokenString = header
	}
	if tokenString == "" {
		tokenString = c.Cookies("access_token", "")
		if tokenString == "" {
			c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Authorization header"})
			return nil
		}
	}
	claims, err := services.ValidateToken(tokenString)
	if err != nil {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		return nil
	}
	c.Locals("user_id", claims.UserID)
	c.Locals("email", claims.Email)
	return c.Next()
}
