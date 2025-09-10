package application

import (
	"context"
	"time"

	"nice-shot/backend/internal/domain"
	"nice-shot/backend/internal/ports"
)

type ShotService struct {
	repository ports.ShotRepository
	now        func() time.Time
}

func NewShotService(repository ports.ShotRepository, now func() time.Time) *ShotService {
	return &ShotService{repository: repository, now: now}
}

func (s *ShotService) Create(ctx context.Context, shot domain.Shot) (domain.Shot, error) {
	// In a real app, validate here and set CreatedAt/UpdatedAt if needed
	return s.repository.Save(ctx, shot)
}

func (s *ShotService) Get(ctx context.Context, id string) (domain.Shot, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *ShotService) List(ctx context.Context, limit int) ([]domain.Shot, error) {
	if limit <= 0 {
		limit = 100
	}
	return s.repository.List(ctx, limit)
}

func (s *ShotService) Delete(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, id)
}
