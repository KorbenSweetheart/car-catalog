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

func (s *CarStore) RecommendedCars(ctx context.Context, viewedIDs []int, excludeID int) ([]domain.Car, error) {
	const op = "usecase.carstore.RecommendedCars"

	log := s.log.With(
		slog.String("op", op),
	)

	const limit = 4

	// 1. If no history return 4 random cars
	if len(viewedIDs) == 0 {
		log.Debug("empty viewedids history")

		randomCars, err := s.repo.RandomCars(ctx, limit)
		if err != nil {
			log.Error("failed to get random cars", slog.Any("error", err))
			return nil, e.Wrap("failed to get random cars: %w", err)
		}

		return randomCars, nil
	}

	// --- RECOMENDATIONS LOGIC ---
	// We count occurrences to find the "Top" 2 cars, 2 top manufacturers and most viewed category

	// Top Viewed Car: "Resume where you left off."
	// 2nd viewed car || top car Manufacturer + top Category && !TopCarID
	// top car Manufacturer + top Category && !TopCarID: "Brand loyalty."
	// 2nd Manufacturer + top Category: "Comparison shopping."
	// Wildcard: same category but none of the prev manufacturers. -> add 2 on top car page.

	// With short personalised list, need to add few random cars.

	// cleanup the viewedIDs, remove duplicates and fetch only unique cars
	// during cleaning, count top 2 most viewed cars
	uniqueIDsMap := make(map[int]int)

	for _, id := range viewedIDs {
		uniqueIDsMap[id]++
	}

	var top1CarID, top2CarID int // WARN: cound be the same brand
	top1CarCount := 0
	top2CarCount := 0

	for id, count := range uniqueIDsMap {
		if count > top1CarCount {
			top2CarID = top1CarID
			top2CarCount = top1CarCount

			top1CarID = id
			top1CarCount = count
		} else if count > top2CarCount {
			top2CarID = id
			top2CarCount = count
		}
	}

	// 2. Fetch cars data to analyze preferences
	// TODO: maybe add cache here to check there first
	foundCars, err := s.repo.CarsByIDs(ctx, uniqueIDsMap)
	if err != nil {
		log.Warn("failed to fetch cars by ids", slog.Any("error", err))

		randomCars, err := s.repo.RandomCars(ctx, limit)
		if err != nil {
			log.Error("failed to get random cars", slog.Any("error", err))
			return nil, e.Wrap("failed to get random cars: %w", err)
		}

		return randomCars, nil
	}

	// create final slice to return
	recommendedCars := make([]domain.Car, 0, 9)

	// 3. Analyze Preferences
	// we need to know:
	// top car manID and catID, topManID, topCatID,
	vendorCounts := make(map[int]int)
	categoryCounts := make(map[int]int)

	var top1ManID, top2ManID, topCategoryID int
	top1ManCount := 0
	top2ManCount := 0
	topCategoryCount := 0

	var top1CarObj, top2CarObj *domain.Car

	// Iterate over the foundCars list to calculate top vendor and top category
	for i := range foundCars {
		if foundCars[i].ID == top1CarID {
			top1CarObj = &foundCars[i]
		}

		if foundCars[i].ID == top2CarID { // && foundCars[i].ID != excludeID
			top2CarObj = &foundCars[i]
		}

		vendorCounts[foundCars[i].Manufacturer.ID]++
		categoryCounts[foundCars[i].Category.ID]++
	}

	for id, count := range vendorCounts {
		if count > top1ManCount {
			top2ManID = top1ManID
			top2ManCount = top1ManCount

			top1ManID = id
			top1ManCount = count
		} else if count > top2ManCount {
			top2ManID = id
			top2ManCount = count
		}
	}

	for id, count := range categoryCounts {
		if count > topCategoryCount {
			topCategoryID = id
			topCategoryCount = count
		}
	}

	if top1CarObj != nil {
		recommendedCars = append(recommendedCars, *top1CarObj)
	}
	if top2CarObj != nil {
		recommendedCars = append(recommendedCars, *top2CarObj)
	}

	additionalVariants, err := s.repo.RecommendedCars(ctx, top1ManID, top2ManID, topCategoryID, top1CarID, top2CarID)
	if err != nil {
		log.Warn("failed to fetch additional variants of cars", slog.Any("error", err))

		randomCars, err := s.repo.RandomCars(ctx, limit)
		if err != nil {
			log.Error("failed to get random cars", slog.Any("error", err))
			return nil, e.Wrap("failed to get random cars: %w", err)
		}

		recommendedCars = append(recommendedCars, randomCars...)
		return recommendedCars, nil
	}

	// add additional variants to top 2 cars
	recommendedCars = append(recommendedCars, additionalVariants...)

	// always 4 random cars
	randomCars, err := s.repo.RandomCars(ctx, limit)
	if err != nil {
		log.Error("failed to get random cars", slog.Any("error", err))
		return recommendedCars, e.Wrap("failed to get random cars: %w", err)
	}

	recommendedCars = append(recommendedCars, randomCars...)

	// need to exclude current page car ID and limit amount of displayed cars to 4
	filteredCars := make([]domain.Car, 0, 4)

	for i := range recommendedCars {
		if recommendedCars[i].ID == excludeID {
			continue
		}

		// Add to list
		filteredCars = append(filteredCars, recommendedCars[i])

		// Stop once we have 4 cars
		if len(filteredCars) == 4 {
			break
		}
	}

	// return top car, second car, car of top man + top cat, car of second man + top cat, car of another man + top cat, 4 random cars

	log.Info("recommended cars loaded",
		slog.Int("cars_count", len(filteredCars)),
	)

	return filteredCars, nil
}
