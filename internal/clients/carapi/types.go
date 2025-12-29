package carapi

import (
	"context"
	"viewer/internal/domain"
)

type Adapter interface {
	GetCarModel(ctx context.Context, id int) (domain.Car, error)
	GetManufacturer(ctx context.Context, id int) (domain.Manufacturer, error)
	// Add other methods like GetAll(), GetCategory() here
}

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	ID      int    `json:"update_id"`
	Message string `json:"message"`
}
