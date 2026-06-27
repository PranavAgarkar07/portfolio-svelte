package handler

import (
	"log/slog"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"

	"portfolio-backend/internal/repository"
)

var (
	likeRateMu sync.Mutex
	likeRate   = map[string]int{}
)

func ResetLikeRate() {
	for range time.Tick(1 * time.Minute) {
		likeRateMu.Lock()
		likeRate = map[string]int{}
		likeRateMu.Unlock()
	}
}

func (h *Handlers) likeRateLimit(c *fiber.Ctx) error {
	ip := c.IP()
	likeRateMu.Lock()
	defer likeRateMu.Unlock()
	likeRate[ip]++
	if likeRate[ip] > 20 {
		return c.Status(429).JSON(fiber.Map{"error": "too many requests"})
	}
	return c.Next()
}

func (h *Handlers) HandleGetProjectLikes(c *fiber.Ctx) error {
	if h.DB == nil {
		return c.Status(503).JSON(fiber.Map{"error": "service unavailable"})
	}

	visitorToken := c.Query("visitor_token", "")

	rows, err := repository.GetProjectLikes(h.DB)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}
	defer rows.Close()

	likes := map[string]int{}
	for rows.Next() {
		var name string
		var cnt int
		if err := rows.Scan(&name, &cnt); err != nil {
			slog.Error("Failed to scan like row", "error", err)
			continue
		}
		likes[name] = cnt
	}

	userLikes := map[string]bool{}
	if visitorToken != "" {
		urows, err := repository.GetUserLikes(h.DB, visitorToken)
		if err != nil {
			slog.Error("Failed to query user likes", "error", err)
		} else {
			defer urows.Close()
			for urows.Next() {
				var name string
				if err := urows.Scan(&name); err == nil {
					userLikes[name] = true
				}
			}
		}
	}

	return c.JSON(fiber.Map{"likes": likes, "user_likes": userLikes})
}

func (h *Handlers) HandleToggleProjectLike(c *fiber.Ctx) error {
	if h.DB == nil {
		return c.Status(503).JSON(fiber.Map{"error": "service unavailable"})
	}

	var input struct {
		ProjectName  string `json:"project_name"`
		VisitorToken string `json:"visitor_token"`
		Liked        bool   `json:"liked"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid JSON body"})
	}
	if input.ProjectName == "" || input.VisitorToken == "" {
		return c.Status(400).JSON(fiber.Map{"error": "project_name and visitor_token are required"})
	}

	if err := repository.UpsertLike(h.DB, input.ProjectName, input.VisitorToken, input.Liked); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}

	count, err := repository.CountProjectLikes(h.DB, input.ProjectName)
	if err != nil {
		count = 0
	}

	return c.JSON(fiber.Map{"likes": count})
}
