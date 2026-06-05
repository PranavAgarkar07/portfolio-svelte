package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
)

// --- Configuration ---
const (
	GitHubUsername = "PranavAgarkar07"
	CacheDuration  = 4 * time.Hour
	RequestTimeout = 15 * time.Second

	QuotaFallback = "Deep in the code — pushing updates across multiple projects. Check back soon!"

	OpenRouterModel = "openrouter/free"

	maxFailures    = 3
	breakerTimeout = 15 * time.Minute
)

// --- Structs ---

type DevLogResponse struct {
	Summary    string `json:"summary"`
	LastUpdate string `json:"last_update"`
	Source     string `json:"source"`
}

type CachedData struct {
	Response  DevLogResponse
	ExpiresAt time.Time
}

type CircuitBreaker struct {
	failures    int
	lastFailure time.Time
	mu          sync.Mutex
}

type Metrics struct {
	mu           sync.Mutex
	cacheHits    int64
	cacheMisses  int64
	llmErrors    int64
	totalLatency time.Duration
	totalCalls   int64
	startTime    time.Time
}

	type ContactSubmission struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Topic   string `json:"topic"`
		Message string `json:"message"`
	}

type Certificate struct {
	ID            int      `json:"id"`
	Title         string   `json:"title"`
	Issuer        string   `json:"issuer"`
	Date          string   `json:"date"`
	CredentialURL string   `json:"credential_url"`
	ImageURL      string   `json:"image_url"`
	Tags          []string `json:"tags"`
	IsVerified    bool     `json:"is_verified"`
	DisplayOrder  int      `json:"display_order"`
	CreatedAt     string   `json:"created_at"`
}

type Badge struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	ImageURL      string `json:"image_url"`
	CredentialURL string `json:"credential_url"`
	Rarity        string `json:"rarity"`
	Category      string `json:"category"`
	Important     bool   `json:"important"`
	DisplayOrder  int    `json:"display_order"`
	CreatedAt     string `json:"created_at"`
}

type MemoryCert struct {
	Certificate
}

var (
	cache         CachedData
	cacheMutex    sync.Mutex
	breaker       CircuitBreaker
	lastResponses []string
	responsesMu   sync.Mutex
	metrics       Metrics
	db            *sql.DB

	memCerts    []MemoryCert
	memCertsMu  sync.Mutex
	memCertID   int
)

func init() {
	metrics.startTime = time.Now()
}

// --- Main ---

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, relying on system env vars")
	} else {
		slog.Info(".env file loaded successfully")
	}

	key := os.Getenv("OPENROUTER_API_KEY")
	if len(key) > 10 {
		slog.Info("OpenRouter API key loaded", "prefix", key[:4], "suffix", key[len(key)-4:])
	} else {
		slog.Warn("OpenRouter API key is empty or invalid")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL != "" {
		var err error
		db, err = sql.Open("postgres", databaseURL)
		if err != nil {
			slog.Warn("Failed to open database", "error", err)
		} else {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := db.PingContext(ctx); err != nil {
				slog.Warn("Failed to ping database", "error", err)
				db.Close()
				db = nil
			} else {
				slog.Info("Database connected successfully")
				migrateQuery := `CREATE TABLE IF NOT EXISTS contact_messages (
					id SERIAL PRIMARY KEY,
					name TEXT NOT NULL,
					email TEXT NOT NULL,
					topic TEXT DEFAULT '',
					message TEXT NOT NULL,
					is_read BOOLEAN DEFAULT FALSE,
					created_at TIMESTAMPTZ DEFAULT NOW()
				)`
			db.Exec(`ALTER TABLE contact_messages ADD COLUMN IF NOT EXISTS topic TEXT DEFAULT ''`)
				if _, err := db.Exec(migrateQuery); err != nil {
					slog.Warn("Auto-migrate failed", "error", err)
				}

			certMigrate := `CREATE TABLE IF NOT EXISTS certificates (
				id SERIAL PRIMARY KEY,
				title TEXT NOT NULL,
				issuer TEXT NOT NULL,
				date TEXT DEFAULT '',
				credential_url TEXT DEFAULT '',
				image_url TEXT DEFAULT '',
				tags TEXT[] DEFAULT '{}',
				is_verified BOOLEAN DEFAULT FALSE,
				display_order INT DEFAULT 0,
				created_at TIMESTAMPTZ DEFAULT NOW()
			)`
			if _, err := db.Exec(certMigrate); err != nil {
				slog.Warn("certificates auto-migrate failed", "error", err)
			}

			badgeMigrate := `CREATE TABLE IF NOT EXISTS badges (
				id SERIAL PRIMARY KEY,
				name TEXT NOT NULL,
				image_url TEXT DEFAULT '',
				credential_url TEXT DEFAULT '',
				rarity TEXT DEFAULT 'common',
				category TEXT DEFAULT '',
				important BOOLEAN DEFAULT false,
				display_order INT DEFAULT 0,
				created_at TIMESTAMPTZ DEFAULT NOW()
			)`
			if _, err := db.Exec(badgeMigrate); err != nil {
				slog.Warn("badges auto-migrate failed", "error", err)
			}
			db.Exec("ALTER TABLE badges ADD COLUMN IF NOT EXISTS category TEXT DEFAULT ''")
			db.Exec("ALTER TABLE badges ADD COLUMN IF NOT EXISTS important BOOLEAN DEFAULT false")

			db.Exec(`UPDATE badges SET image_url = REPLACE(image_url, 'https://sentinel-backend-4x3i.onrender.com/static/uploads/', '/badges/') WHERE image_url LIKE '%/static/uploads/%'`)
			db.Exec(`UPDATE badges SET image_url = REPLACE(image_url, 'http://localhost:8080/static/uploads/', '/badges/') WHERE image_url LIKE '%/static/uploads/%'`)
			}
		}
	}
	if db != nil {
		defer db.Close()
	}

	app := fiber.New()

	app.Use(requestid.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://pranavagarkar07.github.io, http://localhost:5173, http://localhost:5174, http://localhost:4173",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Use(limiter.New(limiter.Config{
		Max:        120,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(fiber.Map{"error": "rate limit exceeded"})
		},
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Sentinel API is Online 🟢")
	})

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "time": time.Now().Unix()})
	})

	app.Get("/api/status", handleStatus)
	app.Get("/api/metrics", handleMetrics)
	app.Post("/api/contact", handleSubmitContact)
	app.Get("/admin/contact", handleAdminContact)
	app.Get("/api/contact/messages", handleGetMessages)
	app.Patch("/api/contact/messages/:id/read", handleMarkRead)

	app.Static("/static", "./static")
	app.Static("/badges", "./static/badges")

	app.Get("/api/certificates", handleGetCertificates)
	app.Get("/api/admin/certificates", handleAdminGetCertificates)
	app.Post("/api/admin/certificates", handleAdminCreateCertificate)
	app.Put("/api/admin/certificates/:id", handleAdminUpdateCertificate)
	app.Delete("/api/admin/certificates/:id", handleAdminDeleteCertificate)
	app.Put("/api/admin/certificates/reorder", handleAdminReorderCertificates)
	app.Post("/api/admin/certificates/upload", handleAdminUploadImage)

	app.Get("/api/badges", handleGetBadges)
	app.Get("/api/admin/badges", handleAdminGetBadges)
	app.Post("/api/admin/badges", handleAdminCreateBadge)
	app.Put("/api/admin/badges/:id", handleAdminUpdateBadge)
	app.Delete("/api/admin/badges/:id", handleAdminDeleteBadge)
	app.Put("/api/admin/badges/reorder", handleAdminReorderBadges)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	slog.Info("Sentinel starting", "port", port)
	if err := app.Listen(":" + port); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}

// --- Handlers ---

func handleStatus(c *fiber.Ctx) error {
	cacheMutex.Lock()
	valid := time.Now().Before(cache.ExpiresAt)
	hasData := cache.Response.Summary != ""
	expired := !valid && hasData
	resp := cache.Response
	cacheMutex.Unlock()

	metrics.mu.Lock()
	metrics.totalCalls++
	if valid {
		metrics.cacheHits++
	} else if expired {
		metrics.cacheHits++
	} else {
		metrics.cacheMisses++
	}
	metrics.mu.Unlock()

	if valid {
		slog.Info("Serving from cache", "source", "cache", "request_id", c.Locals("requestid"))
		resp.Source = "cache"
		return c.JSON(resp)
	}

	if expired {
		slog.Info("Serving stale cache, refreshing in background", "request_id", c.Locals("requestid"), "source", "stale-cache")
		resp.Source = "stale-cache"
		go func() {
			start := time.Now()
			summary, err := generateDevLog()
			if err != nil {
				slog.Error("Background refresh failed", "error", err)
				metrics.mu.Lock()
				metrics.llmErrors++
				metrics.mu.Unlock()
				return
			}
			cacheMutex.Lock()
			cache.Response = DevLogResponse{
				Summary:    summary,
				LastUpdate: time.Now().Format("2006-01-02 15:04:05"),
				Source:     "live",
			}
			cache.ExpiresAt = time.Now().Add(CacheDuration)
			cacheMutex.Unlock()
			slog.Info("Cache refreshed in background", "latency_ms", time.Since(start).Milliseconds())
		}()
		return c.JSON(resp)
	}

	slog.Info("Cache empty, fetching live data", "request_id", c.Locals("requestid"))
	startTime := time.Now()
	summary, err := generateDevLog()
	latency := time.Since(startTime)

	metrics.mu.Lock()
	metrics.totalLatency += latency
	if err != nil {
		metrics.llmErrors++
	}
	metrics.mu.Unlock()

	if err != nil {
		slog.Error("Live fetch failed", "error", err, "latency_ms", latency.Milliseconds(), "request_id", c.Locals("requestid"))
		return c.Status(500).JSON(fiber.Map{"error": "internal server error", "summary": "System Update: Offline (Retrying...)"})
	}

	newResp := DevLogResponse{
		Summary:    summary,
		LastUpdate: time.Now().Format("2006-01-02 15:04:05"),
		Source:     "live",
	}

	cacheMutex.Lock()
	cache.Response = newResp
	cache.ExpiresAt = time.Now().Add(CacheDuration)
	cacheMutex.Unlock()

	return c.JSON(newResp)
}

func handleMetrics(c *fiber.Ctx) error {
	metrics.mu.Lock()
	defer metrics.mu.Unlock()
	avgLatency := time.Duration(0)
	if metrics.totalCalls > 0 {
		avgLatency = metrics.totalLatency / time.Duration(metrics.totalCalls)
	}
	return c.JSON(fiber.Map{
		"cache_hits":     metrics.cacheHits,
		"cache_misses":   metrics.cacheMisses,
		"llm_errors":     metrics.llmErrors,
		"total_requests": metrics.totalCalls,
		"avg_latency_ms": avgLatency.Milliseconds(),
		"uptime_seconds": int64(time.Since(metrics.startTime).Seconds()),
	})
}

func handleSubmitContact(c *fiber.Ctx) error {
	if db == nil {
		return c.Status(503).JSON(fiber.Map{"error": "service unavailable"})
	}

	var sub ContactSubmission
	if err := c.BodyParser(&sub); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	if strings.TrimSpace(sub.Name) == "" {
		return c.Status(400).JSON(fiber.Map{"error": "name is required"})
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(sub.Email) {
		return c.Status(400).JSON(fiber.Map{"error": "invalid email address"})
	}

	if len(strings.TrimSpace(sub.Message)) < 10 {
		return c.Status(400).JSON(fiber.Map{"error": "message must be at least 10 characters"})
	}

	if len(sub.Topic) > 200 {
		return c.Status(400).JSON(fiber.Map{"error": "topic too long (max 200 characters)"})
	}

	_, err := db.Exec("INSERT INTO contact_messages (name, email, topic, message) VALUES ($1, $2, $3, $4)",
		sub.Name, sub.Email, sub.Topic, sub.Message)
	if err != nil {
		slog.Error("Failed to insert contact message", "error", err)
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.Status(201).JSON(fiber.Map{"success": true, "message": "Message sent successfully."})
}

func handleAdminContact(c *fiber.Ctx) error {
	if c.Query("key") != os.Getenv("CONTACT_SECRET") {
		return c.Status(401).SendString("Unauthorized")
	}

	if db == nil {
		return c.Status(503).JSON(fiber.Map{"error": "service unavailable"})
	}

	rows, err := db.Query("SELECT id, name, email, topic, message, is_read, created_at FROM contact_messages ORDER BY created_at DESC LIMIT 50")
	if err != nil {
		slog.Error("Failed to query contact messages", "error", err)
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
			m.ID, htmlEscape(m.Name), htmlEscape(m.Email), htmlEscape(m.Topic), htmlEscape(truncated), m.CreatedAt.Format("Jan 2, 2006 15:04"), statusBadge))
	}

	htmlBuilder.WriteString(`</tbody>
</table>
</body>
</html>`)

	c.Type("text/html")
	return c.SendString(htmlBuilder.String())
}

func handleGetMessages(c *fiber.Ctx) error {
	if c.Query("key") != os.Getenv("CONTACT_SECRET") {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	if db == nil {
		return c.Status(503).JSON(fiber.Map{"error": "service unavailable"})
	}

	rows, err := db.Query("SELECT id, name, email, topic, message, is_read, created_at FROM contact_messages ORDER BY created_at DESC LIMIT 50")
	if err != nil {
		slog.Error("Failed to query contact messages", "error", err)
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

func handleMarkRead(c *fiber.Ctx) error {
	if c.Query("key") != os.Getenv("CONTACT_SECRET") {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	if db == nil {
		return c.Status(503).JSON(fiber.Map{"error": "service unavailable"})
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}

	result, err := db.Exec("UPDATE contact_messages SET is_read = TRUE WHERE id = $1", id)
	if err != nil {
		slog.Error("Failed to mark message as read", "error", err)
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "message not found"})
	}

	return c.JSON(fiber.Map{"success": true})
}

// --- Certificate Handlers ---

func handleGetCertificates(c *fiber.Ctx) error {
	if db == nil {
		return c.Status(503).JSON(fiber.Map{"error": "service unavailable"})
	}

	rows, err := db.Query("SELECT id, title, issuer, date, credential_url, image_url, tags, is_verified, display_order, created_at FROM certificates ORDER BY date DESC NULLS LAST, display_order ASC, id DESC")
	if err != nil {
		slog.Error("Failed to query certificates", "error", err)
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}
	defer rows.Close()

	certs := []Certificate{}
	for rows.Next() {
		var cert Certificate
		var createdAt time.Time
		if err := rows.Scan(&cert.ID, &cert.Title, &cert.Issuer, &cert.Date, &cert.CredentialURL, &cert.ImageURL, pq.Array(&cert.Tags), &cert.IsVerified, &cert.DisplayOrder, &createdAt); err != nil {
			slog.Error("Failed to scan certificate row", "error", err)
			continue
		}
		cert.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
		certs = append(certs, cert)
	}

	return c.JSON(fiber.Map{"certificates": certs})
}

func adminCheck(c *fiber.Ctx) bool {
	return c.Query("key") == os.Getenv("CONTACT_SECRET")
}

func handleAdminGetCertificates(c *fiber.Ctx) error {
	if !adminCheck(c) {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	return handleGetCertificates(c)
}

func handleAdminCreateCertificate(c *fiber.Ctx) error {
	if !adminCheck(c) {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	if db == nil {
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
	if input.ImageURL != "" && !strings.HasPrefix(input.ImageURL, "http://") && !strings.HasPrefix(input.ImageURL, "https://") {
		return c.Status(400).JSON(fiber.Map{"error": "image_url must be a valid HTTP(S) URL"})
	}
	if input.CredentialURL != "" && !strings.HasPrefix(input.CredentialURL, "http://") && !strings.HasPrefix(input.CredentialURL, "https://") {
		return c.Status(400).JSON(fiber.Map{"error": "credential_url must be a valid HTTP(S) URL"})
	}

	if input.Tags == nil {
		input.Tags = []string{}
	}

	var id int
	err := db.QueryRow(
		"INSERT INTO certificates (title, issuer, date, credential_url, image_url, tags, is_verified, display_order) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id",
		input.Title, input.Issuer, input.Date, input.CredentialURL, input.ImageURL, pq.Array(input.Tags), input.IsVerified, input.DisplayOrder,
	).Scan(&id)
	if err != nil {
		slog.Error("Failed to insert certificate", "error", err)
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.Status(201).JSON(fiber.Map{"success": true, "id": id})
}

func handleAdminUpdateCertificate(c *fiber.Ctx) error {
	if !adminCheck(c) {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	if db == nil {
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
	if input.ImageURL != "" && !strings.HasPrefix(input.ImageURL, "http://") && !strings.HasPrefix(input.ImageURL, "https://") {
		return c.Status(400).JSON(fiber.Map{"error": "image_url must be a valid HTTP(S) URL"})
	}
	if input.CredentialURL != "" && !strings.HasPrefix(input.CredentialURL, "http://") && !strings.HasPrefix(input.CredentialURL, "https://") {
		return c.Status(400).JSON(fiber.Map{"error": "credential_url must be a valid HTTP(S) URL"})
	}

	result, err := db.Exec(
		"UPDATE certificates SET title=$1, issuer=$2, date=$3, credential_url=$4, image_url=$5, tags=$6, is_verified=$7, display_order=$8 WHERE id=$9",
		input.Title, input.Issuer, input.Date, input.CredentialURL, input.ImageURL, pq.Array(input.Tags), input.IsVerified, input.DisplayOrder, id,
	)
	if err != nil {
		slog.Error("Failed to update certificate", "error", err)
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "certificate not found"})
	}

	return c.JSON(fiber.Map{"success": true})
}

func handleAdminDeleteCertificate(c *fiber.Ctx) error {
	if !adminCheck(c) {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	if db == nil {
		return c.Status(503).JSON(fiber.Map{"error": "service unavailable"})
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}

	result, err := db.Exec("DELETE FROM certificates WHERE id=$1", id)
	if err != nil {
		slog.Error("Failed to delete certificate", "error", err)
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "certificate not found"})
	}

	return c.JSON(fiber.Map{"success": true})
}

func handleAdminReorderCertificates(c *fiber.Ctx) error {
	if !adminCheck(c) {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	if db == nil {
		return c.Status(503).JSON(fiber.Map{"error": "service unavailable"})
	}

	var input struct {
		Order []int `json:"order"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	for i, id := range input.Order {
		if _, err := db.Exec("UPDATE certificates SET display_order=$1 WHERE id=$2", i, id); err != nil {
			slog.Error("Failed to reorder certificate", "error", err)
			return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
		}
	}

	return c.JSON(fiber.Map{"success": true})
}

func handleAdminUploadImage(c *fiber.Ctx) error {
	if !adminCheck(c) {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "no image file provided"})
	}

	mimeType := file.Header.Get("Content-Type")
	allowedMIME := map[string]bool{
		"image/png":  true,
		"image/jpeg": true,
		"image/gif":  true,
		"image/webp": true,
	}
	if !allowedMIME[mimeType] {
		return c.Status(400).JSON(fiber.Map{"error": "invalid MIME type, use PNG/JPG/GIF/WEBP"})
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := map[string]bool{".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".webp": true}
	if !allowed[ext] {
		return c.Status(400).JSON(fiber.Map{"error": "unsupported format, use PNG/JPG/JPEG/GIF/WEBP"})
	}

	if file.Size > 10<<20 {
		return c.Status(400).JSON(fiber.Map{"error": "file too large, max 10MB"})
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to read file"})
	}
	defer src.Close()

	sniff := make([]byte, 512)
	if _, err := io.ReadFull(src, sniff); err != nil && err != io.ErrUnexpectedEOF {
		return c.Status(500).JSON(fiber.Map{"error": "failed to validate file"})
	}
	sniffMIME := http.DetectContentType(sniff)
	if !allowedMIME[sniffMIME] {
		return c.Status(400).JSON(fiber.Map{"error": "file content does not match image type"})
	}

	src.Seek(0, io.SeekStart)
	data, err := io.ReadAll(src)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to read file data"})
	}

	hash := sha256.Sum256(data)
	filename := hex.EncodeToString(hash[:8]) + ext
	savePath := filepath.Join("static", "uploads", filename)

	if err := os.WriteFile(savePath, data, 0644); err != nil {
		slog.Error("Failed to save uploaded image", "error", err)
		return c.Status(500).JSON(fiber.Map{"error": "failed to save image"})
	}

	scheme := "http"
	if c.Protocol() == "https" || strings.HasPrefix(c.Hostname(), "sentinel") {
		scheme = "https"
	}
	url := fmt.Sprintf("%s://%s/static/uploads/%s", scheme, c.Hostname(), filename)

	return c.JSON(fiber.Map{"url": url, "filename": filename})
}

// --- Badge Handlers ---

func handleGetBadges(c *fiber.Ctx) error {
	if db == nil {
		return c.Status(503).JSON(fiber.Map{"error": "service unavailable"})
	}

	rows, err := db.Query("SELECT id, name, image_url, credential_url, rarity, category, important, display_order, created_at FROM badges ORDER BY important DESC, display_order ASC, id DESC")
	if err != nil {
		slog.Error("Failed to query badges", "error", err)
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}
	defer rows.Close()

	badges := []Badge{}
	for rows.Next() {
		var b Badge
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

func handleAdminGetBadges(c *fiber.Ctx) error {
	if !adminCheck(c) {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	return handleGetBadges(c)
}

func handleAdminCreateBadge(c *fiber.Ctx) error {
	if !adminCheck(c) {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	if db == nil {
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

	if input.ImageURL != "" && !strings.HasPrefix(input.ImageURL, "http://") && !strings.HasPrefix(input.ImageURL, "https://") {
		return c.Status(400).JSON(fiber.Map{"error": "image_url must be a valid HTTP(S) URL"})
	}

	if input.DisplayOrder < 0 {
		input.DisplayOrder = 0
	}

	var id int
	err := db.QueryRow(
		"INSERT INTO badges (name, image_url, credential_url, rarity, category, important, display_order) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		input.Name, input.ImageURL, input.CredentialURL, rarity, input.Category, input.Important, input.DisplayOrder,
	).Scan(&id)
	if err != nil {
		slog.Error("Failed to insert badge", "error", err)
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.Status(201).JSON(fiber.Map{"success": true, "id": id})
}

func handleAdminUpdateBadge(c *fiber.Ctx) error {
	if !adminCheck(c) {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	if db == nil {
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

	result, err := db.Exec(
		"UPDATE badges SET name=$1, image_url=$2, credential_url=$3, rarity=$4, category=$5, important=$6, display_order=$7 WHERE id=$8",
		input.Name, input.ImageURL, input.CredentialURL, rarity, input.Category, input.Important, input.DisplayOrder, id,
	)
	if err != nil {
		slog.Error("Failed to update badge", "error", err)
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "badge not found"})
	}

	return c.JSON(fiber.Map{"success": true})
}

func handleAdminDeleteBadge(c *fiber.Ctx) error {
	if !adminCheck(c) {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	if db == nil {
		return c.Status(503).JSON(fiber.Map{"error": "service unavailable"})
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}

	result, err := db.Exec("DELETE FROM badges WHERE id=$1", id)
	if err != nil {
		slog.Error("Failed to delete badge", "error", err)
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "badge not found"})
	}

	return c.JSON(fiber.Map{"success": true})
}

func handleAdminReorderBadges(c *fiber.Ctx) error {
	if !adminCheck(c) {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	if db == nil {
		return c.Status(503).JSON(fiber.Map{"error": "service unavailable"})
	}

	var input struct {
		Order []int `json:"order"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	for i, id := range input.Order {
		if _, err := db.Exec("UPDATE badges SET display_order=$1 WHERE id=$2", i, id); err != nil {
			slog.Error("Failed to reorder badge", "error", err)
			return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
		}
	}

	return c.JSON(fiber.Map{"success": true})
}

func htmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}

// --- Logic ---

func generateDevLog() (string, error) {
	breaker.mu.Lock()
	if breaker.failures >= maxFailures {
		if time.Since(breaker.lastFailure) < breakerTimeout {
			breaker.mu.Unlock()
			slog.Warn("Circuit breaker open — skipping OpenRouter call")
			return QuotaFallback, nil
		}
		breaker.failures = 0
	}
	breaker.mu.Unlock()

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		return "System Error: Neural Link Disconnected (Missing API Key). Please configure the satellite uplink.", nil
	}

	events, err := fetchGitHubEvents()
	if err != nil {
		return "", err
	}

	responsesMu.Lock()
	diversityInstruction := ""
	if len(lastResponses) > 0 {
		diversityInstruction = fmt.Sprintf("Avoid repeating these recent responses: %s", strings.Join(lastResponses, "; "))
	}
	responsesMu.Unlock()

	prompt := fmt.Sprintf(`You are a senior software engineer named Pranav. A friend just asked "what have you been working on?" — answer naturally.

Today is %s. Here's what you've pushed to GitHub recently:

%s

Write ONE sentence (max 40 words) like you're telling a friend. Mention the project and what you actually did. Be casual, specific, and human.

%s`, time.Now().Weekday().String(), events, diversityInstruction)

	summary, err := callOpenRouter(apiKey, prompt)
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "429") || strings.Contains(errStr, "RESOURCE_EXHAUSTED") || strings.Contains(errStr, "quota") {
			breaker.mu.Lock()
			breaker.failures++
			breaker.lastFailure = time.Now()
			breaker.mu.Unlock()
			slog.Warn("OpenRouter quota exhausted", "failures", breaker.failures)
			return QuotaFallback, nil
		}
		return "", err
	}

	responsesMu.Lock()
	lastResponses = append(lastResponses, summary)
	if len(lastResponses) > 3 {
		lastResponses = lastResponses[len(lastResponses)-3:]
	}
	responsesMu.Unlock()

	return summary, nil
}

func fetchGitHubEvents() (string, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/events/public", GitHubUsername)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("github api error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var events []map[string]interface{}
	if err := json.Unmarshal(body, &events); err != nil {
		return "", err
	}

	relevantTypes := map[string]bool{
		"PushEvent":    true,
		"CreateEvent":  true,
		"ReleaseEvent": true,
		"IssuesEvent":  true,
	}

	var filtered []map[string]interface{}
	for _, event := range events {
		eventType, _ := event["type"].(string)
		if relevantTypes[eventType] {
			filtered = append(filtered, event)
		}
	}
	events = filtered

	var summaryBuilder strings.Builder
	count := 0
	for _, event := range events {
		if count >= 30 {
			break
		}
		eventType, _ := event["type"].(string)
		repo, _ := event["repo"].(map[string]interface{})
		repoName, _ := repo["name"].(string)

		summaryBuilder.WriteString(fmt.Sprintf("- %s on %s", eventType, repoName))

		if payload, ok := event["payload"].(map[string]interface{}); ok {
			if commits, ok := payload["commits"].([]interface{}); ok {
				for _, c := range commits {
					commit, _ := c.(map[string]interface{})
					msg, _ := commit["message"].(string)
					summaryBuilder.WriteString(fmt.Sprintf(": %s", msg))
				}
			}
		}
		summaryBuilder.WriteString("\n")
		count++
	}

	return summaryBuilder.String(), nil
}

func callOpenRouter(apiKey, text string) (string, error) {
	url := "https://openrouter.ai/api/v1/chat/completions"

	requestBody, _ := json.Marshal(map[string]interface{}{
		"model": OpenRouterModel,
		"messages": []interface{}{
			map[string]interface{}{
				"role":    "user",
				"content": text,
			},
		},
	})

	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("openrouter api error: %s", string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if message, ok := choice["message"].(map[string]interface{}); ok {
				if content, ok := message["content"].(string); ok {
					return content, nil
				}
			}
		}
	}

	return "Analysis complete. Systems nominal (Default Response).", nil
}
