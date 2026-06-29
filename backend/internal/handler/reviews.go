package handler

import (
	"log/slog"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"portfolio-backend/internal/model"
	"portfolio-backend/internal/repository"
)

func (h *Handlers) HandleGetProjectReviews(c *fiber.Ctx) error {
	if h.DB == nil {
		return ErrUnavailable(c)
	}

	projectName := c.Params("name")
	if projectName == "" {
		return ErrValidation(c)
	}

	rows, err := repository.GetProjectReviews(h.DB, projectName)
	if err != nil {
		return ErrInternal(c)
	}
	defer rows.Close()

	reviews := []model.ProjectReview{}
	for rows.Next() {
		var r model.ProjectReview
		var createdAt time.Time
		if err := rows.Scan(&r.UserID, &r.UserName, &r.AvatarURL, &r.ProjectName, &r.Rating, &r.Comment, &createdAt); err != nil {
			slog.Error("Failed to scan review", "error", err)
			continue
		}
		r.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
		reviews = append(reviews, r)
	}

	avg, count, _ := repository.GetProjectRatingSummary(h.DB, projectName)

	userReview := map[string]interface{}{}
	userID, ok := c.Locals("user_id").(string)
	if ok && userID != "" {
		row, err := repository.GetUserReview(h.DB, userID, projectName)
		if err == nil && row != nil {
			var r model.ProjectReview
			var createdAt time.Time
			if err := row.Scan(&r.UserID, &r.UserName, &r.AvatarURL, &r.ProjectName, &r.Rating, &r.Comment, &createdAt); err == nil {
				r.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
				userReview = map[string]interface{}{
					"rating":  r.Rating,
					"comment": r.Comment,
				}
			}
		}
	}

	return c.JSON(fiber.Map{
		"reviews":     reviews,
		"avg_rating":  avg,
		"count":       count,
		"user_review": userReview,
	})
}

func (h *Handlers) HandleSubmitReview(c *fiber.Ctx) error {
	if h.DB == nil {
		return ErrUnavailable(c)
	}

	projectName := c.Params("name")
	if projectName == "" {
		return ErrValidation(c)
	}

	userID, _ := c.Locals("user_id").(string)

	var input struct {
		Rating  int    `json:"rating"`
		Comment string `json:"comment"`
	}

	if err := c.BodyParser(&input); err != nil {
		return ErrValidation(c)
	}

	if input.Rating < 1 || input.Rating > 5 {
		return errResponse(c, 400, "ERR_VALIDATION", "rating must be between 1 and 5")
	}

	input.Comment = strings.TrimSpace(input.Comment)

	if err := repository.UpsertReview(h.DB, userID, projectName, input.Rating, input.Comment); err != nil {
		return ErrInternal(c)
	}

	avg, count, _ := repository.GetProjectRatingSummary(h.DB, projectName)

	return c.JSON(fiber.Map{
		"success":    true,
		"avg_rating": avg,
		"count":      count,
	})
}

func (h *Handlers) HandleDeleteReview(c *fiber.Ctx) error {
	if h.DB == nil {
		return ErrUnavailable(c)
	}

	projectName := c.Params("name")
	userID, _ := c.Locals("user_id").(string)

	if err := repository.DeleteReview(h.DB, userID, projectName); err != nil {
		return ErrInternal(c)
	}

	return c.JSON(fiber.Map{"success": true})
}

func (h *Handlers) HandleGetMarquee(c *fiber.Ctx) error {
	if h.DB == nil {
		return ErrUnavailable(c)
	}

	rows, err := repository.GetMarqueeItems(h.DB)
	if err != nil {
		return ErrInternal(c)
	}
	defer rows.Close()

	items := []model.MarqueeItem{}
	for rows.Next() {
		var m model.MarqueeItem
		var createdAt time.Time
		if err := rows.Scan(&m.Type, &m.UserName, &m.AvatarURL, &m.ProjectName, &m.Rating, &m.Comment, &createdAt); err != nil {
			slog.Error("Failed to scan marquee item", "error", err)
			continue
		}
		m.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
		items = append(items, m)
	}

	return c.JSON(fiber.Map{"items": items})
}
