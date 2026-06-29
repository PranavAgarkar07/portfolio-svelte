package repository

import (
	"database/sql"
	"encoding/json"
	"log/slog"

	"github.com/lib/pq"

	"portfolio-backend/internal/model"
)

var postCols = "bp.id, bp.slug, bp.title, bp.content_md, bp.excerpt, bp.cover_image, COALESCE(bp.images, '[]'::jsonb), COALESCE(bp.tags, '{}'::text[]), bp.published, bp.published_at, bp.author_id, u.name, bp.created_at, bp.updated_at"
var postJoin = "FROM blog_posts bp JOIN users u ON bp.author_id = u.id"

func GetPublishedPosts(db *sql.DB) (*sql.Rows, error) {
	rows, err := db.Query(`
		SELECT ` + postCols + `
		` + postJoin + `
		WHERE bp.published = TRUE
		ORDER BY bp.published_at DESC
	`)
	if err != nil {
		slog.Error("Failed to query published posts", "error", err)
	}
	return rows, err
}

func GetPostIDBySlug(db *sql.DB, slug string) (string, error) {
	var id string
	err := db.QueryRow("SELECT id FROM blog_posts WHERE slug = $1", slug).Scan(&id)
	if err != nil {
		slog.Error("Failed to get post ID by slug", "error", err)
	}
	return id, err
}

func GetPostBySlug(db *sql.DB, slug string) (*sql.Row, error) {
	row := db.QueryRow(`
		SELECT ` + postCols + `
		` + postJoin + `
		WHERE bp.slug = $1
	`, slug)
	return row, nil
}

func GetAllPosts(db *sql.DB) (*sql.Rows, error) {
	rows, err := db.Query(`
		SELECT ` + postCols + `
		` + postJoin + `
		ORDER BY bp.created_at DESC
	`)
	if err != nil {
		slog.Error("Failed to query all posts", "error", err)
	}
	return rows, err
}

func GetPostsByAuthor(db *sql.DB, authorID string) (*sql.Rows, error) {
	rows, err := db.Query(`
		SELECT ` + postCols + `
		` + postJoin + `
		WHERE bp.author_id = $1
		ORDER BY bp.created_at DESC
	`, authorID)
	if err != nil {
		slog.Error("Failed to query posts by author", "error", err)
	}
	return rows, err
}

func GetPostByID(db *sql.DB, id string) (*sql.Row, error) {
	row := db.QueryRow(`
		SELECT ` + postCols + `
		` + postJoin + `
		WHERE bp.id = $1
	`, id)
	return row, nil
}

func CreatePost(db *sql.DB, slug, title, contentMD, excerpt, coverImage string, images []model.BlogImage, tags []string, published bool, authorID string) (string, error) {
	imagesJSON, err := json.Marshal(images)
	if err != nil {
		return "", err
	}
	var id string
	err = db.QueryRow(
		`INSERT INTO blog_posts (slug, title, content_md, excerpt, cover_image, images, tags, published, author_id, published_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, CASE WHEN $8 THEN NOW() ELSE NULL END)
		 RETURNING id`,
		slug, title, contentMD, excerpt, coverImage, imagesJSON, pq.Array(tags), published, authorID,
	).Scan(&id)
	if err != nil {
		slog.Error("Failed to create blog post", "error", err)
	}
	return id, err
}

func UpdatePost(db *sql.DB, id, slug, title, contentMD, excerpt, coverImage string, images []model.BlogImage, tags []string, published bool) error {
	imagesJSON, err := json.Marshal(images)
	if err != nil {
		return err
	}
	_, err = db.Exec(
		`UPDATE blog_posts SET slug=$1, title=$2, content_md=$3, excerpt=$4, cover_image=$5, images=$6, tags=$7, published=$8,
		 published_at = CASE WHEN $8 AND published_at IS NULL THEN NOW() ELSE published_at END,
		 updated_at = NOW()
		 WHERE id=$9`,
		slug, title, contentMD, excerpt, coverImage, imagesJSON, pq.Array(tags), published, id,
	)
	if err != nil {
		slog.Error("Failed to update blog post", "error", err)
	}
	return err
}

func DeletePost(db *sql.DB, id string) error {
	_, err := db.Exec("DELETE FROM blog_posts WHERE id = $1", id)
	if err != nil {
		slog.Error("Failed to delete blog post", "error", err)
	}
	return err
}
