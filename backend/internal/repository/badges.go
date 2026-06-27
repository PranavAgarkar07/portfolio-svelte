package repository

import (
	"database/sql"
	"log/slog"
	"time"
)

type BadgeRow struct {
	ID            int
	Name          string
	ImageURL      string
	CredentialURL string
	Rarity        string
	Category      string
	Important     bool
	DisplayOrder  int
	CreatedAt     time.Time
}

func GetBadges(db *sql.DB, limit, offset int) (*sql.Rows, error) {
	rows, err := db.Query(
		"SELECT id, name, image_url, credential_url, rarity, category, important, display_order, created_at FROM badges ORDER BY important DESC, display_order ASC, id DESC LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		slog.Error("Failed to query badges", "error", err)
	}
	return rows, err
}

func InsertBadge(db *sql.DB, name, imageURL, credentialURL, rarity, category string, important bool, displayOrder int) (int, error) {
	var id int
	err := db.QueryRow(
		"INSERT INTO badges (name, image_url, credential_url, rarity, category, important, display_order) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		name, imageURL, credentialURL, rarity, category, important, displayOrder,
	).Scan(&id)
	if err != nil {
		slog.Error("Failed to insert badge", "error", err)
	}
	return id, err
}

func UpdateBadge(db *sql.DB, id int, name, imageURL, credentialURL, rarity, category string, important bool, displayOrder int) (int64, error) {
	result, err := db.Exec(
		"UPDATE badges SET name=$1, image_url=$2, credential_url=$3, rarity=$4, category=$5, important=$6, display_order=$7 WHERE id=$8",
		name, imageURL, credentialURL, rarity, category, important, displayOrder, id,
	)
	if err != nil {
		slog.Error("Failed to update badge", "error", err)
		return 0, err
	}
	rowsAffected, raErr := result.RowsAffected()
	if raErr != nil {
		slog.Error("Failed to get rows affected for badge update", "error", raErr)
	}
	return rowsAffected, nil
}

func DeleteBadge(db *sql.DB, id int) (int64, error) {
	result, err := db.Exec("DELETE FROM badges WHERE id=$1", id)
	if err != nil {
		slog.Error("Failed to delete badge", "error", err)
		return 0, err
	}
	rowsAffected, raErr := result.RowsAffected()
	if raErr != nil {
		slog.Error("Failed to get rows affected for badge delete", "error", raErr)
	}
	return rowsAffected, nil
}

func ReorderBadges(db *sql.DB, order []int) error {
	ids := make([]int, len(order))
	positions := make([]int, len(order))
	for i, id := range order {
		ids[i] = id
		positions[i] = i
	}
	_, err := db.Exec(
		`UPDATE badges SET display_order = data.position
		 FROM (SELECT UNNEST($1::int[]) AS id, UNNEST($2::int[]) AS position) AS data
		 WHERE badges.id = data.id`,
		ids, positions,
	)
	if err != nil {
		slog.Error("Failed to reorder badges", "error", err)
	}
	return err
}
