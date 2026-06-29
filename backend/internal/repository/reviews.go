package repository

import (
	"database/sql"
	"log/slog"
)

func UpsertReview(db *sql.DB, userID string, projectName string, rating int, comment string) error {
	_, err := db.Exec(
		`INSERT INTO project_reviews (user_id, project_name, rating, comment, created_at)
		 VALUES ($1, $2, $3, $4, NOW())
		 ON CONFLICT (user_id, project_name)
		 DO UPDATE SET rating = $3, comment = $4, created_at = NOW()`,
		userID, projectName, rating, comment,
	)
	if err != nil {
		slog.Error("Failed to upsert review", "error", err)
	}
	return err
}

func GetProjectReviews(db *sql.DB, projectName string) (*sql.Rows, error) {
	rows, err := db.Query(`
		SELECT pr.user_id, u.name, u.avatar_url, pr.project_name, pr.rating, pr.comment, pr.created_at
		FROM project_reviews pr JOIN users u ON pr.user_id = u.id
		WHERE pr.project_name = $1
		ORDER BY pr.created_at DESC
	`, projectName)
	if err != nil {
		slog.Error("Failed to query project reviews", "error", err)
	}
	return rows, err
}

func GetProjectRatingSummary(db *sql.DB, projectName string) (float64, int, error) {
	var avg float64
	var count int
	err := db.QueryRow(
		"SELECT COALESCE(AVG(rating), 0), COUNT(*) FROM project_reviews WHERE project_name = $1",
		projectName,
	).Scan(&avg, &count)
	return avg, count, err
}

func GetMarqueeItems(db *sql.DB) (*sql.Rows, error) {
	rows, err := db.Query(`
		SELECT 'review' as type, u.name, u.avatar_url, pr.project_name, pr.rating, pr.comment, pr.created_at
		FROM project_reviews pr JOIN users u ON pr.user_id = u.id
		ORDER BY pr.created_at DESC LIMIT 10
	`)
	if err != nil {
		slog.Error("Failed to query marquee items", "error", err)
	}
	return rows, err
}

func GetUserReview(db *sql.DB, userID, projectName string) (*sql.Row, error) {
	row := db.QueryRow(`
		SELECT pr.user_id, u.name, u.avatar_url, pr.project_name, pr.rating, pr.comment, pr.created_at
		FROM project_reviews pr JOIN users u ON pr.user_id = u.id
		WHERE pr.user_id = $1 AND pr.project_name = $2
	`, userID, projectName)
	return row, nil
}

func DeleteReview(db *sql.DB, userID, projectName string) error {
	_, err := db.Exec("DELETE FROM project_reviews WHERE user_id = $1 AND project_name = $2", userID, projectName)
	if err != nil {
		slog.Error("Failed to delete review", "error", err)
	}
	return err
}

func DeleteReviewByAdmin(db *sql.DB, reviewID int) error {
	_, err := db.Exec("DELETE FROM project_reviews WHERE id = $1", reviewID)
	if err != nil {
		slog.Error("Failed to delete review by admin", "error", err)
	}
	return err
}

func GetAllReviews(db *sql.DB) (*sql.Rows, error) {
	rows, err := db.Query(`
		SELECT pr.id, pr.user_id, u.name, u.avatar_url, pr.project_name, pr.rating, pr.comment, pr.created_at
		FROM project_reviews pr JOIN users u ON pr.user_id = u.id
		ORDER BY pr.created_at DESC LIMIT 100
	`)
	if err != nil {
		slog.Error("Failed to query all reviews", "error", err)
	}
	return rows, err
}
