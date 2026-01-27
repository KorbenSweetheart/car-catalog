package carstore

import (
	"context"
	"log/slog"

	"gitea.kood.tech/ivanandreev/viewer/internal/domain"
	"gitea.kood.tech/ivanandreev/viewer/internal/lib/e"
)

// Filtering for the catalog
// get domain.FilterOptions, and return []domain.Car
// filtering should happen here on the local slice of cars, based on the Filters.
func (s *CarStore) Catalog(ctx context.Context, filters domain.FilterOptions) ([]domain.Car, error) {
	const op = "usecase.carstore.Catalog"

	log := s.log.With(
		slog.String("op", op),
	)

	allCars, err := s.repo.Cars(ctx)
	if err != nil {
		log.Error("failed to get cars catalog", slog.Any("error", err))
		return nil, e.Wrap("failed to get cars catalog: %w", err)
	}

	// Apply "In-Memory" Filtering
	finalList := s.filterCars(allCars, filters)

	return finalList, nil
}

func (s *CarStore) filterCars(allCars []domain.Car, f domain.FilterOptions) []domain.Car {
	var filtered []domain.Car

	for _, car := range allCars {
		// 1. Manufacturer
		if f.ManufacturerID != 0 && car.Manufacturer.ID != f.ManufacturerID {
			continue
		}
		// 2. Category
		if f.CategoryID != 0 && car.Category.ID != f.CategoryID {
			continue
		}
		// 3. Year
		if f.MinYear > 0 && car.Year < f.MinYear {
			continue
		}
		// 4. HP
		if f.MinHP > 0 && car.Specs.HP < f.MinHP {
			continue
		}
		// 5. Transmission
		if f.Transmission != "" && car.Specs.Transmission != f.Transmission {
			continue
		}
		// 6. Drivetrain
		if f.Drivetrain != "" && car.Specs.Drivetrain != f.Drivetrain {
			continue
		}

		filtered = append(filtered, car)
	}
	return filtered
}
