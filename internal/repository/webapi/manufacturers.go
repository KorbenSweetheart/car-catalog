package webapi

import (
	"context"
	"encoding/json"
	"log/slog"

	"gitea.kood.tech/ivanandreev/viewer/internal/domain"
	"gitea.kood.tech/ivanandreev/viewer/internal/lib/e"
)

func (w *WebRepository) Manufacturers(ctx context.Context) ([]domain.Manufacturer, error) {
	const op = "repository.webapi.Manufacturers"

	log := w.log.With(
		slog.String("op", op),
	)

	data, err := w.client.DoRequest(ctx, endpointManufacturers)
	if err != nil {
		log.Error("failed to fetch manufacturers", slog.Any("error", err))
		return []domain.Manufacturer{}, e.Wrap("failed to fetch manufacturers", err)
	}

	var dtos []manufacturerDTO

	if err := json.Unmarshal(data, &dtos); err != nil {
		log.Error("failed to decode API response for manufacturers", slog.Any("error", err))
		return []domain.Manufacturer{}, e.Wrap("failed to decode API response for manufacturers", err)
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
