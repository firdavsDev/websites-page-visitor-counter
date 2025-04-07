package handlers

import (
	"context"
	"visitor-counter/internal/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler struct {
	Store *storage.Storage
}

// @Summary Track page visit
// @Description Track a visitor by IP and page ID
// @Tags Visitor
// @Accept json
// @Produce json
// @Param page_id query string true "Page ID"
// @Success 200 {object} map[string]int
// @Failure 400 {object} map[string]string
// @Router /api/track [get]
func (h *Handler) TrackVisitor(c *fiber.Ctx) error {
	pageID := c.Query("page_id")
	if pageID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "page_id required"})
	}

	pid, err := uuid.Parse(pageID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid page_id"})
	}

	ip := c.IP()
	websiteID := c.Locals("website_id").(string)

	// Check if visitor already counted
	var exists bool
	err = h.Store.Conn.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM visitors WHERE ip=$1 AND page_id=$2)", ip, pid).Scan(&exists)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if !exists {
		// Insert visitor and increment count
		_, err = h.Store.Conn.Exec(context.Background(),
			"INSERT INTO visitors (ip, page_id) VALUES ($1, $2)", ip, pid)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		_, err = h.Store.Conn.Exec(context.Background(),
			"UPDATE pages SET visitor_count = visitor_count + 1 WHERE id=$1 AND website_id=$2", pid, websiteID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	var count int
	err = h.Store.Conn.QueryRow(context.Background(),
		"SELECT visitor_count FROM pages WHERE id=$1 AND website_id=$2", pid, websiteID).Scan(&count)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"visitor_count": count})
}
