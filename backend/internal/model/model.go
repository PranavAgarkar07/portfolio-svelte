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
	Tags          Strings  `json:"tags"`
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

type User struct {
	ID        string `json:"id"`
	GoogleID  string `json:"google_id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

type BlogImage struct {
	URL     string `json:"url"`
	Alt     string `json:"alt,omitempty"`
	Caption string `json:"caption,omitempty"`
}

type BlogPost struct {
	ID          string      `json:"id"`
	Slug        string      `json:"slug"`
	Title       string      `json:"title"`
	ContentMD   string      `json:"content_md"`
	Excerpt     string      `json:"excerpt"`
	CoverImage  string      `json:"cover_image"`
	Images      []BlogImage `json:"images,omitempty"`
	Tags        []string    `json:"tags"`
	Published   bool        `json:"published"`
	PublishedAt string      `json:"published_at,omitempty"`
	AuthorID    string      `json:"author_id,omitempty"`
	AuthorName  string      `json:"author_name,omitempty"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
}

type ProjectReview struct {
	UserID      string `json:"user_id"`
	UserName    string `json:"user_name"`
	AvatarURL   string `json:"avatar_url,omitempty"`
	ProjectName string `json:"project_name"`
	Rating      int    `json:"rating"`
	Comment     string `json:"comment,omitempty"`
	CreatedAt   string `json:"created_at"`
}

type BlogComment struct {
	ID        string  `json:"id"`
	PostID    string  `json:"post_id"`
	ParentID  *string `json:"parent_id,omitempty"`
	UserID    string  `json:"user_id"`
	UserName  string  `json:"user_name"`
	AvatarURL string  `json:"avatar_url"`
	Content   string  `json:"content"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	Replies   []BlogComment `json:"replies,omitempty"`
}

type MarqueeItem struct {
	Type       string `json:"type"`
	UserName   string `json:"user_name"`
	AvatarURL  string `json:"avatar_url,omitempty"`
	ProjectName string `json:"project_name"`
	Rating     int    `json:"rating,omitempty"`
	Comment    string `json:"comment,omitempty"`
	CreatedAt  string `json:"created_at"`
}
