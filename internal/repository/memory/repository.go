package memory

import (
	"context"
	"errors"
	"math/rand"
	"strings"
	"sync"
	"time"
	"viewer/internal/domain"
	"viewer/internal/lib/e"
)

type DataFetcher interface {
	FetchCars(ctx context.Context) ([]domain.Car, error)
	FetchManufacturers(ctx context.Context) ([]domain.Manufacturer, error)
	FetchCategories(ctx context.Context) ([]domain.Category, error)
}

type Repository struct {
	fetcher DataFetcher

	// 1. Hybrid Storage for Cars
	carsSlice []*domain.Car       // For iteration (Search, Random)
	carsMap   map[int]*domain.Car // For fast lookup (GetByID)

	// 2. Metadata (for Sidebar Filters)
	manufacturers []domain.Manufacturer
	categories    []domain.Category

	// NEW: Lists for Dropdowns
	drivetrains   []string
	transmissions []string

	mu sync.RWMutex
}

func New(fetcher DataFetcher) *Repository {
	return &Repository{
		fetcher:   fetcher,
		carsSlice: make([]*domain.Car, 0),
		carsMap:   make(map[int]*domain.Car),
	}
}

// -------------------------------------------------------------------------
// Lifecycle Methods
// -------------------------------------------------------------------------

// Refresh is the "Brain" of the repository.
// It fetches raw data, joins it together, and updates the memory safely.
func (r *Repository) Refresh(ctx context.Context) error {
	// 1. Define variables to hold results and errors
	var (
		rawCars []domain.Car
		vendors []domain.Manufacturer
		cats    []domain.Category

		errCars, errManufs, errCats error

		wg sync.WaitGroup
	)

	// 2. Launch Goroutines
	// We are making 3 independent calls. Running them in parallel
	// reduces the total time to the duration of the slowest request.
	wg.Add(3)

	go func() {
		defer wg.Done()
		rawCars, errCars = r.fetcher.FetchCars(ctx)
	}()

	go func() {
		defer wg.Done()
		vendors, errManufs = r.fetcher.FetchManufacturers(ctx)
	}()

	go func() {
		defer wg.Done()
		cats, errCats = r.fetcher.FetchCategories(ctx)
	}()

	// 3. Wait for all to finish
	wg.Wait()

	// 4. Check Errors
	// If any single part fails, we abort the update to ensure consistency.
	// We wrap the specific error to know which API call failed.
	if errCars != nil {
		return e.Wrap("failed to fetch cars", errCars)
	}
	if errManufs != nil {
		return e.Wrap("failed to fetch manufacturers", errManufs)
	}
	if errCats != nil {
		return e.Wrap("failed to fetch categories", errCats)
	}

	// 2. Build Lookup Maps for Enrichment (The "Join" Preparation)
	manufMap := make(map[int]domain.Manufacturer)
	for _, m := range vendors {
		manufMap[m.ID] = m
	}

	catMap := make(map[int]domain.Category)
	for _, c := range cats {
		catMap[c.ID] = c
	}

	// Sets to deduplicate strings for metadata
	uniqueDrivetrains := make(map[string]bool)
	uniqueTransmissions := make(map[string]bool)

	// 3. Prepare New Storage Structures
	// We use pointers *domain.Car to avoid copying heavy structs during iteration
	newSlice := make([]*domain.Car, 0, len(rawCars))
	newMap := make(map[int]*domain.Car, len(rawCars))

	// 4. The Enrichment Loop
	for i := range rawCars {
		// We modify the car directly in the rawCars slice
		car := &rawCars[i]

		// 2. Prepopulate transmission type based on the gearbox
		car.Specs.Transmission = normalizeGearbox(car.Specs.Gearbox)

		// ENRICHMENT: Fill in the blank Manufacturer details using our map
		if fullManuf, ok := manufMap[car.Manufacturer.ID]; ok {
			car.Manufacturer = fullManuf
		}

		// ENRICHMENT: Fill in the blank Category details
		if fullCat, ok := catMap[car.Category.ID]; ok {
			car.Category = fullCat
		}

		// Collect Metadata Strings
		if car.Specs.Drivetrain != "" {
			uniqueDrivetrains[car.Specs.Drivetrain] = true
		}
		uniqueTransmissions[car.Specs.Transmission] = true

		// Add to our new structures
		newSlice = append(newSlice, car)
		newMap[car.ID] = car
	}

	// 8. Convert Maps to Slices
	drivetrains := make([]string, 0, len(uniqueDrivetrains))
	for d := range uniqueDrivetrains {
		drivetrains = append(drivetrains, d)
	}

	transmissions := make([]string, 0, len(uniqueTransmissions))
	for t := range uniqueTransmissions {
		transmissions = append(transmissions, t)
	}

	// 5. Atomic Swap (Thread-Safe Update)
	r.mu.Lock()
	r.carsSlice = newSlice
	r.carsMap = newMap
	r.manufacturers = vendors
	r.categories = cats
	r.drivetrains = drivetrains
	r.transmissions = transmissions
	r.mu.Unlock()

	return nil
}

// Helper function to prepopulate transmission field
func normalizeGearbox(gearbox string) string {

	// If it contains "manual", it's a manual
	if strings.Contains(strings.ToLower(gearbox), "manual") {
		return domain.TransmissionManual
	}

	// Everything else (Automatic, CVT, DSG, Dual Clutch, Single-Speed)
	// counts as "Automatic" for a general filter.
	return domain.TransmissionAutomatic
}

// -------------------------------------------------------------------------
// Access Methods (The "Read" Side)
// -------------------------------------------------------------------------

func (r *Repository) GetByID(ctx context.Context, id int) (domain.Car, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	car, exists := r.carsMap[id]
	if !exists {
		return domain.Car{}, errors.New("car not found")
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

// Search handles filtering for the catalog
func (r *Repository) Search(ctx context.Context, f domain.FilterOptions) ([]domain.Car, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Heuristic: Allocate for ~20% of cars to avoid resizing
	results := make([]domain.Car, 0, len(r.carsSlice)/5)

	for _, car := range r.carsSlice {
		// 1. Manufacturer (Exact Match)
		if f.ManufacturerID > 0 && car.Manufacturer.ID != f.ManufacturerID {
			continue
		}

		// 2. Category (Exact Match)
		if f.CategoryID > 0 && car.Category.ID != f.CategoryID {
			continue
		}

		// 3. Min Year (Range)
		if f.MinYear > 0 && car.Year < f.MinYear {
			continue
		}

		// 4. Min HP (Range)
		if f.MinHP > 0 && car.Specs.HP < f.MinHP {
			continue
		}

		// We compare the filter (e.g., "Automatic") against our clean field
		if f.Transmission != "" && car.Specs.Transmission != f.Transmission {
			continue
		}

		// 6. Drivetrain (Exact Match)
		if f.Drivetrain != "" && car.Specs.Drivetrain != f.Drivetrain {
			continue
		}

		results = append(results, *car)
	}

	return results, nil
}

// GetMetadata returns the lists needed for the Sidebar
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
