package store

import (
	"context"
	"database/sql"
)

type Comment struct {
	ID        int    `json:"id"`
	PostID    int    `json:"post_id"`
	UserID    int    `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	User      User   `json:"user"`
}

type CommentStore struct {
	db *sql.DB
}

func (c *CommentStore) GetByPostID(ctx context.Context, postID int) ([]Comment, error) {
	query := `
            SELECT
                c.id, c.post_id, c.user_id, c.content, c.created_at, users.username, users.id
            FROM
                comments c
                JOIN users ON c.user_id = users.id
            WHERE
                post_id = $1
            ORDER BY
                c.created_at DESC
            `

	rows, err := c.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	comments := make([]Comment, 0)
	for rows.Next() {
		var c Comment
		if err := rows.Scan(
			&c.ID,
			&c.PostID,
			&c.UserID,
			&c.Content,
			&c.CreatedAt,
			&c.User.Username,
			&c.User.ID,
		); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
