package middleware

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

func AdminAuth(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Get("Authorization") != "Bearer "+secret {
			return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
		}
		return c.Next()
	}
}

func AdminAuthQuery(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Query("key") != secret {
			return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
		}
		return c.Next()
	}
}

var (
	LikeRateMu sync.Mutex
	LikeRate   = map[string]int{}
)

func LikeRateLimit(c *fiber.Ctx) error {
	ip := c.IP()
	LikeRateMu.Lock()
	defer LikeRateMu.Unlock()
	LikeRate[ip]++
	if LikeRate[ip] > 20 {
		return c.Status(429).JSON(fiber.Map{"error": "too many requests"})
	}
	return c.Next()
}

func ResetLikeRate() {
	for range time.Tick(1 * time.Minute) {
		LikeRateMu.Lock()
		LikeRate = map[string]int{}
		LikeRateMu.Unlock()
	}
}
