package handler

import (
	"log/slog"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"portfolio-backend/internal/model"
	"portfolio-backend/internal/repository"
)

func (h *Handlers) HandleGetBadges(c *fiber.Ctx) error {
	if h.DB == nil {
		return c.Status(503).JSON(fiber.Map{"error": "service unavailable"})
	}

	rows, err := repository.GetBadges(h.DB, 100, 0)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}
	defer rows.Close()

	badges := []model.Badge{}
	for rows.Next() {
		var b model.Badge
		var createdAt time.Time
		if err := rows.Scan(&b.ID, &b.Name, &b.ImageURL, &b.CredentialURL, &b.Rarity, &b.Category, &b.Important, &b.DisplayOrder, &createdAt); err != nil {
			slog.Error("Failed to scan badge row", "error", err)
			continue
		}
		b.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
		badges = append(badges, b)
	}

	return c.JSON(fiber.Map{"badges": badges})
}

func (h *Handlers) HandleAdminGetBadges(c *fiber.Ctx) error {
	if !h.checkAdminKey(c) {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	return h.HandleGetBadges(c)
}

func (h *Handlers) HandleAdminCreateBadge(c *fiber.Ctx) error {
	if !h.checkAdminKey(c) {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	if h.DB == nil {
		return c.Status(503).JSON(fiber.Map{"error": "service unavailable"})
	}

	var input struct {
		Name          string `json:"name"`
		ImageURL      string `json:"image_url"`
		CredentialURL string `json:"credential_url"`
		Rarity        string `json:"rarity"`
		Category      string `json:"category"`
		Important     bool   `json:"important"`
		DisplayOrder  int    `json:"display_order"`
	}

	if err := c.BodyParser(&input); err != nil {
		slog.Error("Badge create body parse error", "error", err)
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	if strings.TrimSpace(input.Name) == "" {
		return c.Status(400).JSON(fiber.Map{"error": "name is required"})
	}

	rarity := input.Rarity
	if rarity != "rare" && rarity != "uncommon" {
		rarity = "common"
	}

	if input.ImageURL != "" && !strings.HasPrefix(input.ImageURL, "http://") && !strings.HasPrefix(input.ImageURL, "https://") && !strings.HasPrefix(input.ImageURL, "/") {
		return c.Status(400).JSON(fiber.Map{"error": "image_url must be a valid HTTP(S) or relative URL"})
	}

	if input.DisplayOrder < 0 {
		input.DisplayOrder = 0
	}

	id, err := repository.InsertBadge(h.DB, input.Name, input.ImageURL, input.CredentialURL, rarity, input.Category, input.Important, input.DisplayOrder)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.Status(201).JSON(fiber.Map{"success": true, "id": id})
}

func (h *Handlers) HandleAdminUpdateBadge(c *fiber.Ctx) error {
	if !h.checkAdminKey(c) {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	if h.DB == nil {
		return c.Status(503).JSON(fiber.Map{"error": "service unavailable"})
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}

	var input struct {
		Name          string `json:"name"`
		ImageURL      string `json:"image_url"`
		CredentialURL string `json:"credential_url"`
		Rarity        string `json:"rarity"`
		Category      string `json:"category"`
		Important     bool   `json:"important"`
		DisplayOrder  int    `json:"display_order"`
	}

	if err := c.BodyParser(&input); err != nil {
		slog.Error("Badge update body parse error", "error", err)
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	rarity := input.Rarity
	if rarity != "rare" && rarity != "uncommon" {
		rarity = "common"
	}

	if input.DisplayOrder < 0 {
		input.DisplayOrder = 0
	}

	rowsAffected, err := repository.UpdateBadge(h.DB, id, input.Name, input.ImageURL, input.CredentialURL, rarity, input.Category, input.Important, input.DisplayOrder)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "badge not found"})
	}

	return c.JSON(fiber.Map{"success": true})
}

func (h *Handlers) HandleAdminDeleteBadge(c *fiber.Ctx) error {
	if !h.checkAdminKey(c) {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	if h.DB == nil {
		return c.Status(503).JSON(fiber.Map{"error": "service unavailable"})
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}

	rowsAffected, err := repository.DeleteBadge(h.DB, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "badge not found"})
	}

	return c.JSON(fiber.Map{"success": true})
}

func (h *Handlers) HandleAdminReorderBadges(c *fiber.Ctx) error {
	if !h.checkAdminKey(c) {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	if h.DB == nil {
		return c.Status(503).JSON(fiber.Map{"error": "service unavailable"})
	}

	var input struct {
		Order []int `json:"order"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	if err := repository.ReorderBadges(h.DB, input.Order); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.JSON(fiber.Map{"success": true})
}
