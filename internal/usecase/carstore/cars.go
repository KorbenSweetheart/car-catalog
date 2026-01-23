package carstore

import (
	"context"
	"log/slog"

	"gitea.kood.tech/ivanandreev/viewer/internal/domain"
	"gitea.kood.tech/ivanandreev/viewer/internal/lib/e"
)

// The Business Logic -> provide car/cars

const limit = 4

type CarProvider interface {
	Car(ctx context.Context, ID int) (domain.Car, error)
	Cars(ctx context.Context) ([]domain.Car, error)
	RandomCars(ctx context.Context, limit int) ([]domain.Car, error)
	Metadata(ctx context.Context) (domain.Metadata, error)
}

type CacheProvider interface {
	Get(ctx context.Context, id int) (domain.Car, bool)
	Set(ctx context.Context, c domain.Car)
	GetMetadata(ctx context.Context) (domain.Metadata, bool)
	SetMetadata(ctx context.Context, m domain.Metadata)
}

type CarStore struct {
	log   *slog.Logger
	repo  CarProvider
	cache CacheProvider
}

func New(log *slog.Logger, r CarProvider, c CacheProvider) *CarStore {
	return &CarStore{
		log:   log,
		repo:  r,
		cache: c,
	}
}

func (s *CarStore) Car(ctx context.Context, ID int) (domain.Car, error) {
	const op = "usecase.carstore.Car"

	log := s.log.With(
		slog.String("op", op),
	)

	if car, found := s.cache.Get(ctx, ID); found {
		s.log.Debug("loaded car from cache", "id", ID)
		return car, nil
	}

	car, err := s.repo.Car(ctx, ID)
	if err != nil {
		log.Error("failed to get car by id", slog.Any("error", err))
		return domain.Car{}, e.Wrap("failed to get car by id: %w", err)
	}

	s.cache.Set(ctx, car)

	log.Info("car loaded",
		slog.Int("car_id", ID),
	)

	return car, nil
}

func (s *CarStore) Cars(ctx context.Context) ([]domain.Car, error) {
	const op = "usecase.carstore.Cars"

	log := s.log.With(
		slog.String("op", op),
	)

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

	log := s.log.With(
		slog.String("op", op),
	)

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

func (s *CarStore) Filters(ctx context.Context) (domain.Metadata, error) {
	const op = "usecase.carstore.Filters"

	log := s.log.With(
		slog.String("op", op),
	)

	if meta, found := s.cache.GetMetadata(ctx); found {
		s.log.Debug("loaded metadata from cache")
		return meta, nil
	}

	filters, err := s.repo.Metadata(ctx)
	if err != nil {
		log.Error("failed to get metadata", slog.Any("error", err))
		return domain.Metadata{}, e.Wrap("failed to get metadata: %w", err)
	}

	s.cache.SetMetadata(ctx, filters)

	log.Info("metadata loaded")

	return filters, nil
}

/* Catalog filters
Manufacturer
Category
Drive Train
Transmission
Year min (user input)
Power (hp) (user input)
*/
