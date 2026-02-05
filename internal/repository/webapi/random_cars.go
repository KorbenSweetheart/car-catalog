package webapi

import (
	"context"
	"log/slog"
	"math/rand/v2"

	"gitea.kood.tech/ivanandreev/viewer/internal/domain"
	"gitea.kood.tech/ivanandreev/viewer/internal/lib/e"
)

func (w *WebRepository) RandomCars(ctx context.Context, limit int) ([]domain.Car, error) {
	const op = "repository.webapi.RandomCars"

	log := w.log.With(
		slog.String("op", op),
	)

	cars, err := w.Cars(ctx)
	if err != nil {
		log.Error("failed to get cars", slog.Any("error", err))
		return []domain.Car{}, e.Wrap("failed to get cars", err)
	}

	// 2. Shuffle in Go Memory (not efficient on a big scale, but with current webapi it would be ok)

	// its better to use one more copy for security, but in this case its ok to oporeate directly on a copy of cars slice.
	shuffled := make([]domain.Car, len(cars))
	copy(shuffled, cars)

	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	if limit > len(shuffled) {
		limit = len(shuffled)
	}

	randomCars := make([]domain.Car, limit)

	copy(randomCars, shuffled)

	return randomCars, nil
}
