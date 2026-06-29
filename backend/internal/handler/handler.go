package handler

import (
	"database/sql"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"

	"portfolio-backend/internal/config"
	"portfolio-backend/internal/middleware"
	"portfolio-backend/internal/model"
	"portfolio-backend/internal/repository"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func reqID(c *fiber.Ctx) slog.Attr {
	if id, ok := c.Locals("requestid").(string); ok {
		return slog.String("request_id", id)
	}
	return slog.String("request_id", "")
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func errResponse(c *fiber.Ctx, status int, code, message string) error {
	return c.Status(status).JSON(fiber.Map{"error": APIError{Code: code, Message: message}})
}

func (h *Handlers) checkAdminKey(c *fiber.Ctx) bool {
	auth := c.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") && strings.TrimPrefix(auth, "Bearer ") == h.Config.ContactSecret {
		return true
	}
	if c.Query("key") == h.Config.ContactSecret {
		return true
	}
	return false
}

var (
	ErrUnauthorized    = func(c *fiber.Ctx) error { return errResponse(c, 401, "ERR_UNAUTHORIZED", "unauthorized") }
	ErrUnavailable     = func(c *fiber.Ctx) error { return errResponse(c, 503, "ERR_UNAVAILABLE", "service unavailable") }
	ErrInternal        = func(c *fiber.Ctx) error { return errResponse(c, 500, "ERR_INTERNAL", "internal server error") }
	ErrNotFound        = func(c *fiber.Ctx) error { return errResponse(c, 404, "ERR_NOT_FOUND", "resource not found") }
	ErrValidation      = func(c *fiber.Ctx) error { return errResponse(c, 400, "ERR_VALIDATION", "invalid request") }
	ErrRateLimited     = func(c *fiber.Ctx) error { return errResponse(c, 429, "ERR_RATE_LIMIT", "rate limit exceeded") }
)

type Handlers struct {
	DB            *sql.DB
	S3Client      *s3.Client
	Config        *config.Config
	HTTPClient    *http.Client
	Cache         model.CachedData
	CacheMu       sync.RWMutex
	Metrics       model.Metrics
	Breaker       model.CircuitBreaker
	LastResponses []string
	ResponsesMu   sync.Mutex
	Auth          *AuthHandler
}

func New(cfg *config.Config, db *sql.DB, s3Client *s3.Client) *Handlers {
	return &Handlers{
		DB:         db,
		S3Client:   s3Client,
		Config:     cfg,
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        20,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		Metrics: model.Metrics{StartTime: time.Now()},
		Auth:    NewAuthHandler(cfg, db),
	}
}

func (h *Handlers) RegisterRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Sentinel API is Online 🟢")
	})
	app.Get("/healthz", h.HandleHealthz)

	app.Get("/api/status", h.HandleStatus)
	app.Get("/api/metrics", h.HandleMetrics)
	app.Post("/api/contact", h.HandleSubmitContact)
	app.Get("/admin/contact", h.HandleAdminContact)
	app.Get("/api/contact/messages", h.HandleGetMessages)
	app.Patch("/api/contact/messages/:id/read", h.HandleMarkRead)

	app.Get("/api/certificates", h.HandleGetCertificates)
	app.Get("/api/admin/certificates", h.HandleAdminGetCertificates)
	app.Post("/api/admin/certificates", h.HandleAdminCreateCertificate)
	app.Put("/api/admin/certificates/reorder", h.HandleAdminReorderCertificates)
	app.Put("/api/admin/certificates/:id", h.HandleAdminUpdateCertificate)
	app.Delete("/api/admin/certificates/:id", h.HandleAdminDeleteCertificate)
	app.Post("/api/admin/certificates/upload", h.HandleAdminUploadImage)

	app.Get("/api/badges", h.HandleGetBadges)
	app.Get("/api/admin/badges", h.HandleAdminGetBadges)
	app.Post("/api/admin/badges", h.HandleAdminCreateBadge)
	app.Put("/api/admin/badges/reorder", h.HandleAdminReorderBadges)
	app.Put("/api/admin/badges/:id", h.HandleAdminUpdateBadge)
	app.Delete("/api/admin/badges/:id", h.HandleAdminDeleteBadge)

	app.Get("/api/projects/likes", h.HandleGetProjectLikes)
	app.Post("/api/projects/like", h.likeRateLimit, h.HandleToggleProjectLike)

	auth := app.Group("/api/auth")
	auth.Get("/google-url", h.Auth.HandleGoogleURL)
	auth.All("/google/callback", h.Auth.HandleGoogleCallback)
	auth.Get("/me", middleware.AuthRequired(h.Config.JWTSecret), h.Auth.HandleAuthMe)

	authMw := middleware.AuthRequired(h.Config.JWTSecret)
	authorMw := middleware.RequireRole("author", "admin")
	adminMw := middleware.RequireRole("admin")

	blogPublic := app.Group("/api/blog")
	blogPublic.Get("/", h.HandleGetPublishedPosts)
	blogPublic.Get("/:slug", h.HandleGetPostBySlug)

	blogComments := app.Group("/api/blog/:slug/comments")
	blogComments.Use(func(c *fiber.Ctx) error {
		postID, err := repository.GetPostIDBySlug(h.DB, c.Params("slug"))
		if err != nil {
			return ErrNotFound(c)
		}
		c.Locals("post_id", postID)
		return c.Next()
	})
	blogComments.Get("/", h.HandleGetComments)
	blogComments.Post("/", authMw, h.HandleCreateComment)
	blogComments.Delete("/:commentID", authMw, h.HandleDeleteComment)

	blogAdmin := app.Group("/api/admin/blog", authMw, authorMw)
	blogAdmin.Get("/", h.HandleAdminGetPosts)
	blogAdmin.Get("/:id", h.HandleAdminGetPost)
	blogAdmin.Post("/", h.HandleAdminCreatePost)
	blogAdmin.Put("/:id", h.HandleAdminUpdatePost)
	blogAdmin.Delete("/:id", h.HandleAdminDeletePost)
	blogAdmin.Post("/upload-image", h.HandleBlogUploadImage)
	blogAdmin.Post("/delete-image", h.HandleAdminDeleteImage)

	reviews := app.Group("/api/projects/:name/reviews")
	reviews.Get("/", func(c *fiber.Ctx) error {
		if c.Get("Authorization") != "" {
			return authMw(c)
		}
		return c.Next()
	}, h.HandleGetProjectReviews)
	reviews.Post("/", authMw, h.HandleSubmitReview)
	reviews.Delete("/", authMw, h.HandleDeleteReview)

	app.Get("/api/feedback/marquee", h.HandleGetMarquee)

	admin := app.Group("/api/admin", authMw, adminMw)
	admin.Get("/users", h.HandleAdminGetUsers)
	admin.Put("/users/:id/role", h.HandleAdminUpdateUserRole)
	admin.Get("/reviews", h.HandleAdminGetReviews)
	admin.Delete("/reviews/:id", h.HandleAdminDeleteReview)
}
