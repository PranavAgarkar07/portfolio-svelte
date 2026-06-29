package handler

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"portfolio-backend/internal/config"
	"portfolio-backend/internal/middleware"
	"portfolio-backend/internal/repository"
)

type AuthHandler struct {
	OAuthConfig *oauth2.Config
	JWTSecret   string
	FrontendURL string
	DB          *sql.DB
	stateMu     sync.Mutex
	states      map[string]bool
}

func NewAuthHandler(cfg *config.Config, db *sql.DB) *AuthHandler {
	return &AuthHandler{
		OAuthConfig: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  cfg.GoogleCallbackURL,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint:     google.Endpoint,
		},
		JWTSecret:   cfg.JWTSecret,
		FrontendURL: cfg.FrontendURL,
		DB:          db,
		states:      make(map[string]bool),
	}
}

func (a *AuthHandler) generateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	state := hex.EncodeToString(b)
	a.stateMu.Lock()
	a.states[state] = true
	a.stateMu.Unlock()
	go func() {
		time.Sleep(10 * time.Minute)
		a.stateMu.Lock()
		delete(a.states, state)
		a.stateMu.Unlock()
	}()
	return state
}

type googleUserInfo struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func (a *AuthHandler) HandleGoogleURL(c *fiber.Ctx) error {
	state := a.generateState()
	url := a.OAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return c.JSON(fiber.Map{"url": url})
}

func (a *AuthHandler) HandleGoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		return c.Status(400).JSON(fiber.Map{"error": "missing authorization code"})
	}

	if state != "" {
		a.stateMu.Lock()
		delete(a.states, state)
		a.stateMu.Unlock()
	}

	token, err := a.OAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		slog.Error("Failed to exchange OAuth code", "error", err)
		return c.Status(401).JSON(fiber.Map{"error": "failed to authenticate"})
	}

	client := a.OAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		slog.Error("Failed to get user info", "error", err)
		return c.Status(500).JSON(fiber.Map{"error": "failed to get user info"})
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to read user info"})
	}

	var googleUser googleUserInfo
	if err := json.Unmarshal(body, &googleUser); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to parse user info"})
	}

	var userID, role string
	row, err := repository.GetUserByGoogleID(a.DB, googleUser.ID)
	if err == nil && row != nil {
		var id, gid, email, name, avatarURL, r string
		var createdAt time.Time
		if err := row.Scan(&id, &gid, &email, &name, &avatarURL, &r, &createdAt); err == nil {
			userID = id
			role = r
			if avatarURL != googleUser.Picture {
				repository.UpdateUserAvatar(a.DB, id, googleUser.Picture)
			}
		}
	}

	if userID == "" {
		userID, err = repository.CreateUser(a.DB, googleUser.ID, googleUser.Email, googleUser.Name, googleUser.Picture)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to create user"})
		}
		role = "user"
	}

	jwt, err := middleware.GenerateJWT(a.JWTSecret, userID, googleUser.Email, role, googleUser.Name)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to generate token"})
	}

	redirectURL := fmt.Sprintf("%s#token=%s", a.FrontendURL, url.QueryEscape(jwt))
	return c.Redirect(redirectURL, 302)
}

func (a *AuthHandler) HandleAuthMe(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(string)

	row, err := repository.GetUserByID(a.DB, userID)
	if err != nil || row == nil {
		return c.Status(401).JSON(fiber.Map{"error": "user not found"})
	}

	var id, googleID, email, name, avatarURL, role string
	var createdAt time.Time
	if err := row.Scan(&id, &googleID, &email, &name, &avatarURL, &role, &createdAt); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to read user"})
	}

	return c.JSON(fiber.Map{
		"user": fiber.Map{
			"id":         id,
			"email":      email,
			"name":       name,
			"avatar_url": avatarURL,
			"role":       role,
		},
	})
}
