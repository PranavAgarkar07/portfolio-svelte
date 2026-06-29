package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"

	"portfolio-backend/internal/model"
	"portfolio-backend/internal/repository"
)

func (h *Handlers) deleteImageURLs(urls ...string) {
	if h.S3Client == nil {
		return
	}
	for _, imageURL := range urls {
		if imageURL == "" {
			continue
		}
		parsed, err := url.Parse(imageURL)
		if err != nil {
			continue
		}
		parts := strings.Split(parsed.Path, "/")
		if len(parts) < 2 {
			continue
		}
		key := strings.TrimPrefix(parsed.Path, "/")
		h.S3Client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
			Bucket: aws.String(model.S3Bucket),
			Key:    aws.String(key),
		})
	}
}

func scanPost(scanner interface {
	Scan(dest ...interface{}) error
}) (model.BlogPost, error) {
	var p model.BlogPost
	var publishedAt, createdAt, updatedAt time.Time
	var tags []string
	var imagesBytes []byte
	err := scanner.Scan(&p.ID, &p.Slug, &p.Title, &p.ContentMD, &p.Excerpt, &p.CoverImage, &imagesBytes, pq.Array(&tags), &p.Published, &publishedAt, &p.AuthorID, &p.AuthorName, &createdAt, &updatedAt)
	if err == nil {
		p.Tags = tags
		if len(imagesBytes) > 0 {
			if err := json.Unmarshal(imagesBytes, &p.Images); err != nil {
				slog.Error("Failed to unmarshal blog images", "error", err)
			}
		}
		if !publishedAt.IsZero() {
			p.PublishedAt = publishedAt.Format("2006-01-02 15:04:05")
		}
		p.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
		p.UpdatedAt = updatedAt.Format("2006-01-02 15:04:05")
	}
	return p, err
}

func (h *Handlers) HandleGetPublishedPosts(c *fiber.Ctx) error {
	if h.DB == nil {
		return ErrUnavailable(c)
	}

	rows, err := repository.GetPublishedPosts(h.DB)
	if err != nil {
		return ErrInternal(c)
	}
	defer rows.Close()

	posts := []model.BlogPost{}
	for rows.Next() {
		p, err := scanPost(rows)
		if err != nil {
			slog.Error("Failed to scan blog post", "error", err)
			continue
		}
		posts = append(posts, p)
	}
	return c.JSON(fiber.Map{"posts": posts})
}

func (h *Handlers) HandleGetPostBySlug(c *fiber.Ctx) error {
	if h.DB == nil {
		return ErrUnavailable(c)
	}

	slug := c.Params("slug")
	row, err := repository.GetPostBySlug(h.DB, slug)
	if err != nil {
		return ErrInternal(c)
	}

	p, err := scanPost(row)
	if err != nil {
		return ErrNotFound(c)
	}

	return c.JSON(fiber.Map{"post": p})
}

func (h *Handlers) HandleAdminGetPosts(c *fiber.Ctx) error {
	if h.DB == nil {
		return ErrUnavailable(c)
	}

	role, _ := c.Locals("user_role").(string)
	userID, _ := c.Locals("user_id").(string)

	var rows interface{ Close() error }
	var err error

	if role == "admin" {
		r, e := repository.GetAllPosts(h.DB)
		rows = r
		err = e
	} else {
		r, e := repository.GetPostsByAuthor(h.DB, userID)
		rows = r
		err = e
	}

	if err != nil {
		return ErrInternal(c)
	}

	posts := []model.BlogPost{}
	for rows.(interface {
		Next() bool
		Scan(dest ...interface{}) error
		Close() error
	}).Next() {
		p, err := scanPost(rows.(interface {
			Scan(dest ...interface{}) error
		}))
		if err != nil {
			slog.Error("Failed to scan blog post", "error", err)
			continue
		}
		posts = append(posts, p)
	}
	rows.Close()

	return c.JSON(fiber.Map{"posts": posts})
}

func (h *Handlers) HandleAdminGetPost(c *fiber.Ctx) error {
	if h.DB == nil {
		return ErrUnavailable(c)
	}

	id := c.Params("id")
	row, err := repository.GetPostByID(h.DB, id)
	if err != nil {
		return ErrInternal(c)
	}

	p, err := scanPost(row)
	if err != nil {
		return ErrNotFound(c)
	}

	return c.JSON(fiber.Map{"post": p})
}

func (h *Handlers) HandleAdminCreatePost(c *fiber.Ctx) error {
	if h.DB == nil {
		return ErrUnavailable(c)
	}

	userID, _ := c.Locals("user_id").(string)

	var input struct {
		Slug       string            `json:"slug"`
		Title      string            `json:"title"`
		ContentMD  string            `json:"content_md"`
		Excerpt    string            `json:"excerpt"`
		CoverImage string            `json:"cover_image"`
		Images     []model.BlogImage `json:"images"`
		Tags       []string          `json:"tags"`
		Published  bool              `json:"published"`
	}

	if err := c.BodyParser(&input); err != nil {
		return ErrValidation(c)
	}

	if strings.TrimSpace(input.Title) == "" {
		return errResponse(c, 400, "ERR_VALIDATION", "title is required")
	}
	if strings.TrimSpace(input.Slug) == "" {
		input.Slug = strings.ToLower(strings.ReplaceAll(strings.TrimSpace(input.Title), " ", "-"))
	}
	if strings.TrimSpace(input.ContentMD) == "" {
		return errResponse(c, 400, "ERR_VALIDATION", "content is required")
	}
	if input.Tags == nil {
		input.Tags = []string{}
	}
	if input.Images == nil {
		input.Images = []model.BlogImage{}
	}

	id, err := repository.CreatePost(h.DB, input.Slug, input.Title, input.ContentMD, input.Excerpt, input.CoverImage, input.Images, input.Tags, input.Published, userID)
	if err != nil {
		return ErrInternal(c)
	}

	return c.Status(201).JSON(fiber.Map{"success": true, "id": id})
}

func (h *Handlers) HandleAdminUpdatePost(c *fiber.Ctx) error {
	if h.DB == nil {
		return ErrUnavailable(c)
	}

	id := c.Params("id")
	userID, _ := c.Locals("user_id").(string)
	role, _ := c.Locals("user_role").(string)

	row, err := repository.GetPostByID(h.DB, id)
	if err != nil {
		return ErrNotFound(c)
	}
	var existing model.BlogPost
	existing, err = scanPost(row)
	if existing.ID == "" {
		return ErrNotFound(c)
	}
	if role != "admin" && existing.AuthorID != userID {
		return c.Status(403).JSON(fiber.Map{"error": "you can only edit your own posts"})
	}

	var input struct {
		Slug       string            `json:"slug"`
		Title      string            `json:"title"`
		ContentMD  string            `json:"content_md"`
		Excerpt    string            `json:"excerpt"`
		CoverImage string            `json:"cover_image"`
		Images     []model.BlogImage `json:"images"`
		Tags       []string          `json:"tags"`
		Published  bool              `json:"published"`
	}

	if err := c.BodyParser(&input); err != nil {
		return ErrValidation(c)
	}
	if input.Slug == "" {
		input.Slug = existing.Slug
	}
	if input.Title == "" {
		input.Title = existing.Title
	}
	if input.ContentMD == "" {
		input.ContentMD = existing.ContentMD
	}
	if input.Tags == nil {
		input.Tags = existing.Tags
	}
	if input.Images == nil {
		input.Images = existing.Images
	}

	orphanURLs := []string{}
	if existing.CoverImage != "" && existing.CoverImage != input.CoverImage {
		orphanURLs = append(orphanURLs, existing.CoverImage)
	}
	oldURLs := make(map[string]bool)
	for _, img := range existing.Images {
		oldURLs[img.URL] = true
	}
	for _, img := range input.Images {
		delete(oldURLs, img.URL)
	}
	for url := range oldURLs {
		orphanURLs = append(orphanURLs, url)
	}
	h.deleteImageURLs(orphanURLs...)

	if err := repository.UpdatePost(h.DB, id, input.Slug, input.Title, input.ContentMD, input.Excerpt, input.CoverImage, input.Images, input.Tags, input.Published); err != nil {
		return ErrInternal(c)
	}

	return c.JSON(fiber.Map{"success": true})
}

func (h *Handlers) HandleAdminDeleteImage(c *fiber.Ctx) error {
	var input struct {
		URL string `json:"url"`
	}
	if err := c.BodyParser(&input); err != nil {
		return ErrValidation(c)
	}
	if input.URL == "" {
		return errResponse(c, 400, "ERR_VALIDATION", "image url is required")
	}
	h.deleteImageURLs(input.URL)
	return c.JSON(fiber.Map{"success": true})
}

func (h *Handlers) HandleAdminDeletePost(c *fiber.Ctx) error {
	if h.DB == nil {
		return ErrUnavailable(c)
	}

	id := c.Params("id")
	userID, _ := c.Locals("user_id").(string)
	role, _ := c.Locals("user_role").(string)

	row, err := repository.GetPostByID(h.DB, id)
	if err != nil {
		return ErrNotFound(c)
	}
	existing, err := scanPost(row)
	if existing.ID == "" {
		return ErrNotFound(c)
	}
	if role != "admin" && existing.AuthorID != userID {
		return c.Status(403).JSON(fiber.Map{"error": "you can only delete your own posts"})
	}

	imageURLs := []string{}
	if existing.CoverImage != "" {
		imageURLs = append(imageURLs, existing.CoverImage)
	}
	for _, img := range existing.Images {
		imageURLs = append(imageURLs, img.URL)
	}
	h.deleteImageURLs(imageURLs...)

	if err := repository.DeletePost(h.DB, id); err != nil {
		return ErrInternal(c)
	}

	return c.JSON(fiber.Map{"success": true})
}
