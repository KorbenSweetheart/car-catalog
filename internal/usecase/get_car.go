package usecase

import (
	"context"
	"fmt"
	"viewer/internal/domain"
)

func (c *CarModel) GetCarModel(ctx context.Context, id int) (domain.Car, error) {

	carModel, err := c.repo.GetCarModel(ctx, id)
	if err != nil {
		return domain.Car{}, fmt.Errorf("get profile from postgres: %w", err)
	}

	return carModel, nil
}
