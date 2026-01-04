package home

import (
	"context"
	"fmt"
	"viewer/internal/domain"
)

type CarFetcher interface { // CarSummaryFetcher
	FetchCarSummary(ctx context.Context, number int) ([]domain.CarSummary, error)
}

type Usecase struct {
	fetcher CarFetcher
}

// 3. The Constructor
// When called from app, it looks like: home.New(adapter)
func New(fetcher CarFetcher) *Usecase {
	return &Usecase{
		fetcher: fetcher,
	}
}

// 4. The Business Logic
func (u *Usecase) GetFeaturedCars(ctx context.Context, id int) (domain.Car, error) {

	// Logic: Get 4 cars, maybe filter them, maybe handle errors specifically

	carModel, err := u.fetcher.FetchCarSummary(ctx, id)
	if err != nil {
		return domain.Car{}, fmt.Errorf("get profile from postgres: %w", err)
	}

	return carModel, nil
}
