package handler

import (
	"log/slog"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"portfolio-backend/internal/model"
	"portfolio-backend/internal/repository"
)

func deriveThumbURL(imageURL string) string {
	idx := strings.Index(imageURL, "/static/")
	if idx < 0 {
		return ""
	}
	afterStatic := imageURL[idx+8:]
	if strings.HasPrefix(afterStatic, "thumbs/") {
		return ""
	}
	baseURL := imageURL[:idx+8]
	parts := strings.Split(afterStatic, "/")
	filename := parts[len(parts)-1]
	extIdx := strings.LastIndexByte(filename, '.')
	base := filename
	if extIdx >= 0 {
		base = filename[:extIdx]
	}
	return baseURL + "thumbs/" + base + "_thumb.jpg"
}

func (h *Handlers) HandleGetCertificates(c *fiber.Ctx) error {
	if h.DB == nil {
		return c.Status(503).JSON(fiber.Map{"error": "service unavailable"})
	}

	rows, err := repository.GetCertificates(h.DB, 100, 0)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}
	defer rows.Close()

	certs := []model.Certificate{}
	for rows.Next() {
		var cert model.Certificate
		var createdAt time.Time
		if err := rows.Scan(&cert.ID, &cert.Title, &cert.Issuer, &cert.Date, &cert.CredentialURL, &cert.ImageURL, &cert.Tags, &cert.IsVerified, &cert.DisplayOrder, &createdAt); err != nil {
			slog.Error("Failed to scan certificate row", "error", err)
			continue
		}
		cert.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
		cert.ThumbURL = deriveThumbURL(cert.ImageURL)
		certs = append(certs, cert)
	}

	return c.JSON(fiber.Map{"certificates": certs})
}

func (h *Handlers) HandleAdminGetCertificates(c *fiber.Ctx) error {
	if !h.checkAdminKey(c) {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	return h.HandleGetCertificates(c)
}

func (h *Handlers) HandleAdminCreateCertificate(c *fiber.Ctx) error {
	if !h.checkAdminKey(c) {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	if h.DB == nil {
		return c.Status(503).JSON(fiber.Map{"error": "service unavailable"})
	}

	var input struct {
		Title         string   `json:"title"`
		Issuer        string   `json:"issuer"`
		Date          string   `json:"date"`
		CredentialURL string   `json:"credential_url"`
		ImageURL      string   `json:"image_url"`
		Tags          []string `json:"tags"`
		IsVerified    bool     `json:"is_verified"`
		DisplayOrder  int      `json:"display_order"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	if strings.TrimSpace(input.Title) == "" {
		return c.Status(400).JSON(fiber.Map{"error": "title is required"})
	}
	if strings.TrimSpace(input.Issuer) == "" {
		return c.Status(400).JSON(fiber.Map{"error": "issuer is required"})
	}
	if input.Date != "" {
		if _, err := time.Parse("2006-01-02", input.Date); err != nil {
			if _, err2 := time.Parse("2006-01", input.Date); err2 != nil {
				return c.Status(400).JSON(fiber.Map{"error": "date must be YYYY-MM-DD or YYYY-MM"})
			}
		}
	}
	if input.ImageURL != "" && !strings.HasPrefix(input.ImageURL, "http://") && !strings.HasPrefix(input.ImageURL, "https://") && !strings.HasPrefix(input.ImageURL, "/") {
		return c.Status(400).JSON(fiber.Map{"error": "image_url must be a valid HTTP(S) or relative URL"})
	}
	if input.CredentialURL != "" && !strings.HasPrefix(input.CredentialURL, "http://") && !strings.HasPrefix(input.CredentialURL, "https://") {
		return c.Status(400).JSON(fiber.Map{"error": "credential_url must be a valid HTTP(S) URL"})
	}
	if input.Tags == nil {
		input.Tags = []string{}
	}

	id, err := repository.InsertCertificate(h.DB, input.Title, input.Issuer, input.Date, input.CredentialURL, input.ImageURL, input.Tags, input.IsVerified, input.DisplayOrder)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.Status(201).JSON(fiber.Map{"success": true, "id": id})
}

func (h *Handlers) HandleAdminUpdateCertificate(c *fiber.Ctx) error {
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
		Title         string   `json:"title"`
		Issuer        string   `json:"issuer"`
		Date          string   `json:"date"`
		CredentialURL string   `json:"credential_url"`
		ImageURL      string   `json:"image_url"`
		Tags          []string `json:"tags"`
		IsVerified    bool     `json:"is_verified"`
		DisplayOrder  int      `json:"display_order"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}
	if input.Tags == nil {
		input.Tags = []string{}
	}
	if input.Date != "" {
		if _, err := time.Parse("2006-01-02", input.Date); err != nil {
			if _, err2 := time.Parse("2006-01", input.Date); err2 != nil {
				return c.Status(400).JSON(fiber.Map{"error": "date must be YYYY-MM-DD or YYYY-MM"})
			}
		}
	}
	if input.ImageURL != "" && !strings.HasPrefix(input.ImageURL, "http://") && !strings.HasPrefix(input.ImageURL, "https://") && !strings.HasPrefix(input.ImageURL, "/") {
		return c.Status(400).JSON(fiber.Map{"error": "image_url must be a valid HTTP(S) or relative URL"})
	}
	if input.CredentialURL != "" && !strings.HasPrefix(input.CredentialURL, "http://") && !strings.HasPrefix(input.CredentialURL, "https://") {
		return c.Status(400).JSON(fiber.Map{"error": "credential_url must be a valid HTTP(S) URL"})
	}

	rowsAffected, err := repository.UpdateCertificate(h.DB, id, input.Title, input.Issuer, input.Date, input.CredentialURL, input.ImageURL, input.Tags, input.IsVerified, input.DisplayOrder)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "certificate not found"})
	}

	return c.JSON(fiber.Map{"success": true})
}

func (h *Handlers) HandleAdminDeleteCertificate(c *fiber.Ctx) error {
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

	rowsAffected, err := repository.DeleteCertificate(h.DB, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "certificate not found"})
	}

	return c.JSON(fiber.Map{"success": true})
}

func (h *Handlers) HandleAdminReorderCertificates(c *fiber.Ctx) error {
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

	if err := repository.ReorderCertificates(h.DB, input.Order); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.JSON(fiber.Map{"success": true})
}
