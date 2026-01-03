package usecase

import (
	"context"
	"viewer/internal/domain"
)

type Adapter interface {
	FetchCarSummary(ctx context.Context, number int) (domain.CarSummary, error)
	FetchCar(ctx context.Context, id int) (domain.Car, error)
	FetchCars(ctx context.Context, id int) ([]domain.Car, error)
	FetchManufacturer(ctx context.Context, id int) (domain.Manufacturer, error)
	FetchManufacturers(ctx context.Context, id int) ([]domain.Manufacturer, error)
}

type UseCase struct {
	adapter Adapter
}

func NewUseCase(a Adapter) *UseCase {
	return &UseCase{
		adapter: a,
	}
}
