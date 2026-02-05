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
	var slot1, slot2, slot3 domain.Car
	var found1, found2, found3 bool

	for i := range cars {
		c := &cars[i]

		if c.ID == top1CarID || c.ID == top2CarID || c.Category.ID != topCatID {
			continue
		}

		if !found1 && c.Manufacturer.ID == topManID { // car of top vendor + top category && !topCarID && !secondCarID
			slot1 = *c
			found1 = true
		} else if !found2 && c.Manufacturer.ID == secondManID { // car of second vendor + top category && !topCarID && !secondCarID
			slot2 = *c
			found2 = true
		} else if !found3 && c.Manufacturer.ID != topManID && c.Manufacturer.ID != secondManID { // car of another vendor + top cat && !topCarID && !secondCarID // NOTE: it might be better to show random manufacturer + top category
			slot3 = *c
			found3 = true
		}

		if found1 && found2 && found3 {
			break
		}
	}

	if slot1.ID != 0 {
		result = append(result, slot1)
	}
	if slot2.ID != 0 {
		result = append(result, slot2)
	}
	if slot3.ID != 0 {
		result = append(result, slot3)
	}

	return result, nil
}
