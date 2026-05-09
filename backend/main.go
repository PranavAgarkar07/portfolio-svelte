package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/joho/godotenv"
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

var (
	cache         CachedData
	cacheMutex    sync.Mutex
	breaker       CircuitBreaker
	lastResponses []string
	responsesMu   sync.Mutex
	metrics       Metrics
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

	prompt := fmt.Sprintf(`You are a senior software engineer named Pranav writing a quick personal status update.
Today is %s. Your recent GitHub activity:

%s

Examples of good responses (use these as tone & style reference):
- "Just landed a full CI/CD overhaul — Dockerized the backend and cut deploy times in half."
- "Spent the week refactoring the auth layer. Token validation is 40ms faster now."
- "Dropped a new P2P transfer feature in BeamSync — QR pairing works on the first try."
- "Been debugging a nasty race condition in the WebSocket handler. Found it — was a missing mutex unlock."

Write ONE sentence (max 20 words) summarizing the work you've been doing recently. 
Be specific about which project and what kind of work. Sound human, not robotic.

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
