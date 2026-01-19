package webapi

import (
	"context"
	"log/slog"

	"gitea.kood.tech/ivanandreev/viewer/internal/domain"
	"gitea.kood.tech/ivanandreev/viewer/internal/lib/e"
)

func (w *WebRepository) Metadata(ctx context.Context) (domain.Metadata, error) {
	const op = "repository.webapi.Metadata"

	log := w.log.With(
		slog.String("op", op),
	)

	vendors, err := w.Manufacturers(ctx)
	if err != nil {
		log.Error("failed to get manufacturers", slog.Any("error", err))
		return domain.Metadata{}, e.Wrap("failed to get manufacturers", err)
	}

	categories, err := w.Categories(ctx)
	if err != nil {
		log.Error("failed to get categories", slog.Any("error", err))
		return domain.Metadata{}, e.Wrap("failed to get categories", err)
	}

	cars, err := w.Cars(ctx)
	if err != nil {
		log.Error("failed to get cars", slog.Any("error", err))
		return domain.Metadata{}, e.Wrap("failed to get cars", err)
	}

	// Sets to deduplicate strings for metadata
	uniqueDrivetrains := make(map[string]bool)
	uniqueTransmissions := make(map[string]bool)

	// The Enrichment Loop
	for _, car := range cars {
		// Collect Metadata Strings
		uniqueDrivetrains[car.Specs.Drivetrain] = true
		uniqueTransmissions[car.Specs.Transmission] = true
	}

	// Convert Maps to Slices
	drivetrains := make([]string, 0, len(uniqueDrivetrains))
	for d := range uniqueDrivetrains {
		drivetrains = append(drivetrains, d)
	}

	transmissions := make([]string, 0, len(uniqueTransmissions))
	for t := range uniqueTransmissions {
		transmissions = append(transmissions, t)
	}

	log.Info("metadata loaded",
		slog.Int("manufacturers_count", len(vendors)),
		slog.Int("categories_count", len(categories)),
		slog.Int("drivetrains_count", len(drivetrains)),
		slog.Int("transmissions_count", len(transmissions)),
	)

	return domain.Metadata{
		Manufacturers: vendors,
		Categories:    categories,
		Drivetrains:   drivetrains,
		Transmissions: transmissions,
	}, nil
}
