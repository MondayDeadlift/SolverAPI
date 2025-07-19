package service

import (
	"SolverAPI/internal/model"
	"SolverAPI/internal/repository"
	"SolverAPI/pkg/codewars"
	"context"
	"time"
)

type UserService struct {
	repo repository.UserRepository
	cw   *codewars.Client
}

func NewUserService(repo repository.UserRepository, cw *codewars.Client) *UserService {
	return &UserService{repo: repo, cw: cw}
}

func (s *UserService) SyncUser(ctx context.Context, username string) (*model.User, error) {
	//Получаем данные из Codewars API
	cwUser, err := s.cw.GetUser(username)
	if err != nil {
		return nil, err
	}

	//Преобразуем в нашу модель
	user := &model.User{
		CodewarsUser: *cwUser,
		CreatedAt:    time.Now(),
	}

	//Сохраняем в БД
	if err := s.repo.CreateOrUpdateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
