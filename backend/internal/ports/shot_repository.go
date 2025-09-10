package ports

import (
	"context"
	"nice-shot/backend/internal/domain"
)

type ShotRepository interface {
	Save(ctx context.Context, s domain.Shot) (domain.Shot, error)
	FindByID(ctx context.Context, id string) (domain.Shot, error)
	List(ctx context.Context, limit int) ([]domain.Shot, error)
	Delete(ctx context.Context, id string) error
}
