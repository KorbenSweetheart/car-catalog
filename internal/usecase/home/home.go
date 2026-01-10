package home

import (
	"context"
	"log/slog"
	"viewer/internal/domain"
	"viewer/internal/lib/e"
)

const limit = 4

type CarRepository interface {
	GetRandom(ctx context.Context, limit int) ([]domain.Car, error)
}

type Usecase struct {
	log  *slog.Logger
	repo CarRepository
}

// 3. The Constructor
// When called from app, it looks like: home.New(log, repo)
func New(log *slog.Logger, r CarRepository) *Usecase {
	return &Usecase{
		log:  log,
		repo: r,
	}
}

// The Business Logic
func (u *Usecase) Load(ctx context.Context) ([]domain.Car, error) {
	const op = "usecase.home.Load"

	log := u.log.With("op", op)

	// Logic: Get 4 cars, maybe filter them, maybe handle errors specifically
	cars, err := u.repo.GetRandom(ctx, limit)
	if err != nil {
		log.Error("failed to get random cars", slog.Any("error", err))
		return nil, e.Wrap("failed to get random cars: %w", err)
	}

	log.Info("home page loaded",
		slog.Int("cars_count", len(cars)),
	)

	return cars, nil
}
