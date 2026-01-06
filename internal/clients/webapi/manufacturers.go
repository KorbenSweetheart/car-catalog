package webapi

import (
	"context"
	"encoding/json"
	"viewer/internal/domain"
	"viewer/internal/lib/e"
)

func (c *Client) FetchManufacturers(ctx context.Context) ([]domain.Manufacturer, error) {

	data, err := c.doRequest(ctx, endpointManufacturers)
	if err != nil {
		return []domain.Manufacturer{}, e.Wrap("failed to fetch manufacturers", err)
	}

	var dtos []manufacturerDTO

	if err := json.Unmarshal(data, &dtos); err != nil {
		return []domain.Manufacturer{}, e.Wrap("failed to decode API response for manufacturer", err)
	}

	vendors := make([]domain.Manufacturer, 0, len(dtos))

	for _, d := range dtos {
		vendor := domain.Manufacturer{
			ID:           d.ID,
			Name:         d.Name,
			Country:      d.Country,
			FoundingYear: d.FoundingYear,
		}
		vendors = append(vendors, vendor)
	}

	return vendors, nil
}
