package usecase

import (
	"context"
	"viewer/internal/clients/carapi"
	"viewer/internal/domain"
)

// 3. The UseCase Interface
// The Handler uses this to talk to the UseCase.
type Repository interface {
	GetCarModel(ctx context.Context, id int) (domain.Car, error)
	GetManufacturer(ctx context.Context, id int) (domain.Manufacturer, error)
}

type CarModel struct {
	repo Repository
}

func NewCarModel(c *carapi.Client) *CarModel {
	return &CarModel{
		repo: c,
	}
}
