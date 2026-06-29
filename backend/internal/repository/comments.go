package repository

import (
	"database/sql"
	"log/slog"
)

const commentCols = "bc.id, bc.post_id, bc.parent_id, bc.user_id, u.name, COALESCE(u.avatar_url, ''), bc.content, bc.created_at, bc.updated_at"
const commentJoin = "FROM blog_comments bc JOIN users u ON bc.user_id = u.id"

func GetCommentsByPostID(db *sql.DB, postID string) (*sql.Rows, error) {
	rows, err := db.Query(`
		SELECT `+commentCols+`
		`+commentJoin+`
		WHERE bc.post_id = $1
		ORDER BY bc.created_at ASC
	`, postID)
	if err != nil {
		slog.Error("Failed to query comments", "error", err)
	}
	return rows, err
}

func CreateComment(db *sql.DB, postID, userID, content string, parentID *string) (string, error) {
	var id string
	err := db.QueryRow(
		`INSERT INTO blog_comments (post_id, user_id, content, parent_id)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id`,
		postID, userID, content, parentID,
	).Scan(&id)
	if err != nil {
		slog.Error("Failed to create comment", "error", err)
	}
	return id, err
}

func GetCommentByID(db *sql.DB, id string) (*sql.Row, error) {
	row := db.QueryRow(`
		SELECT `+commentCols+`
		`+commentJoin+`
		WHERE bc.id = $1
	`, id)
	return row, nil
}

func DeleteComment(db *sql.DB, id string) error {
	_, err := db.Exec("DELETE FROM blog_comments WHERE id = $1", id)
	if err != nil {
		slog.Error("Failed to delete comment", "error", err)
	}
	return err
}

func DeleteCommentTree(db *sql.DB, rootID string) error {
	_, err := db.Exec(`
		WITH RECURSIVE subtree AS (
			SELECT id FROM blog_comments WHERE id = $1
			UNION ALL
			SELECT bc.id FROM blog_comments bc JOIN subtree s ON bc.parent_id = s.id
		)
		DELETE FROM blog_comments WHERE id IN (SELECT id FROM subtree)
	`, rootID)
	if err != nil {
		slog.Error("Failed to delete comment tree", "error", err)
	}
	return err
}
