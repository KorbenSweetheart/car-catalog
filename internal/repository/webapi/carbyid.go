package webapi

import (
	"context"
	"log/slog"

	"gitea.kood.tech/ivanandreev/viewer/internal/domain"
	"gitea.kood.tech/ivanandreev/viewer/internal/lib/e"
)

func (w *WebRepository) CarsByIDs(ctx context.Context, viewedIDs map[int]int) ([]domain.Car, error) {
	const op = "repository.webapi.CarsByIDs"

	log := w.log.With(
		slog.String("op", op),
	)

	// Its better to have a separate request and fetch exact cars with 1 request,
	// but unfortunatelly webapi provided for the study project doesn't support it,
	// For DB, it would also be a single retrieval request,
	// So, featching full list of cars here, and with already written method is a fast and convenient solution,

	cars, err := w.Cars(ctx)
	if err != nil {
		log.Error("failed to get cars", slog.Any("error", err))
		return []domain.Car{}, e.Wrap("failed to get cars", err)
	}

	foundCars := make([]domain.Car, 0, len(viewedIDs))

	// Find the matching car object
	for _, car := range cars {
		if _, ok := viewedIDs[car.ID]; ok {
			foundCars = append(foundCars, car)

			if len(foundCars) == len(viewedIDs) {
				break
			}
		}
	}

	log.Info("cars loaded",
		slog.Int("cars_count", len(foundCars)),
	)

	return foundCars, nil
}
