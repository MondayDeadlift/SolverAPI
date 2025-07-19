package postgres

import (
	"SolverAPI/internal/model"
	"SolverAPI/internal/repository"
	"context"
	"database/sql"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateOrUpdateUser(ctx context.Context, user *model.User) error {
	query := `
        INSERT INTO users (username, honor, created_at, updated_at)
        VALUES ($1, $2, NOW(), NOW())
        ON CONFLICT (username) DO UPDATE
        SET honor = $2, updated_at = NOW()
    `
	_, err := r.db.ExecContext(ctx, query, user.Username, user.Honor)
	return err
}

func (r *UserRepo) GetUser(ctx context.Context, username string) (*model.User, error) {
	query := `SELECT username, honor, created_at FROM users WHERE username = $1`
	row := r.db.QueryRowContext(ctx, query, username)

	var user model.User
	err := row.Scan(&user.Username, &user.Honor, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}
