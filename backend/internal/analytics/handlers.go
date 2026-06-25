package analytics

import (
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
)

func HandleCreateSession(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var s SessionPayload
		if err := c.BodyParser(&s); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid JSON body"})
		}
		if s.ID == "" || s.IPHash == "" {
			return c.Status(400).JSON(fiber.Map{"error": "id and ip_hash are required"})
		}
		if err := InsertSession(db, s); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to create session"})
		}
		return c.Status(201).JSON(fiber.Map{"status": "created"})
	}
}

func HandleCreateEvents(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var events []EventPayload
		if err := c.BodyParser(&events); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid JSON body"})
		}
		if len(events) > 50 {
			events = events[:50]
		}
		if len(events) > 0 {
			for _, e := range events {
				if e.SessionID == "" || e.Type == "" {
					return c.Status(400).JSON(fiber.Map{"error": "each event requires session_id and type"})
				}
			}
			if err := InsertEvents(db, events); err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "failed to store events"})
			}
		}
		return c.Status(202).JSON(fiber.Map{"status": "accepted"})
	}
}

func HandleDashboard(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sinceStr := c.Query("since", "")
		since := time.Now().Add(-30 * 24 * time.Hour)
		if sinceStr != "" {
			parsed, err := time.Parse(time.RFC3339, sinceStr)
			if err == nil {
				since = parsed
			}
		}
		stats := GetDashboardStats(db, since)
		return c.JSON(stats)
	}
}
