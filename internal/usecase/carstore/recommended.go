package carstore

import (
	"context"
	"log/slog"

	"gitea.kood.tech/ivanandreev/viewer/internal/domain"
	"gitea.kood.tech/ivanandreev/viewer/internal/lib/e"
)

// car handler reads and sets cookies
// handlers that need personalisation only read cookies
// those handlers trigger recommendation usecase and pass cookies data to it, to get slice of cars (recommended)
// use case must parse the cookie, get needed info, after that fetch needed cars from repo, and return them to needed handler

func (s *CarStore) RecommendedCars(ctx context.Context, viewedIDs []int) ([]domain.Car, error) {
	const op = "usecase.carstore.RecommendedCars"

	log := s.log.With(
		slog.String("op", op),
	)

	const limit = 4

	// 1. Cold Start: No history? Return random cars
	if len(viewedIDs) == 0 {
		log.Debug("empty viewedids history")

		randomCars, err := s.repo.RandomCars(ctx, limit)
		if err != nil {
			log.Error("failed to get random cars", slog.Any("error", err))
			return nil, e.Wrap("failed to get random cars: %w", err)
		}

		return randomCars, nil
	}

	// TODO: clean the viewedIDs, remove duplicates and fetch only unique cars
	// during cleaning, count top 2-3 most visited cars
	uniqueIDs := make([]int, 0, len(viewedIDs)/2)

	// 2. Fetch User History (To analyze preferences)
	// TODO: maybe add cache here to check there first
	historyCars, err := s.repo.CarsByIDs(ctx, uniqueIDs)
	if err != nil {
		log.Warn("failed to fetch cars by ids", slog.Any("error", err))

		randomCars, err := s.repo.RandomCars(ctx, limit)
		if err != nil {
			log.Error("failed to get random cars", slog.Any("error", err))
			return nil, e.Wrap("failed to get random cars: %w", err)
		}

		return randomCars, nil
	}

	// 3. Analyze Preferences (Business Logic)
	// We count occurrences to find the "Top" 2 cars, 2 manufacturers and most viewed category
	carCounts := make(map[int]int)
	vendorCounts := make(map[int]int)
	categoryCounts := make(map[int]int)

	// Iterate over the raw ID list (not the unique car objects) to respect frequency
	for _, id := range uniqueIDs {
		// Find the matching car object
		var car domain.Car
		found := false
		for i := range historyCars {
			if historyCars[i].ID == id {
				car = historyCars[i]
				found = true
				break
			}
		}
		if !found {
			continue
		}

		vendorCounts[car.Manufacturer.ID]++
		categoryCounts[car.Category.ID]++
	}

	log.Info("recommended cars loaded",
		slog.Int("cars_count", len(cars)),
	)

	return cars, nil
}
