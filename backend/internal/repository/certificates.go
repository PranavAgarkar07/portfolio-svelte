package repository

import (
	"database/sql"
	"log/slog"
	"time"
)

func InsertContactMessage(db *sql.DB, name, email, topic, message string) error {
	_, err := db.Exec(
		"INSERT INTO contact_messages (name, email, topic, message) VALUES ($1, $2, $3, $4)",
		name, email, topic, message,
	)
	if err != nil {
		slog.Error("Failed to insert contact message", "error", err)
	}
	return err
}

func GetMessages(db *sql.DB, limit, offset int) (*sql.Rows, error) {
	rows, err := db.Query(
		"SELECT id, name, email, topic, message, is_read, created_at FROM contact_messages ORDER BY created_at DESC LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		slog.Error("Failed to query contact messages", "error", err)
	}
	return rows, err
}

func MarkMessageRead(db *sql.DB, id int) (int64, error) {
	result, err := db.Exec("UPDATE contact_messages SET is_read = TRUE WHERE id = $1", id)
	if err != nil {
		slog.Error("Failed to mark message as read", "error", err)
		return 0, err
	}
	rowsAffected, raErr := result.RowsAffected()
	if raErr != nil {
		slog.Error("Failed to get rows affected for mark-read", "error", raErr)
	}
	return rowsAffected, nil
}

type CertificateRow struct {
	ID            int
	Title         string
	Issuer        string
	Date          string
	CredentialURL string
	ImageURL      string
	Tags          []string
	IsVerified    bool
	DisplayOrder  int
	CreatedAt     time.Time
}

func GetCertificates(db *sql.DB, limit, offset int) (*sql.Rows, error) {
	rows, err := db.Query(
		"SELECT id, title, issuer, date, credential_url, image_url, tags, is_verified, display_order, created_at FROM certificates ORDER BY date DESC NULLS LAST, display_order ASC, id DESC LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		slog.Error("Failed to query certificates", "error", err)
	}
	return rows, err
}

func InsertCertificate(db *sql.DB, title, issuer, date, credentialURL, imageURL string, tags []string, isVerified bool, displayOrder int) (int, error) {
	var id int
	err := db.QueryRow(
		"INSERT INTO certificates (title, issuer, date, credential_url, image_url, tags, is_verified, display_order) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id",
		title, issuer, date, credentialURL, imageURL, tags, isVerified, displayOrder,
	).Scan(&id)
	if err != nil {
		slog.Error("Failed to insert certificate", "error", err)
	}
	return id, err
}

func UpdateCertificate(db *sql.DB, id int, title, issuer, date, credentialURL, imageURL string, tags []string, isVerified bool, displayOrder int) (int64, error) {
	result, err := db.Exec(
		"UPDATE certificates SET title=$1, issuer=$2, date=$3, credential_url=$4, image_url=$5, tags=$6, is_verified=$7, display_order=$8 WHERE id=$9",
		title, issuer, date, credentialURL, imageURL, tags, isVerified, displayOrder, id,
	)
	if err != nil {
		slog.Error("Failed to update certificate", "error", err)
		return 0, err
	}
	rowsAffected, raErr := result.RowsAffected()
	if raErr != nil {
		slog.Error("Failed to get rows affected for cert update", "error", raErr)
	}
	return rowsAffected, nil
}

func DeleteCertificate(db *sql.DB, id int) (int64, error) {
	result, err := db.Exec("DELETE FROM certificates WHERE id=$1", id)
	if err != nil {
		slog.Error("Failed to delete certificate", "error", err)
		return 0, err
	}
	rowsAffected, raErr := result.RowsAffected()
	if raErr != nil {
		slog.Error("Failed to get rows affected for cert delete", "error", raErr)
	}
	return rowsAffected, nil
}

func ReorderCertificates(db *sql.DB, order []int) error {
	ids := make([]int, len(order))
	positions := make([]int, len(order))
	for i, id := range order {
		ids[i] = id
		positions[i] = i
	}
	_, err := db.Exec(
		`UPDATE certificates SET display_order = data.position
		 FROM (SELECT UNNEST($1::int[]) AS id, UNNEST($2::int[]) AS position) AS data
		 WHERE certificates.id = data.id`,
		ids, positions,
	)
	if err != nil {
		slog.Error("Failed to reorder certificates", "error", err)
	}
	return err
}
