package analytics

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

var (
	rateBuckets sync.Map
	rateMu      sync.Mutex
)

type bucket struct {
	tokens    float64
	lastCheck time.Time
}

func startRateCleanup() {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for range ticker.C {
			rateBuckets.Range(func(key, value interface{}) bool {
				b := value.(*bucket)
				if time.Since(b.lastCheck) > 2*time.Minute {
					rateBuckets.Delete(key)
				}
				return true
			})
		}
	}()
}

func getToken(key string, max float64, interval time.Duration) bool {
	now := time.Now()
	val, _ := rateBuckets.LoadOrStore(key, &bucket{tokens: max, lastCheck: now})
	b := val.(*bucket)

	rateMu.Lock()
	defer rateMu.Unlock()

	elapsed := now.Sub(b.lastCheck).Seconds()
	b.tokens = b.tokens + elapsed*(max/interval.Seconds())
	if b.tokens > max {
		b.tokens = max
	}
	b.lastCheck = now

	if b.tokens >= 1 {
		b.tokens--
		return true
	}
	return false
}

func RateLimit(max float64, interval time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := c.IP() + c.Path()
		if !getToken(key, max, interval) {
			return c.Status(429).JSON(fiber.Map{"error": "rate limit exceeded"})
		}
		return c.Next()
	}
}

func CORS() fiber.Handler {
	return func(c *fiber.Ctx) error {
		origin := c.Get("Origin")
		if origin == "http://localhost:5173" || origin == "http://localhost:4173" || origin == "https://pranavagarkar07.github.io" {
			c.Set("Access-Control-Allow-Origin", origin)
		}
		c.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type, X-Analytics-Key")

		if c.Method() == "OPTIONS" {
			return c.SendStatus(200)
		}
		return c.Next()
	}
}

func Auth(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := c.Get("X-Analytics-Key")
		if key == "" || key != secret {
			return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
		}
		return c.Next()
	}
}
