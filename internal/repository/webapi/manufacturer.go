package webapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"viewer/internal/domain"
	"viewer/internal/lib/e"
)

func (w *WebRepository) Manufacturer(ctx context.Context, ID int) (domain.Manufacturer, error) {
	const op = "repository.webapi.Manufacturer"

	log := w.log.With(
		slog.String("op", op),
	)

	// Fetch manufacturer to fill in the missing car info
	URL, err := url.JoinPath(endpointManufacturers, strconv.Itoa(ID))
	if err != nil {
		log.Error("invalid url",
			slog.String("url:", URL),
			slog.Any("error", err),
		)
		return domain.Manufacturer{}, e.Wrap("invalid url", fmt.Errorf("url: %s, error: %w", URL, err))
	}

	vendorData, err := w.client.DoRequest(ctx, URL)
	if err != nil {
		log.Error("failed to fetch manufacturer", slog.Any("error", err))
		return domain.Manufacturer{}, e.Wrap("failed to fetch manufacturer", err)
	}

	var vendorDTO manufacturerDTO

	if err := json.Unmarshal(vendorData, &vendorDTO); err != nil {
		log.Error("failed to decode API response for manufacturer", slog.Any("error", err))
		return domain.Manufacturer{}, e.Wrap("failed to decode API response for manufacturer", err)
	}

	// Map DTO to Domain object
	vendor := domain.Manufacturer{
		ID:           ID,
		Name:         vendorDTO.Name,
		Country:      vendorDTO.Country,
		FoundingYear: vendorDTO.FoundingYear,
	}

	return vendor, nil
}
