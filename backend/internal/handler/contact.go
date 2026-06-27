package handler

import (
	"fmt"
	"html"
	"log/slog"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"portfolio-backend/internal/model"
	"portfolio-backend/internal/repository"
)

func (h *Handlers) HandleSubmitContact(c *fiber.Ctx) error {
	if h.DB == nil {
		return c.Status(503).JSON(fiber.Map{"error": "service unavailable"})
	}

	var sub model.ContactSubmission
	if err := c.BodyParser(&sub); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	if strings.TrimSpace(sub.Name) == "" {
		return c.Status(400).JSON(fiber.Map{"error": "name is required"})
	}

	if !emailRegex.MatchString(sub.Email) {
		return c.Status(400).JSON(fiber.Map{"error": "invalid email address"})
	}

	if len(strings.TrimSpace(sub.Message)) < 10 {
		return c.Status(400).JSON(fiber.Map{"error": "message must be at least 10 characters"})
	}

	if len(sub.Topic) > 200 {
		return c.Status(400).JSON(fiber.Map{"error": "topic too long (max 200 characters)"})
	}

	if err := repository.InsertContactMessage(h.DB, sub.Name, sub.Email, sub.Topic, sub.Message); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.Status(201).JSON(fiber.Map{"success": true, "message": "Message sent successfully."})
}

func (h *Handlers) HandleAdminContact(c *fiber.Ctx) error {
	if c.Query("key") != h.Config.ContactSecret {
		return c.Status(401).SendString("Unauthorized")
	}
	if h.DB == nil {
		return c.Status(503).JSON(fiber.Map{"error": "service unavailable"})
	}

	rows, err := repository.GetMessages(h.DB, 50, 0)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}
	defer rows.Close()

	type MessageRow struct {
		ID        int
		Name      string
		Email     string
		Topic     string
		Message   string
		IsRead    bool
		CreatedAt time.Time
	}

	var messages []MessageRow
	for rows.Next() {
		var m MessageRow
		if err := rows.Scan(&m.ID, &m.Name, &m.Email, &m.Topic, &m.Message, &m.IsRead, &m.CreatedAt); err != nil {
			slog.Error("Failed to scan contact message row", "error", err)
			continue
		}
		messages = append(messages, m)
	}

	var htmlBuilder strings.Builder
	htmlBuilder.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Contact Messages — Admin</title>
<style>
body { background: #1a1a2e; color: #e0e0e0; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; padding: 2rem; }
h1 { color: #e94560; }
table { width: 100%; border-collapse: collapse; margin-top: 1rem; }
th, td { padding: 0.75rem; text-align: left; border-bottom: 1px solid #333; }
th { background: #16213e; color: #e94560; font-weight: 600; }
tr:hover { background: #16213e; }
.badge-read { background: #0f3460; color: #e94560; padding: 0.25rem 0.5rem; border-radius: 4px; font-size: 0.8rem; }
.badge-new { background: #e94560; color: #fff; padding: 0.25rem 0.5rem; border-radius: 4px; font-size: 0.8rem; }
</style>
</head>
<body>
<h1>Contact Messages — Admin</h1>
<table>
<thead><tr><th>#</th><th>Name</th><th>Email</th><th>Topic</th><th>Message</th><th>Date</th><th>Status</th></tr></thead>
<tbody>`)

	for _, m := range messages {
		truncated := m.Message
		if len(truncated) > 100 {
			truncated = truncated[:100] + "..."
		}
		statusBadge := `<span class="badge-read">Read</span>`
		if !m.IsRead {
			statusBadge = `<span class="badge-new">New</span>`
		}
		htmlBuilder.WriteString(fmt.Sprintf(`<tr><td>%d</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>`,
			m.ID, html.EscapeString(m.Name), html.EscapeString(m.Email), html.EscapeString(m.Topic), html.EscapeString(truncated), m.CreatedAt.Format("Jan 2, 2006 15:04"), statusBadge))
	}

	htmlBuilder.WriteString(`</tbody>
</table>
</body>
</html>`)

	c.Type("text/html")
	return c.SendString(htmlBuilder.String())
}

func (h *Handlers) HandleGetMessages(c *fiber.Ctx) error {
	if c.Query("key") != h.Config.ContactSecret {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	if h.DB == nil {
		return c.Status(503).JSON(fiber.Map{"error": "service unavailable"})
	}

	rows, err := repository.GetMessages(h.DB, 50, 0)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}
	defer rows.Close()

	type MessageJSON struct {
		ID        int       `json:"id"`
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		Topic     string    `json:"topic"`
		Message   string    `json:"message"`
		IsRead    bool      `json:"is_read"`
		CreatedAt time.Time `json:"created_at"`
	}

	var messages []MessageJSON
	for rows.Next() {
		var m MessageJSON
		if err := rows.Scan(&m.ID, &m.Name, &m.Email, &m.Topic, &m.Message, &m.IsRead, &m.CreatedAt); err != nil {
			slog.Error("Failed to scan message row", "error", err)
			continue
		}
		messages = append(messages, m)
	}

	return c.JSON(fiber.Map{"messages": messages})
}

func (h *Handlers) HandleMarkRead(c *fiber.Ctx) error {
	if c.Query("key") != h.Config.ContactSecret {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	if h.DB == nil {
		return c.Status(503).JSON(fiber.Map{"error": "service unavailable"})
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}

	rowsAffected, err := repository.MarkMessageRead(h.DB, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "message not found"})
	}

	return c.JSON(fiber.Map{"success": true})
}
