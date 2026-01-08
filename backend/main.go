package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

// --- Configuration ---
const (
	GitHubUsername = "PranavAgarkar07"
	CacheDuration  = 5 * time.Minute
)

// --- Structs ---

type DevLogResponse struct {
	Summary    string `json:"summary"`
	LastUpdate string `json:"last_update"`
	Source     string `json:"source"` // "cache" or "live"
}

type CachedData struct {
	Response  DevLogResponse
	ExpiresAt time.Time
}

var (
	cache      CachedData
	cacheMutex sync.Mutex
)

// --- Main ---

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system env vars")
	} else {
		log.Println(".env file loaded successfully")
	}

	key := os.Getenv("GEMINI_API_KEY")
	if len(key) > 10 {
		log.Printf("Loaded API Key: %s...%s", key[:4], key[len(key)-4:])
	} else {
		log.Println("Loaded API Key: [EMPTY] or [INVALID LENGTH]")
	}

	app := fiber.New()

	// Enable CORS for all origins (since frontend is on different port/domain)
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Sentinel API is Online 🟢")
	})

	app.Get("/api/status", handleStatus)

	// Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Sentinel listening on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

// --- Handlers ---

func handleStatus(c *fiber.Ctx) error {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// 1. Check Cache
	if time.Now().Before(cache.ExpiresAt) {
		log.Println("⚡ Serving from CACHE")
		resp := cache.Response
		resp.Source = "cache"
		return c.JSON(resp)
	}

	log.Println("🔄 Cache expired or empty. Fetching LIVE data...")

	// 2. Fetch Live Data
	summary, err := generateDevLog()
	if err != nil {
		// On error, try to return stale cache if available
		if cache.Response.Summary != "" {
			log.Printf("⚠️ Error fetching fresh data: %v. Returning STALE cache.", err)
			resp := cache.Response
			resp.Source = "stale-cache"
			return c.JSON(resp)
		}
		// If no cache, return the error
		return c.Status(500).JSON(fiber.Map{"error": err.Error(), "summary": "System Update: Offline (Retrying...)"})
	}

	// 3. Update Cache
	newResp := DevLogResponse{
		Summary:    summary,
		LastUpdate: time.Now().Format("2006-01-02 15:04:05"),
		Source:     "live",
	}
	cache.Response = newResp
	cache.ExpiresAt = time.Now().Add(CacheDuration)

	return c.JSON(newResp)
}

// --- Logic ---

func generateDevLog() (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "System Error: Neural Link Disconnected (Missing API Key). Please configure the satellite uplink.", nil
	}

	// 1. Fetch GitHub Events
	events, err := fetchGitHubEvents()
	if err != nil {
		return "", err
	}

	// 2. Generate Prompt
	prompt := fmt.Sprintf(`You are Pranav Agarkar, a Senior Software Engineer writing a quick status update for your personal portfolio.
	Analyze your following raw GitHub commit history:
	
	Raw Data:
	%s

	Task: Write a SINGLE, casual but professional sentence about what you've been building lately.
	Tone: Highly personal and authentic. Use "I". Sound like a human engineer talking to a friend.
	Examples:
	- "I've been deep in the backend refactoring the auth middleware for better security."
	- "Just pushed some major updates to the rendering engine to smooth out animations."
	- "Spending the weekend optimizing database queries for the TaskVault project."
	
	Constraint: Keep it under 20 words. No robotic "Pranav has updated" language. Be you.`, events)

	// 3. Call Gemini API
	summary, err := callGemini(apiKey, prompt)
	if err != nil {
		return "", err
	}

	return summary, nil
}

func fetchGitHubEvents() (string, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/events/public", GitHubUsername)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
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

	// Parse JSON to extract "type" and "repo.name" and commit messages pushed
	var events []map[string]interface{}
	if err := json.Unmarshal(body, &events); err != nil {
		return "", err
	}

	// Summarize locally first to save token count
	var summaryBuilder strings.Builder
	count := 0
	for _, event := range events {
		if count >= 30 { // Increased from 15 to 30 to grab more context
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
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-flash-latest:generateContent?key=" + apiKey

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

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
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

	// Extract text from response structure
	// candidates[0].content.parts[0].text
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
