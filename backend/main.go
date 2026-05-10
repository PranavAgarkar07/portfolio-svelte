package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// --- Configuration ---
const (
	GitHubUsername = "PranavAgarkar07"
	CacheDuration  = 4 * time.Hour
	RequestTimeout = 15 * time.Second

	QuotaFallback = "Deep in the code — pushing updates across multiple projects. Check back soon!"

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
	geminiErrors int64
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

var (
	cache         CachedData
	cacheMutex    sync.Mutex
	breaker       CircuitBreaker
	lastResponses []string
	responsesMu   sync.Mutex
	metrics       Metrics
	db            *sql.DB
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

	key := os.Getenv("GEMINI_API_KEY")
	if len(key) > 10 {
		slog.Info("API key loaded", "prefix", key[:4], "suffix", key[len(key)-4:])
	} else {
		slog.Warn("API key is empty or invalid")
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
			}
		}
	}
	if db != nil {
		defer db.Close()
	}

	app := fiber.New()

	app.Use(requestid.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://pranavagarkar07.github.io, http://localhost:5173, http://localhost:4173",
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
				metrics.geminiErrors++
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
		metrics.geminiErrors++
	}
	metrics.mu.Unlock()

	if err != nil {
		slog.Error("Live fetch failed", "error", err, "latency_ms", latency.Milliseconds(), "request_id", c.Locals("requestid"))
		return c.Status(500).JSON(fiber.Map{"error": err.Error(), "summary": "System Update: Offline (Retrying...)"})
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
		"gemini_errors":  metrics.geminiErrors,
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
			slog.Warn("Circuit breaker open — skipping Gemini call")
			return QuotaFallback, nil
		}
		breaker.failures = 0
	}
	breaker.mu.Unlock()

	apiKey := os.Getenv("GEMINI_API_KEY")
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

	summary, err := callGemini(apiKey, prompt)
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "429") || strings.Contains(errStr, "RESOURCE_EXHAUSTED") || strings.Contains(errStr, "quota") {
			breaker.mu.Lock()
			breaker.failures++
			breaker.lastFailure = time.Now()
			breaker.mu.Unlock()
			slog.Warn("Gemini quota exhausted", "failures", breaker.failures)
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

func callGemini(apiKey, text string) (string, error) {
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=" + apiKey

	requestBody, _ := json.Marshal(map[string]interface{}{
		"contents": []interface{}{
			map[string]interface{}{
				"parts": []interface{}{
					map[string]interface{}{
						"text": text,
					},
				},
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

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("gemini api error: %s", string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if candidates, ok := result["candidates"].([]interface{}); ok && len(candidates) > 0 {
		if candidate, ok := candidates[0].(map[string]interface{}); ok {
			if content, ok := candidate["content"].(map[string]interface{}); ok {
				if parts, ok := content["parts"].([]interface{}); ok && len(parts) > 0 {
					if part, ok := parts[0].(map[string]interface{}); ok {
						if textVal, ok := part["text"].(string); ok {
							return textVal, nil
						}
					}
				}
			}
		}
	}

	return "Analysis complete. Systems nominal (Default Response).", nil
}
