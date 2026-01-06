package catalog

import (
	"context"
	"viewer/internal/domain"
)

type MetadataProvider interface {
	GetManufacturers(ctx context.Context) ([]domain.Manufacturer, error)
	GetCategories(ctx context.Context) ([]domain.Category, error)
}

type Usecase struct {
	carsAPI CarFetcher
	metaAPI MetadataProvider // The interface we defined above
}

type CatalogData struct {
	Cars          []domain.CarSummary
	Manufacturers []domain.Manufacturer
	Categories    []domain.Category
}

func (u *Usecase) LoadCatalog(ctx context.Context) (CatalogData, error) {
	// 1. Start fetching everything (could be parallelized!)
	cars, err := u.carsAPI.FetchCarSummaries(ctx, 100)
	if err != nil {
		return CatalogData{}, err
	}

	// 2. These calls will be INSTANT if cached
	mans, err := u.metaAPI.GetManufacturers(ctx)
	if err != nil {
		return CatalogData{}, err
	}

	cats, err := u.metaAPI.GetCategories(ctx)
	if err != nil {
		return CatalogData{}, err
	}

	// 3. (Optional) Enrich the Car Summaries
	// If CarSummary only has ID, we can look up the Name here
	// using the 'mans' slice we just fetched.

	return CatalogData{
		Cars:          cars,
		Manufacturers: mans,
		Categories:    cats,
	}, nil
}
