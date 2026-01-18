package webapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"

	"gitea.kood.tech/ivanandreev/viewer/internal/domain"
	"gitea.kood.tech/ivanandreev/viewer/internal/lib/e"
)

// Car(ctx context.Context, id int) (domain.Car, error)
func (w *WebRepository) Car(ctx context.Context, id int) (domain.Car, error) {
	const op = "repository.webapi.Car"

	log := w.log.With(
		slog.String("op", op),
	)

	URL, err := url.JoinPath(endpointModels, strconv.Itoa(id))
	if err != nil {
		log.Error("invalid url",
			slog.String("url:", URL),
			slog.Any("error", err),
		)
		return domain.Car{}, e.Wrap("invalid url", fmt.Errorf("url: %s, error: %w", URL, err))
	}

	carData, err := w.client.DoRequest(ctx, URL)
	if err != nil {
		log.Error("failed to fetch car", slog.Any("error", err))
		return domain.Car{}, e.Wrap("failed to fetch car", err)
	}

	var carDTO carDTO

	if err := json.Unmarshal(carData, &carDTO); err != nil {
		log.Error("failed to decode API response for cars", slog.Any("error", err))
		return domain.Car{}, e.Wrap("failed to decode API response for cars", err)
	}

	// Fetch manufacturer to fill in the missing car info
	vendor, err := w.Manufacturer(ctx, carDTO.ManufacturerId)
	if err != nil {
		log.Error("failed to get manufacturer", slog.Any("error", err))
		return domain.Car{}, e.Wrap("failed to get manufacturer", err)
	}

	// Fetch category to fill in the missing car info
	category, err := w.Category(ctx, carDTO.CategoryId)
	if err != nil {
		log.Error("failed to get category", slog.Any("error", err))
		return domain.Car{}, e.Wrap("failed to get category", err)
	}

	// Map DTO to Domain object
	car := domain.Car{

		ID:    carDTO.ID,
		Name:  carDTO.Name,
		Year:  carDTO.Year,
		Image: w.imageURL(carDTO.Image),

		// Map the nested Specs struct
		Specs: domain.Specs{
			Engine:     carDTO.Specs.Engine,
			HP:         carDTO.Specs.HP,
			Gearbox:    carDTO.Specs.Gearbox,
			Drivetrain: carDTO.Specs.Drivetrain,
		},

		// Map the nested Vendor struct
		Manufacturer: domain.Manufacturer{
			ID:           carDTO.ManufacturerId,
			Name:         vendor.Name,
			Country:      vendor.Country,
			FoundingYear: vendor.FoundingYear,
		},

		// Map the nested Category struct
		Category: domain.Category{
			ID:   carDTO.CategoryId,
			Name: category.Name,
		},
	}

	return car, nil
}
