package repository

import (
	"database/sql"
	"log/slog"
)

func GetProjectLikes(db *sql.DB) (*sql.Rows, error) {
	rows, err := db.Query(`SELECT project_name, COUNT(*) as cnt FROM project_likes WHERE liked = TRUE GROUP BY project_name`)
	if err != nil {
		slog.Error("Failed to query project likes", "error", err)
	}
	return rows, err
}

func GetUserLikes(db *sql.DB, visitorToken string) (*sql.Rows, error) {
	rows, err := db.Query(`SELECT project_name FROM project_likes WHERE visitor_token = $1 AND liked = TRUE`, visitorToken)
	if err != nil {
		slog.Error("Failed to query user likes", "error", err)
	}
	return rows, err
}

func UpsertLike(db *sql.DB, projectName, visitorToken string, liked bool) error {
	_, err := db.Exec(
		`INSERT INTO project_likes (project_name, visitor_token, liked, updated_at)
		 VALUES ($1, $2, $3, NOW())
		 ON CONFLICT (project_name, visitor_token)
		 DO UPDATE SET liked = $3, updated_at = NOW()`,
		projectName, visitorToken, liked,
	)
	if err != nil {
		slog.Error("Failed to upsert project like", "error", err)
	}
	return err
}

func CountProjectLikes(db *sql.DB, projectName string) (int, error) {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM project_likes WHERE project_name = $1 AND liked = TRUE`, projectName).Scan(&count)
	if err != nil {
		slog.Error("Failed to count likes", "error", err)
	}
	return count, err
}
