package middleware

import (
	"context"
	"visitor-counter/internal/storage"

	"github.com/gofiber/fiber/v2"
)

func Authenticate(st *storage.Storage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token required"})
		}

		var websiteID string
		err := st.Conn.QueryRow(context.Background(),
			"SELECT id FROM websites WHERE token = $1", token).Scan(&websiteID)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		c.Locals("website_id", websiteID)
		return c.Next()
	}
}
