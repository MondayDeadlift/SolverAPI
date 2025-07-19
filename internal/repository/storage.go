package repository

import (
	"SolverAPI/internal/model"
	"context"
	"errors"
)

// UserRepository определяет контракт для работы с пользователями
type UserRepository interface {
	CreateOrUpdateUser(ctx context.Context, user *model.User) error
	GetUser(ctx context.Context, username string) (*model.User, error)
}

type KataRepository interface {
	SaveKata(ctx context.Context, kata *model.Kata) error
	GetRandomKata(ctx context.Context) (*model.Kata, error)
}

var (
	ErrUserNotFound = errors.New("user not found")
)
