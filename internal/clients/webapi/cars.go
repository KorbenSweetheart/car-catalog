package webapi

import (
	"context"
	"encoding/json"
	"viewer/internal/domain"
	"viewer/internal/lib/e"
)

func (c *Client) FetchCars(ctx context.Context) ([]domain.Car, error) {

	data, err := c.doRequest(ctx, endpointModels)
	if err != nil {
		return []domain.Car{}, e.Wrap("failed to fetch cars", err)
	}

	var dtos []carDTO

	if err := json.Unmarshal(data, &dtos); err != nil {
		return []domain.Car{}, e.Wrap("failed to decode API response for cars", err)
	}

	// Map DTO to Domain
	cars := make([]domain.Car, 0, len(dtos))

	for _, d := range dtos {
		car := domain.Car{

			ID:    d.ID,
			Name:  d.Name,
			Year:  d.Year,
			Image: d.Image,

			// Map the Nested Specs Struct
			Specs: domain.Specs{
				Engine:     d.Specs.Engine,
				HP:         d.Specs.HP,
				Gearbox:    d.Specs.Gearbox,
				Drivetrain: d.Specs.Drivetrain,
			},

			// PARTIAL FILL: We only know the ID right now.
			// The Repository will use this ID to look up the Name/Country later.
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
