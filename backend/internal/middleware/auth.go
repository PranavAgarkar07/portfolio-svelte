package middleware

import (
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

func GenerateJWT(secret string, userID, email, role, name string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		Name:   name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func AuthRequired(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(401).JSON(fiber.Map{"error": "missing authorization header"})
		}
		tokenStr := ""
		fmt.Sscanf(auth, "Bearer %s", &tokenStr)
		if tokenStr == "" {
			return c.Status(401).JSON(fiber.Map{"error": "invalid authorization format"})
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{"error": "invalid or expired token"})
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)
		c.Locals("user_role", claims.Role)
		c.Locals("user_name", claims.Name)
		return c.Next()
	}
}

func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("user_role").(string)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
		}
		for _, r := range roles {
			if role == r {
				return c.Next()
			}
		}
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
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
