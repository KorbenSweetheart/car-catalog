package usecase

import (
	"context"
	"viewer/internal/domain"
)

type Reporistory interface {
	GetByID(ctx context.Context, id int) (domain.Car, error)
	GetMetadata(ctx context.Context) (domain.Metadata, error)
	// GetRandom(ctx context.Context, limit int) ([]domain.Car, error)
	Refresh(ctx context.Context) error
	Search(ctx context.Context, f domain.FilterOptions) ([]domain.Car, error)
}

type UseCase struct {
	repo Reporistory
}

func NewUseCase(r Reporistory) *UseCase {
	return &UseCase{
		repo: r,
	}
}
