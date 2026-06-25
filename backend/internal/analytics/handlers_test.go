package analytics

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestCreateSession_RejectsMissingFields(t *testing.T) {
	app := fiber.New()
	app.Post("/session", HandleCreateSession(nil))

	tests := []struct {
		name   string
		body   string
		status int
	}{
		{"empty body", `{}`, 400},
		{"missing id", `{}`, 400},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/session", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req, 100)
			if resp.StatusCode != tt.status {
				t.Errorf("expected status %d, got %d", tt.status, resp.StatusCode)
			}
		})
	}
}

func TestCreateEvents_RejectsInvalidEvents(t *testing.T) {
	app := fiber.New()
	app.Post("/events", HandleCreateEvents(nil))

	tests := []struct {
		name   string
		body   string
		status int
	}{
		{"empty body", `[]`, 202},
		{"missing session_id", `[{"type": "click"}]`, 400},
		{"missing type", `[{"session_id": "abc"}]`, 400},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/events", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req, 100)
			if resp.StatusCode != tt.status {
				t.Errorf("expected status %d, got %d", tt.status, resp.StatusCode)
			}
		})
	}
}

func TestCreateEvents_RejectsInvalidJSON(t *testing.T) {
	app := fiber.New()
	app.Post("/events", HandleCreateEvents(nil))

	req := httptest.NewRequest("POST", "/events", strings.NewReader("not json"))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 100)
	if resp.StatusCode != 400 {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestDashboard_RejectsNoAuth(t *testing.T) {
	app := fiber.New()
	app.Get("/dashboard", Auth("secret"), func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	req := httptest.NewRequest("GET", "/dashboard", nil)
	resp, _ := app.Test(req, 100)
	if resp.StatusCode != 401 {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}
}

func TestDashboard_AcceptsValidAuth(t *testing.T) {
	app := fiber.New()
	app.Get("/dashboard", Auth("secret"), func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	req := httptest.NewRequest("GET", "/dashboard", nil)
	req.Header.Set("X-Analytics-Key", "secret")
	resp, _ := app.Test(req, 100)
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}
