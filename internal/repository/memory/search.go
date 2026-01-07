package memory

import (
	"context"
	"viewer/internal/domain"
)

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
