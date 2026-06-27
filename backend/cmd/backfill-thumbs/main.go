package main

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/image/draw"
)

const (
	S3Bucket = "portfolio-uploads-sentinel"
	S3Region = "ap-south-1"
)

type Certificate struct {
	ID       int
	ImageURL string
}

func main() {
	slog.Info("Starting thumbnail backfill...")

	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, relying on system env vars")
	}

	// --- DB ---
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		slog.Error("DATABASE_URL not set")
		os.Exit(1)
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		slog.Error("Failed to open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		slog.Error("Failed to ping database", "error", err)
		os.Exit(1)
	}
	slog.Info("Database connected")

	// --- S3 ---
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(S3Region))
	if err != nil {
		slog.Error("Failed to load AWS config", "error", err)
		os.Exit(1)
	}
	s3Client := s3.NewFromConfig(cfg)
	slog.Info("S3 client initialized")

	// --- Fetch certificates with S3 image URLs ---
	rows, err := db.Query(`SELECT id, image_url FROM certificates WHERE image_url LIKE '%amazonaws.com/static/%' AND image_url NOT LIKE '%amazonaws.com/static/thumbs/%'`)
	if err != nil {
		slog.Error("Failed to query certificates", "error", err)
		os.Exit(1)
	}
	defer rows.Close()

	var certs []Certificate
	for rows.Next() {
		var c Certificate
		if err := rows.Scan(&c.ID, &c.ImageURL); err != nil {
			slog.Error("Failed to scan row", "error", err)
			continue
		}
		certs = append(certs, c)
	}

	slog.Info("Found certificates to process", "count", len(certs))

	processed := 0
	skipped := 0
	errors := 0

	for _, cert := range certs {
		if err := processCertificate(s3Client, cert); err != nil {
			slog.Error("Failed to process certificate", "id", cert.ID, "error", err)
			errors++
			continue
		}
		processed++
	}

	slog.Info("Backfill complete",
		"total", len(certs),
		"processed", processed,
		"skipped", skipped,
		"errors", errors,
	)
}

func processCertificate(s3Client *s3.Client, cert Certificate) error {
	// Extract key from S3 URL
	// URL: https://portfolio-uploads-sentinel.s3.ap-south-1.amazonaws.com/static/uploads/filename.ext
	uploadKey := extractKey(cert.ImageURL)
	if uploadKey == "" {
		return fmt.Errorf("could not extract S3 key from URL: %s", cert.ImageURL)
	}

	// Check if thumbnail already exists
	thumbName := thumbFilename(uploadKey)
	thumbKey := "static/thumbs/" + thumbName

	_, err := s3Client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(S3Bucket),
		Key:    aws.String(thumbKey),
	})
	if err == nil {
		slog.Info("Thumbnail already exists, skipping", "id", cert.ID, "key", thumbKey)
		return nil
	}

	// Download original from S3
	data, err := downloadFromS3(s3Client, uploadKey)
	if err != nil {
		return fmt.Errorf("failed to download %s: %w", uploadKey, err)
	}

	// Generate thumbnail
	thumbData, err := generateThumb(data)
	if err != nil {
		return fmt.Errorf("failed to generate thumbnail: %w", err)
	}

	if thumbData == nil {
		slog.Warn("Could not decode image, skipping", "id", cert.ID, "key", uploadKey)
		return nil
	}

	// Upload thumbnail
	_, err = s3Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(S3Bucket),
		Key:         aws.String(thumbKey),
		Body:        bytes.NewReader(thumbData),
		ContentType: aws.String("image/jpeg"),
	})
	if err != nil {
		return fmt.Errorf("failed to upload thumbnail: %w", err)
	}

	slog.Info("Generated thumbnail",
		"id", cert.ID,
		"original", uploadKey,
		"thumb", thumbKey,
		"original_bytes", len(data),
		"thumb_bytes", len(thumbData),
	)
	return nil
}

func downloadFromS3(s3Client *s3.Client, key string) ([]byte, error) {
	result, err := s3Client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(S3Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()
	return io.ReadAll(result.Body)
}

func generateThumb(data []byte) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, nil
	}

	const maxThumbWidth = 640
	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	thumbW, thumbH := w, h
	if w > maxThumbWidth {
		thumbW = maxThumbWidth
		thumbH = h * maxThumbWidth / w
	}

	thumb := image.NewRGBA(image.Rect(0, 0, thumbW, thumbH))
	draw.BiLinear.Scale(thumb, thumb.Bounds(), img, bounds, draw.Over, nil)

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, thumb, &jpeg.Options{Quality: 70}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func extractKey(url string) string {
	idx := strings.Index(url, "/static/")
	if idx < 0 {
		return ""
	}
	afterStatic := url[idx+8:]
	if strings.HasPrefix(afterStatic, "thumbs/") {
		return ""
	}
	return "static/" + afterStatic
}

func thumbFilename(key string) string {
	base := filepath.Base(key)
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext) + "_thumb.jpg"
}


