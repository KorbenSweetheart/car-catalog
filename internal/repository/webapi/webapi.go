package webapi

import (
	"context"
	"log/slog"
)

/*
The api exposes the following endpoints:
```
GET /api/models
GET /api/models/{id}
GET /api/manufacturers
GET /api/manufacturers/{id}
GET /api/categories
GET /api/categories/{id}
```
*/

// The `image` property relates to an image for a `carModel`, and can be found in the `/api/images` directory.

const (
	endpointModels        = "models"
	endpointManufacturers = "manufacturers"
	endpointCategories    = "categories"
)

type Client interface {
	DoRequest(ctx context.Context, path string) (data []byte, err error)
}

type WebRepository struct {
	log       *slog.Logger
	client    Client
	mediaHost string
}

func New(log *slog.Logger, client Client, mediaHost string) *WebRepository {
	return &WebRepository{
		log:       log,
		client:    client,
		mediaHost: mediaHost,
	}
}
