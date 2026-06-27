package model

import (
	"sync"
	"time"
)

const (
	GitHubUsername = "PranavAgarkar07"
	CacheDuration  = 4 * time.Hour
	RequestTimeout = 15 * time.Second
	QuotaFallback  = "Deep in the code — pushing updates across multiple projects. Check back soon!"
	OpenRouterModel = "openrouter/free"
	MaxFailures    = 3
	BreakerTimeout = 15 * time.Minute
	S3Bucket       = "portfolio-uploads-sentinel"
	S3Region       = "ap-south-1"
)

type DevLogResponse struct {
	Summary    string `json:"summary"`
	LastUpdate string `json:"last_update"`
	Source     string `json:"source"`
}

type CachedData struct {
	Response  DevLogResponse
	ExpiresAt time.Time
}

type CircuitBreaker struct {
	Mu          sync.Mutex
	Failures    int
	LastFailure time.Time
}

type Metrics struct {
	Mu           sync.Mutex
	CacheHits    int64
	CacheMisses  int64
	LLMErrors    int64
	TotalLatency time.Duration
	TotalCalls   int64
	StartTime    time.Time
}

type ContactSubmission struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Topic   string `json:"topic"`
	Message string `json:"message"`
}

type Certificate struct {
	ID            int      `json:"id"`
	Title         string   `json:"title"`
	Issuer        string   `json:"issuer"`
	Date          string   `json:"date"`
	CredentialURL string   `json:"credential_url"`
	ImageURL      string   `json:"image_url"`
	ThumbURL      string   `json:"thumb_url"`
	Tags          []string `json:"tags"`
	IsVerified    bool     `json:"is_verified"`
	DisplayOrder  int      `json:"display_order"`
	CreatedAt     string   `json:"created_at"`
}

type Badge struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	ImageURL      string `json:"image_url"`
	CredentialURL string `json:"credential_url"`
	Rarity        string `json:"rarity"`
	Category      string `json:"category"`
	Important     bool   `json:"important"`
	DisplayOrder  int    `json:"display_order"`
	CreatedAt     string `json:"created_at"`
}

type LikeRow struct {
	ProjectName string `json:"project_name"`
	Count       int    `json:"count"`
	Liked       bool   `json:"liked"`
}
