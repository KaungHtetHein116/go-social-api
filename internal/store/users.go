package store

import (
	"context"
	"database/sql"
)

type User struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"-"`
	Created_At string `json:"created_at"`
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (username, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	if err := s.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Password,
	).Scan(&user.ID, &user.Created_At); err != nil {
		return err
	}

	return nil
}

func (s *UserStore) GetByUserID(ctx context.Context, userID int) (*User, error) {
	query := `
		SELECT id, username, email, created_at
		FROM users
		WHERE id = $1
	`

	user := &User{}
	if err := s.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Created_At,
	); err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrRecordNotFound

		default:
			return nil, err
		}
	}

	return user, nil
}
