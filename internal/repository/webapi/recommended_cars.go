package webapi

import (
	"context"
	"log/slog"

	"gitea.kood.tech/ivanandreev/viewer/internal/domain"
	"gitea.kood.tech/ivanandreev/viewer/internal/lib/e"
)

func (w *WebRepository) RecommendedCars(ctx context.Context, topManID int, secondManID int, topCatID int, top1CarID int, top2CarID int) ([]domain.Car, error) {
	const op = "repository.webapi.RandomCars"

	log := w.log.With(
		slog.String("op", op),
	)

	cars, err := w.Cars(ctx)
	if err != nil {
		log.Error("failed to get cars", slog.Any("error", err))
		return []domain.Car{}, e.Wrap("failed to get cars", err)
	}

	// return car of top man + top cat && !topCarID, car of second man + top cat, car of another man + top cat
	result := make([]domain.Car, 0, 3)
	var slot1, slot2, slot3 *domain.Car

	for i := range cars {
		c := &cars[i]

		if c.ID == top1CarID || c.ID == top2CarID || c.Category.ID != topCatID {
			continue
		}

		// car of top man + top cat && !topCarID && !secondCarID
		if slot1 == nil && c.Manufacturer.ID == topManID {
			slot1 = c
		}
		// car of second man + top cat && !topCarID && !secondCarID
		if slot2 == nil && c.Manufacturer.ID == secondManID {
			slot2 = c
		}
		// car of another man + top cat && !topCarID && !secondCarID
		// NOTE: it might be better to show random manufacturer + top category
		if slot3 == nil && c.Manufacturer.ID != topManID && c.Manufacturer.ID != secondManID {
			slot3 = c
		}

		if slot1 != nil && slot2 != nil && slot3 != nil {
			break
		}
	}

	if slot1 != nil {
		result = append(result, *slot1)
	}
	if slot2 != nil {
		result = append(result, *slot2)
	}
	if slot3 != nil {
		result = append(result, *slot3)
	}

	return result, nil
}
