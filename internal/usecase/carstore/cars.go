package carstore

import (
	"context"
	"log/slog"
	"viewer/internal/domain"
	"viewer/internal/lib/e"
)

// The Business Logic -> provide car/cars

const limit = 4

type CarProvider interface {
	Car(ctx context.Context, ID int) (domain.Car, error)
	Cars(ctx context.Context) ([]domain.Car, error)
	Manufacturer(ctx context.Context, ID int) (domain.Manufacturer, error)
	Manufacturers(ctx context.Context) ([]domain.Manufacturer, error)
	Category(ctx context.Context, ID int) (domain.Category, error)
	Categories(ctx context.Context) ([]domain.Category, error)
	RandomCars(ctx context.Context, limit int) ([]domain.Car, error)
	Metadata(ctx context.Context) (domain.Metadata, error)
}

type CarStore struct {
	log  *slog.Logger
	repo CarProvider
}

func New(log *slog.Logger, r CarProvider) *CarStore {
	return &CarStore{
		log:  log,
		repo: r,
	}
}

func (s *CarStore) Car(ctx context.Context, ID int) (domain.Car, error) {
	const op = "usecase.carstore.Car"

	log := s.log.With("op", op)

	car, err := s.repo.Car(ctx, ID)
	if err != nil {
		log.Error("failed to get car by id", slog.Any("error", err))
		return domain.Car{}, e.Wrap("failed to get car by id: %w", err)
	}

	log.Info("car loaded",
		slog.Int("car_id", ID),
	)

	return car, nil
}

func (s *CarStore) Cars(ctx context.Context) ([]domain.Car, error) {
	const op = "usecase.carstore.Cars"

	log := s.log.With("op", op)

	cars, err := s.repo.Cars(ctx)
	if err != nil {
		log.Error("failed to get cars catalog", slog.Any("error", err))
		return nil, e.Wrap("failed to get cars catalog: %w", err)
	}

	log.Info("catalog loaded",
		slog.Int("cars_count", len(cars)),
	)

	return cars, nil
}

func (s *CarStore) RandomCars(ctx context.Context) ([]domain.Car, error) {
	const op = "usecase.carstore.RandomCars"

	log := s.log.With("op", op)

	// Logic: Get 4 cars, maybe filter them, maybe handle errors specifically
	cars, err := s.repo.RandomCars(ctx, limit)
	if err != nil {
		log.Error("failed to get random cars", slog.Any("error", err))
		return nil, e.Wrap("failed to get random cars: %w", err)
	}

	log.Info("random cars loaded",
		slog.Int("cars_count", len(cars)),
	)

	return cars, nil
}

/* Catalog filters
Manufacturer
Category
Drive Train
Transmission
Year min (user input)
Power (hp) (user input)
*/

func (s *CarStore) Filters(ctx context.Context) (domain.Metadata, error) {
	const op = "usecase.carstore.Filters"

	log := s.log.With("op", op)

	filters, err := s.repo.Metadata(ctx)
	if err != nil {
		log.Error("failed to get metadata", slog.Any("error", err))
		return domain.Metadata{}, e.Wrap("failed to get metadata: %w", err)
	}

	return filters, nil
}
