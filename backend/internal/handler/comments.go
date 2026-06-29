package handler

import (
	"log/slog"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"portfolio-backend/internal/model"
	"portfolio-backend/internal/repository"
)

func scanComment(scanner interface{ Scan(dest ...interface{}) error }) (model.BlogComment, error) {
	var c model.BlogComment
	var parentID *string
	var createdAt, updatedAt time.Time
	err := scanner.Scan(&c.ID, &c.PostID, &parentID, &c.UserID, &c.UserName, &c.AvatarURL, &c.Content, &createdAt, &updatedAt)
	if err == nil {
		c.ParentID = parentID
		c.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
		c.UpdatedAt = updatedAt.Format("2006-01-02 15:04:05")
	}
	return c, err
}

func nestComments(flat []model.BlogComment) []model.BlogComment {
	byID := make(map[string]*model.BlogComment, len(flat))

	for i := range flat {
		byID[flat[i].ID] = &flat[i]
	}

	// Process in reverse so deepest children are fully populated
	// before being copied into their parent's Replies slice.
	for i := len(flat) - 1; i >= 0; i-- {
		c := &flat[i]
		if c.ParentID != nil {
			if parent, ok := byID[*c.ParentID]; ok {
				parent.Replies = append(parent.Replies, *c)
			}
		}
	}

	var roots []model.BlogComment
	for i := range flat {
		if flat[i].ParentID == nil {
			roots = append(roots, flat[i])
		}
	}
	return roots
}

func (h *Handlers) HandleGetComments(c *fiber.Ctx) error {
	if h.DB == nil {
		return ErrUnavailable(c)
	}

	postID := c.Locals("post_id").(string)
	rows, err := repository.GetCommentsByPostID(h.DB, postID)
	if err != nil {
		return ErrInternal(c)
	}
	defer rows.Close()

	var flat []model.BlogComment
	for rows.Next() {
		cm, err := scanComment(rows)
		if err != nil {
			slog.Error("Failed to scan comment", "error", err)
			continue
		}
		flat = append(flat, cm)
	}

	nested := nestComments(flat)
	return c.JSON(fiber.Map{"comments": nested})
}

func (h *Handlers) HandleCreateComment(c *fiber.Ctx) error {
	if h.DB == nil {
		return ErrUnavailable(c)
	}

	postID := c.Locals("post_id").(string)
	userID, _ := c.Locals("user_id").(string)

	var input struct {
		Content  string  `json:"content"`
		ParentID *string `json:"parent_id"`
	}
	if err := c.BodyParser(&input); err != nil {
		return ErrValidation(c)
	}
	input.Content = strings.TrimSpace(input.Content)
	if input.Content == "" {
		return errResponse(c, 400, "ERR_VALIDATION", "content is required")
	}

	id, err := repository.CreateComment(h.DB, postID, userID, input.Content, input.ParentID)
	if err != nil {
		return ErrInternal(c)
	}

	row, err := repository.GetCommentByID(h.DB, id)
	if err != nil {
		return ErrInternal(c)
	}
	cm, err := scanComment(row)
	if err != nil {
		return ErrInternal(c)
	}

	return c.Status(201).JSON(fiber.Map{"comment": cm})
}

func (h *Handlers) HandleDeleteComment(c *fiber.Ctx) error {
	if h.DB == nil {
		return ErrUnavailable(c)
	}

	id := c.Params("commentID")
	userID, _ := c.Locals("user_id").(string)
	role, _ := c.Locals("user_role").(string)

	row, err := repository.GetCommentByID(h.DB, id)
	if err != nil {
		return ErrNotFound(c)
	}
	cm, err := scanComment(row)
	if err != nil || cm.ID == "" {
		return ErrNotFound(c)
	}

	if role != "admin" && cm.UserID != userID {
		return ErrUnauthorized(c)
	}

	if err := repository.DeleteCommentTree(h.DB, id); err != nil {
		return ErrInternal(c)
	}

	return c.JSON(fiber.Map{"success": true})
}
