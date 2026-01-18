package webapi

import (
	"context"
	"encoding/json"
	"log/slog"
	"viewer/internal/domain"
	"viewer/internal/lib/e"
)

func (w *WebRepository) Cars(ctx context.Context) ([]domain.Car, error) {
	const op = "repository.webapi.Cars"

	log := w.log.With(
		slog.String("op", op),
	)

	data, err := w.client.DoRequest(ctx, endpointModels)
	if err != nil {
		log.Error("failed to fetch cars",
			slog.String("endpoint:", endpointModels),
			slog.Any("error", err),
		)
		return []domain.Car{}, e.Wrap("failed to fetch cars", err)
	}

	var dtos []carDTO

	if err := json.Unmarshal(data, &dtos); err != nil {
		log.Error("failed to decode API response for cars", slog.Any("error", err))
		return []domain.Car{}, e.Wrap("failed to decode API response for cars", err)
	}

	// Map DTO to Domain
	cars := make([]domain.Car, 0, len(dtos))

	for _, d := range dtos {
		car := domain.Car{

			ID:    d.ID,
			Name:  d.Name,
			Year:  d.Year,
			Image: w.imageURL(d.Image),

			// Map the Nested Specs Struct
			Specs: domain.Specs{
				Engine:     d.Specs.Engine,
				HP:         d.Specs.HP,
				Gearbox:    d.Specs.Gearbox,
				Drivetrain: d.Specs.Drivetrain,
			},

			// PARTIAL FILL: We only know and need the ID right now.
			Manufacturer: domain.Manufacturer{
				ID: d.ManufacturerId,
			},
			Category: domain.Category{
				ID: d.CategoryId,
			},
		}
		cars = append(cars, car)
	}

	return cars, nil
}
