package service

import (
	"SolverAPI/internal/model"
	"SolverAPI/internal/repository"
	"SolverAPI/pkg/codewars"
	"fmt"

	"context"
	"time"
)

type KataService struct {
	repo     repository.KataRepository
	cwClient *codewars.Client
}

func NewKataService(repo repository.KataRepository, cwClient *codewars.Client) *KataService {
	return &KataService{
		repo:     repo,
		cwClient: cwClient,
	}
}

func (s *KataService) GetRandomKata(ctx context.Context) (*model.Kata, error) {
	//Обновляем буфер при необходимости
	s.cwClient.RefreshBuffer()

	//Получаем случайный ID из буфера
	randID, err := s.cwClient.GetRandomKataID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get random ID: %w", err)
	}

	//Получаем полные данные по задаче
	cwKata, err := s.cwClient.GetKataByID(ctx, randID)
	if err != nil {
		return nil, fmt.Errorf("failed to get kata details: %w", err)
	}

	//Сохраняем в БД
	kata := &model.Kata{
		CodewarsKata: *cwKata,
		AddedAt:      time.Now(),
	}

	if err := s.repo.SaveKata(ctx, kata); err != nil {
		return nil, fmt.Errorf("failed to save kata: %w", err)
	}

	return kata, nil
}
