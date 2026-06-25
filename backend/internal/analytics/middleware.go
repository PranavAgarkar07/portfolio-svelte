package analytics

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

var (
	rateBuckets map[string]*bucket
	rateMu      sync.Mutex
)

type bucket struct {
	tokens    float64
	lastCheck time.Time
}

func init() {
	resetRateLimiter()
}

func resetRateLimiter() {
	rateMu.Lock()
	defer rateMu.Unlock()
	rateBuckets = make(map[string]*bucket)
}

func startRateCleanup() {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for range ticker.C {
			rateMu.Lock()
			for key, b := range rateBuckets {
				if time.Since(b.lastCheck) > 2*time.Minute {
					delete(rateBuckets, key)
				}
			}
			rateMu.Unlock()
		}
	}()
}

func getToken(key string, max float64, interval time.Duration) bool {
	rateMu.Lock()
	defer rateMu.Unlock()

	now := time.Now()
	b, ok := rateBuckets[key]
	if !ok {
		rateBuckets[key] = &bucket{tokens: max - 1, lastCheck: now}
		return true
	}

	elapsed := now.Sub(b.lastCheck).Seconds()
	b.tokens += elapsed * (max / interval.Seconds())
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

func Auth(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := c.Get("X-Analytics-Key")
		if key == "" || key != secret {
			return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
		}
		return c.Next()
	}
}
