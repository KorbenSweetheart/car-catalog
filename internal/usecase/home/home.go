package home

import (
	"context"
	"fmt"
	"viewer/internal/domain"
)

type CarFetcher interface { // CarSummaryFetcher
	FetchCarSummaries(ctx context.Context, limit int) ([]domain.CarSummary, error)
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
func (u *Usecase) GetFeaturedCars(ctx context.Context, limit int) ([]domain.CarSummary, error) {

	// Logic: Get 4 cars, maybe filter them, maybe handle errors specifically

	CarSummaries, err := u.fetcher.FetchCarSummaries(ctx, limit)
	if err != nil {
		return []domain.CarSummary{}, fmt.Errorf("can't fetch car summaries: %w", err)
	}

	return CarSummaries, nil
}
