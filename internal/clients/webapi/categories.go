package webapi

import (
	"context"
	"encoding/json"
	"viewer/internal/domain"
	"viewer/internal/lib/e"
)

func (c *Client) FetchCategories(ctx context.Context) ([]domain.Category, error) {

	data, err := c.doRequest(ctx, endpointCategories)
	if err != nil {
		return []domain.Category{}, e.Wrap("failed to fetch categories", err)
	}

	var dtos []categoryDTO

	if err := json.Unmarshal(data, &dtos); err != nil {
		return []domain.Category{}, e.Wrap("failed to decode API response for categories", err)
	}

	categories := make([]domain.Category, 0, len(dtos))

	for _, d := range dtos {
		category := domain.Category{
			ID:   d.ID,
			Name: d.Name,
		}
		categories = append(categories, category)
	}

	return categories, nil
}
