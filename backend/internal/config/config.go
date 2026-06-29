package config

import (
	"fmt"
	"os"
	"log/slog"

	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	DatabaseURL       string
	ContactSecret     string
	AnalyticsSecret   string
	OpenRouterAPIKey  string
	GitHubToken       string
	IsLambda          bool
	GoogleClientID    string
	GoogleClientSecret string
	JWTSecret         string
	FrontendURL       string
	GoogleCallbackURL string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		slog.Info("No .env file found, relying on system env vars")
	}

	cfg := &Config{
		Port:               getEnv("PORT", "8080"),
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		ContactSecret:      os.Getenv("CONTACT_SECRET"),
		AnalyticsSecret:    os.Getenv("ANALYTICS_SECRET"),
		OpenRouterAPIKey:   os.Getenv("OPENROUTER_API_KEY"),
		GitHubToken:        os.Getenv("GITHUB_TOKEN"),
		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		JWTSecret:          os.Getenv("JWT_SECRET"),
		FrontendURL:        getEnv("FRONTEND_URL", "http://localhost:5173"),
		GoogleCallbackURL:  os.Getenv("GOOGLE_CALLBACK_URL"),
		IsLambda:           os.Getenv("AWS_LAMBDA_RUNTIME_API") != "",
	}

	if cfg.AnalyticsSecret == "" {
		cfg.AnalyticsSecret = cfg.ContactSecret
	}

	if cfg.OpenRouterAPIKey != "" {
		slog.Info("OpenRouter API key loaded", "present", true)
	} else {
		slog.Warn("OpenRouter API key is empty or invalid")
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	required := map[string]string{
		"DATABASE_URL":       c.DatabaseURL,
	}

	var missing []string
	for name, val := range required {
		if val == "" {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %v", missing)
	}
	return nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
