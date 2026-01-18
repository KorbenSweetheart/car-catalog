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

func (w *WebRepository) Category(ctx context.Context, ID int) (domain.Category, error) {
	const op = "repository.webapi.Category"

	log := w.log.With(
		slog.String("op", op),
	)

	// Fetch category to fill in the missing car info
	URL, err := url.JoinPath(endpointCategories, strconv.Itoa(ID))
	if err != nil {
		log.Error("invalid url",
			slog.String("url:", URL),
			slog.Any("error", err),
		)
		return domain.Category{}, e.Wrap("invalid url", fmt.Errorf("url: %s, error: %w", URL, err))
	}

	categoryData, err := w.client.DoRequest(ctx, URL)
	if err != nil {
		log.Error("failed to fetch category", slog.Any("error", err))
		return domain.Category{}, e.Wrap("failed to fetch category", err)
	}

	var categoryDTO categoryDTO

	if err := json.Unmarshal(categoryData, &categoryDTO); err != nil {
		log.Error("failed to decode API response for category", slog.Any("error", err))
		return domain.Category{}, e.Wrap("failed to decode API response for category", err)
	}

	// Map DTO to Domain object
	category := domain.Category{
		ID:   ID,
		Name: categoryDTO.Name,
	}

	return category, nil
}
