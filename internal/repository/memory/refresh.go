package memory

import (
	"context"
	"log/slog"
	"strings"
	"sync"
	"time"

	"gitea.kood.tech/ivanandreev/viewer/internal/domain"
	"gitea.kood.tech/ivanandreev/viewer/internal/lib/e"
)

// Refresh is the "Brain" of the repository.
// It fetches raw data, joins it together, and updates the memory safely.
func (r *Repository) Refresh(ctx context.Context) error {
	const op = "repository.memory.Refresh"

	log := r.log.With(
		slog.String("op", op),
	)

	log.Info("refresh started")

	t1 := time.Now()

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
		log.Error("refresh failed on step: fetch_cars", slog.Any("error", errCars))
		return e.Wrap("failed to fetch cars", errCars)
	}
	if errManufs != nil {
		log.Error("refresh failed on step: fetch_manufacturers", slog.Any("error", errManufs))
		return e.Wrap("failed to fetch manufacturers", errManufs)
	}
	if errCats != nil {
		log.Error("refresh failed on step: fetch_categories", slog.Any("error", errCats))
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

	log.Info("refresh completed",
		slog.Int("cars_count", len(newSlice)),
		slog.String("duration", time.Since(t1).String()),
	)

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
