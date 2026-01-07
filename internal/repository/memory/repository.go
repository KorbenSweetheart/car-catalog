package memory

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"sync"
	"time"
	"viewer/internal/domain"
	"viewer/internal/lib/e"
)

type DataFetcher interface {
	FetchCars(ctx context.Context) ([]domain.Car, error)
	FetchCategories(ctx context.Context) ([]domain.Category, error)
	FetchManufacturers(ctx context.Context) ([]domain.Manufacturer, error)
}

type Repository struct {
	log     *slog.Logger
	fetcher DataFetcher

	// Hybrid Storage for Cars
	carsSlice []*domain.Car
	carsMap   map[int]*domain.Car

	// Metadata (for Sidebar Filters)
	manufacturers []domain.Manufacturer
	categories    []domain.Category
	drivetrains   []string
	transmissions []string

	mu sync.RWMutex
}

func New(log *slog.Logger, fetcher DataFetcher) *Repository {
	return &Repository{
		log:       log,
		fetcher:   fetcher,
		carsSlice: make([]*domain.Car, 0),
		carsMap:   make(map[int]*domain.Car),
	}
}

// -------------------------------------------------------------------------
// Access Methods (The "Read" Side)
// -------------------------------------------------------------------------

func (r *Repository) GetByID(ctx context.Context, id int) (domain.Car, error) {

	r.mu.RLock()
	defer r.mu.RUnlock()

	car, exists := r.carsMap[id]
	if !exists {
		return domain.Car{}, e.Wrap("car not found", fmt.Errorf("id: %d", id))
	}
	// Return a copy (value) so the caller can't mutate our memory
	return *car, nil
}

func (r *Repository) GetRandom(ctx context.Context, limit int) ([]domain.Car, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := len(r.carsSlice)
	if count == 0 {
		return []domain.Car{}, nil
	}

	if count <= limit {
		// Return everything we have
		result := make([]domain.Car, 0, count)
		for _, c := range r.carsSlice {
			result = append(result, *c)
		}
		return result, nil
	}

	// Random selection
	// Note: For production, consider using crypto/rand or a seeded source
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	perm := rng.Perm(count)

	result := make([]domain.Car, 0, limit)
	for i := 0; i < limit; i++ {
		idx := perm[i]
		result = append(result, *r.carsSlice[idx])
	}

	return result, nil
}

// GetMetadata returns the lists needed for filters
func (r *Repository) GetMetadata(ctx context.Context) (domain.Metadata, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// NOTE: if you ever marshal this to JSON, it might come out as null instead of [].
	// Optional: Ensure non-nil slices for JSON friendliness
	// manufs := r.manufacturers
	// if manufs == nil { manufs = []domain.Manufacturer{} }

	return domain.Metadata{
		Manufacturers: append([]domain.Manufacturer(nil), r.manufacturers...),
		Categories:    append([]domain.Category(nil), r.categories...),
		Drivetrains:   append([]string(nil), r.drivetrains...),
		Transmissions: append([]string(nil), r.transmissions...),
	}, nil
}
