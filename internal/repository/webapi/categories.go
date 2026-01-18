package webapi

import (
	"context"
	"encoding/json"
	"log/slog"

	"gitea.kood.tech/ivanandreev/viewer/internal/domain"
	"gitea.kood.tech/ivanandreev/viewer/internal/lib/e"
)

func (w *WebRepository) Categories(ctx context.Context) ([]domain.Category, error) {
	const op = "repository.webapi.Categories"

	log := w.log.With(
		slog.String("op", op),
	)

	data, err := w.client.DoRequest(ctx, endpointCategories)
	if err != nil {
		log.Error("failed to fetch categories", slog.Any("error", err))
		return []domain.Category{}, e.Wrap("failed to fetch categories", err)
	}

	var dtos []categoryDTO

	if err := json.Unmarshal(data, &dtos); err != nil {
		log.Error("failed to decode API response for categories", slog.Any("error", err))
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
