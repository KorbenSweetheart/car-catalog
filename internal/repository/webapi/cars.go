package webapi

import (
	"context"
	"encoding/json"
	"log/slog"

	"gitea.kood.tech/ivanandreev/viewer/internal/domain"
	"gitea.kood.tech/ivanandreev/viewer/internal/lib/e"
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

	for i := range dtos {
		car := domain.Car{

			ID:    dtos[i].ID,
			Name:  dtos[i].Name,
			Year:  dtos[i].Year,
			Image: w.imageURL(dtos[i].Image),

			// Map the Nested Specs Struct
			Specs: domain.Specs{
				Engine:       dtos[i].Specs.Engine,
				HP:           dtos[i].Specs.HP,
				Gearbox:      dtos[i].Specs.Gearbox,
				Transmission: NormalizeGearbox(dtos[i].Specs.Gearbox),
				Drivetrain:   dtos[i].Specs.Drivetrain,
			},

			// PARTIAL FILL: We only know and need the ID right now.
			Manufacturer: domain.Manufacturer{
				ID: dtos[i].ManufacturerId,
			},
			Category: domain.Category{
				ID: dtos[i].CategoryId,
			},
		}
		cars = append(cars, car)
	}

	return cars, nil
}
