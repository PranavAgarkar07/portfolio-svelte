package repository

import (
	"database/sql"
	"log/slog"
)

func GetUserByGoogleID(db *sql.DB, googleID string) (*sql.Row, error) {
	row := db.QueryRow("SELECT id, google_id, email, name, avatar_url, role, created_at FROM users WHERE google_id = $1", googleID)
	return row, nil
}

func GetUserByID(db *sql.DB, id string) (*sql.Row, error) {
	row := db.QueryRow("SELECT id, google_id, email, name, avatar_url, role, created_at FROM users WHERE id = $1", id)
	return row, nil
}

func CreateUser(db *sql.DB, googleID, email, name, avatarURL string) (string, error) {
	var id string
	err := db.QueryRow(
		"INSERT INTO users (google_id, email, name, avatar_url) VALUES ($1, $2, $3, $4) RETURNING id",
		googleID, email, name, avatarURL,
	).Scan(&id)
	if err != nil {
		slog.Error("Failed to create user", "error", err)
	}
	return id, err
}

func GetAllUsers(db *sql.DB) (*sql.Rows, error) {
	rows, err := db.Query("SELECT id, google_id, email, name, avatar_url, role, created_at FROM users ORDER BY created_at DESC")
	if err != nil {
		slog.Error("Failed to query users", "error", err)
	}
	return rows, err
}

func UpdateUserRole(db *sql.DB, userID, role string) error {
	_, err := db.Exec("UPDATE users SET role = $1 WHERE id = $2", role, userID)
	if err != nil {
		slog.Error("Failed to update user role", "error", err)
	}
	return err
}

func UpdateUserAvatar(db *sql.DB, userID, avatarURL string) error {
	_, err := db.Exec("UPDATE users SET avatar_url = $1 WHERE id = $2", avatarURL, userID)
	if err != nil {
		slog.Error("Failed to update user avatar", "error", err)
	}
	return err
}

func GetUserByEmail(db *sql.DB, email string) (*sql.Row, error) {
	row := db.QueryRow("SELECT id, google_id, email, name, avatar_url, role, created_at FROM users WHERE email = $1", email)
	return row, nil
}
