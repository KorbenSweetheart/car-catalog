package carstore

import (
	"context"
	"log/slog"
	"viewer/internal/domain"
	"viewer/internal/lib/e"
)

// The Business Logic -> provide metadata to form catalog filters

type MetadataProvider interface {
	Metadata(ctx context.Context) (domain.Metadata, error)
}

type Metadata struct {
	log      *slog.Logger
	metadata MetadataProvider
}

/* Catalog filters
Manufacturer
Category
Drive Train
Transmission
Year min (user input)
Power (hp) (user input)
*/

func (f *Metadata) Filters(ctx context.Context) (domain.Metadata, error) {
	const op = "usecase.carstore.Filters"

	log := f.log.With("op", op)

	filters, err := f.metadata.Metadata(ctx)
	if err != nil {
		log.Error("failed to get metadata", slog.Any("error", err))
		return domain.Metadata{}, e.Wrap("failed to get metadata: %w", err)
	}

	return filters, nil
}
