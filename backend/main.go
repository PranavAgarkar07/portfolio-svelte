package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	_ "github.com/jackc/pgx/v5/stdlib"

	"portfolio-backend/internal/analytics"
	"portfolio-backend/internal/config"
	"portfolio-backend/internal/handler"
	"portfolio-backend/internal/model"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	var db *sql.DB
	if cfg.DatabaseURL != "" {
		db, err = sql.Open("pgx", cfg.DatabaseURL)
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
				db.SetMaxOpenConns(10)
				db.SetMaxIdleConns(5)
				db.SetConnMaxLifetime(30 * time.Minute)
				db.SetConnMaxIdleTime(5 * time.Minute)

				if os.Getenv("RUN_MIGRATIONS") == "true" {
					slog.Info("Running database migrations")
					runMigrations(db)
					runURLFixups(db)
				}
			}
		}
	}
	if db != nil {
		defer db.Close()
	}

	var s3Client *s3.Client
	if cfg.IsLambda {
		awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(), awsconfig.WithRegion(model.S3Region))
		if err != nil {
			slog.Warn("Failed to load AWS config", "error", err)
		} else {
			s3Client = s3.NewFromConfig(awsCfg)
			slog.Info("S3 client initialized")
		}
	}

	app := fiber.New()
	app.Use(requestid.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://pranavagarkar07.github.io, http://localhost:5173, http://localhost:5174, http://localhost:4173",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, X-Analytics-Key",
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

	h := handler.New(cfg, db, s3Client)
	h.RegisterRoutes(app)

	if db != nil {
		analyticsSecret := cfg.AnalyticsSecret
		analytics.RegisterRoutes(app, db, analyticsSecret)
	}

	if !cfg.IsLambda {
		app.Static("/static", "./static")
		app.Static("/badges", "./static/badges")
	}

	go handler.ResetLikeRate()

	if cfg.IsLambda {
		slog.Info("Sentinel starting in Lambda mode")
		adapter := fiberadapter.New(app)
		lambda.Start(adapter.ProxyV2)
	} else {
		slog.Info("Sentinel starting", "port", cfg.Port)

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		go func() {
			if err := app.Listen(":" + cfg.Port); err != nil {
				slog.Error("Server failed to start", "error", err)
				os.Exit(1)
			}
		}()

		<-ctx.Done()
		slog.Info("Shutting down gracefully...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := app.ShutdownWithContext(shutdownCtx); err != nil {
			slog.Error("Server forced to shutdown", "error", err)
		}
	}
}

func runMigrations(db *sql.DB) {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS contact_messages (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL,
			topic TEXT DEFAULT '',
			message TEXT NOT NULL,
			is_read BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`ALTER TABLE contact_messages ADD COLUMN IF NOT EXISTS topic TEXT DEFAULT ''`,
		`CREATE TABLE IF NOT EXISTS certificates (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			issuer TEXT NOT NULL,
			date TEXT DEFAULT '',
			credential_url TEXT DEFAULT '',
			image_url TEXT DEFAULT '',
			tags TEXT[] DEFAULT '{}',
			is_verified BOOLEAN DEFAULT FALSE,
			display_order INT DEFAULT 0,
			created_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS badges (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			image_url TEXT DEFAULT '',
			credential_url TEXT DEFAULT '',
			rarity TEXT DEFAULT 'common',
			category TEXT DEFAULT '',
			important BOOLEAN DEFAULT false,
			display_order INT DEFAULT 0,
			created_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`ALTER TABLE badges ADD COLUMN IF NOT EXISTS category TEXT DEFAULT ''`,
		`ALTER TABLE badges ADD COLUMN IF NOT EXISTS important BOOLEAN DEFAULT false`,
		`CREATE TABLE IF NOT EXISTS sessions (
			id TEXT PRIMARY KEY,
			ip_hash TEXT NOT NULL,
			country TEXT DEFAULT '',
			city TEXT DEFAULT '',
			referrer TEXT DEFAULT '',
			device TEXT DEFAULT '',
			os TEXT DEFAULT '',
			browser TEXT DEFAULT '',
			theme TEXT DEFAULT '',
			created_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS events (
			id SERIAL PRIMARY KEY,
			session_id TEXT NOT NULL,
			type TEXT NOT NULL,
			target TEXT DEFAULT '',
			value TEXT DEFAULT '',
			ts TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_events_ts ON events(ts)`,
		`CREATE INDEX IF NOT EXISTS idx_events_type_target ON events(type, target)`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_created ON sessions(created_at)`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_ip ON sessions(ip_hash)`,
		`ALTER TABLE sessions ADD COLUMN IF NOT EXISTS source TEXT DEFAULT ''`,
		`CREATE TABLE IF NOT EXISTS project_likes (
			project_name TEXT NOT NULL,
			visitor_token TEXT NOT NULL,
			liked BOOLEAN NOT NULL DEFAULT TRUE,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW(),
			PRIMARY KEY (project_name, visitor_token)
		)`,
	}
	for _, m := range migrations {
		if _, err := db.Exec(m); err != nil {
			slog.Warn("Migration failed", "error", err, "query", m[:60])
		}
	}
}

func runURLFixups(db *sql.DB) {
	fixups := []string{
		`UPDATE badges SET image_url = REPLACE(image_url, 'https://sentinel-backend-4x3i.onrender.com/static/uploads/', '/badges/') WHERE image_url LIKE '%/static/uploads/%'`,
		`UPDATE badges SET image_url = REPLACE(image_url, 'http://localhost:8080/static/uploads/', '/badges/') WHERE image_url LIKE '%/static/uploads/%'`,
		`UPDATE certificates SET image_url = REPLACE(image_url, 'https://sentinel-backend-4x3i.onrender.com/static/', 'https://portfolio-uploads-sentinel.s3.ap-south-1.amazonaws.com/static/') WHERE image_url LIKE '%sentinel-backend%'`,
		`UPDATE certificates SET image_url = REPLACE(image_url, 'http://localhost:8080/static/', 'https://portfolio-uploads-sentinel.s3.ap-south-1.amazonaws.com/static/') WHERE image_url LIKE '%localhost%'`,
	}
	for _, f := range fixups {
		if _, err := db.Exec(f); err != nil {
			slog.Warn("URL fixup failed", "error", err)
		}
	}
}
