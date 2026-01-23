package adapter

import (
	"context"
	"fmt"
	"log/slog"

	"gitea.kood.tech/ivanandreev/viewer/internal/domain"
	"gitea.kood.tech/ivanandreev/viewer/pkg/cache"
)

// Adapter wraps the generic cache to make it type-safe for the domain
type CacheAdapter struct {
	cache *cache.Cache
	log   *slog.Logger
}

func NewAdapter(c *cache.Cache, logger *slog.Logger) *CacheAdapter {
	return &CacheAdapter{cache: c, log: logger}
}

func (a *CacheAdapter) Get(ctx context.Context, id int) (domain.Car, bool) {
	const op = "repository.adapter.Get"

	log := a.log.With(
		slog.String("op", op),
	)

	key := fmt.Sprintf("car:%d", id)
	val, found := a.cache.Get(key)
	if !found {
		log.Debug("car not found in cache",
			slog.Int("car_id", id),
		)
		return domain.Car{}, false
	}

	// Type assertion: Unbox the interface{} back to a Car struct
	car, ok := val.(domain.Car)
	if !ok {
		log.Error("failed to assert cached object to car", slog.Any("error", fmt.Errorf("failed to assert object:%v to car", val)))
		return domain.Car{}, false
	}

	log.Debug("car loaded from cache",
		slog.Int("car_id", id),
	)

	return car, true
}

func (a *CacheAdapter) Set(ctx context.Context, car domain.Car) {
	const op = "repository.adapter.Set"

	log := a.log.With(
		slog.String("op", op),
	)

	key := fmt.Sprintf("car:%d", car.ID)
	a.cache.Set(key, car, cache.DefaultExpiration)

	log.Debug("car added to cache",
		slog.Int("car_id", car.ID),
	)
}

func (a *CacheAdapter) GetMetadata(ctx context.Context) (domain.Metadata, bool) {
	const op = "repository.adapter.Get"

	log := a.log.With(
		slog.String("op", op),
	)

	val, found := a.cache.Get("metadata")
	if !found {
		log.Debug("metadata not found in cache")
		return domain.Metadata{}, false
	}
	meta, ok := val.(domain.Metadata)
	if !ok {
		log.Error("failed to assert cached object to metadata", slog.Any("error", fmt.Errorf("failed to assert object:%v to metadata", val)))
		return domain.Metadata{}, false
	}

	return meta, ok
}

func (a *CacheAdapter) SetMetadata(ctx context.Context, m domain.Metadata) {
	const op = "repository.adapter.SetMetadata"

	log := a.log.With(
		slog.String("op", op),
	)

	// Maybe use a longer TTL??? default is 10 min
	a.cache.Set("metadata", m, cache.DefaultExpiration)

	log.Debug("metadata added to cache")
}
