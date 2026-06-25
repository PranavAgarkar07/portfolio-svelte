package analytics

import (
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, db *sql.DB, secret string) {
	g := app.Group("/api/v1/analytics")

	startRateCleanup()

	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		CleanupOldData(db)
		for range ticker.C {
			CleanupOldData(db)
		}
	}()

	g.Post("/session", RateLimit(5, 1*time.Minute), HandleCreateSession(db))
	g.Post("/events", RateLimit(10, 1*time.Minute), HandleCreateEvents(db))
	g.Get("/dashboard", Auth(secret), HandleDashboard(db))
}
