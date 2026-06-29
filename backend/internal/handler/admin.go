package handler

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"

	"portfolio-backend/internal/repository"
)

func (h *Handlers) HandleAdminGetUsers(c *fiber.Ctx) error {
	if h.DB == nil {
		return ErrUnavailable(c)
	}

	rows, err := repository.GetAllUsers(h.DB)
	if err != nil {
		return ErrInternal(c)
	}
	defer rows.Close()

	type UserRow struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
		Role      string `json:"role"`
		CreatedAt string `json:"created_at"`
	}

	users := []UserRow{}
	for rows.Next() {
		var u UserRow
		var createdAt time.Time
		if err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.AvatarURL, &u.Role, &createdAt); err != nil {
			slog.Error("Failed to scan user", "error", err)
			continue
		}
		u.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
		users = append(users, u)
	}

	return c.JSON(fiber.Map{"users": users})
}

func (h *Handlers) HandleAdminUpdateUserRole(c *fiber.Ctx) error {
	if h.DB == nil {
		return ErrUnavailable(c)
	}

	userID := c.Params("id")

	var input struct {
		Role string `json:"role"`
	}

	if err := c.BodyParser(&input); err != nil {
		return ErrValidation(c)
	}

	if input.Role != "user" && input.Role != "author" && input.Role != "admin" {
		return errResponse(c, 400, "ERR_VALIDATION", "role must be user, author, or admin")
	}

	if err := repository.UpdateUserRole(h.DB, userID, input.Role); err != nil {
		return ErrInternal(c)
	}

	return c.JSON(fiber.Map{"success": true})
}

func (h *Handlers) HandleAdminGetReviews(c *fiber.Ctx) error {
	if h.DB == nil {
		return ErrUnavailable(c)
	}

	rows, err := repository.GetAllReviews(h.DB)
	if err != nil {
		return ErrInternal(c)
	}
	defer rows.Close()

	type ReviewRow struct {
		ID          int    `json:"id"`
		UserID      string `json:"user_id"`
		UserName    string `json:"user_name"`
		AvatarURL   string `json:"avatar_url"`
		ProjectName string `json:"project_name"`
		Rating      int    `json:"rating"`
		Comment     string `json:"comment"`
		CreatedAt   string `json:"created_at"`
	}

	reviews := []ReviewRow{}
	for rows.Next() {
		var r ReviewRow
		var createdAt time.Time
		if err := rows.Scan(&r.ID, &r.UserID, &r.UserName, &r.AvatarURL, &r.ProjectName, &r.Rating, &r.Comment, &createdAt); err != nil {
			slog.Error("Failed to scan review", "error", err)
			continue
		}
		r.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
		reviews = append(reviews, r)
	}

	return c.JSON(fiber.Map{"reviews": reviews})
}

func (h *Handlers) HandleAdminDeleteReview(c *fiber.Ctx) error {
	if h.DB == nil {
		return ErrUnavailable(c)
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return ErrValidation(c)
	}

	if err := repository.DeleteReviewByAdmin(h.DB, id); err != nil {
		return ErrInternal(c)
	}

	return c.JSON(fiber.Map{"success": true})
}
