package carapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
	"viewer/internal/domain"
	"viewer/internal/lib/e"
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

const (
	endpointModels        = "models"
	endpointManufacturers = "manufacturers"
	endpointCategories    = "categories"
)

type Client struct {
	host    string
	timeout time.Duration
	client  http.Client
}

func New(host string, timeout time.Duration) *Client {
	return &Client{
		host:   host,
		client: http.Client{Timeout: timeout},
	}
}

func (c *Client) FetchCarSummary(ctx context.Context, number int) (domain.CarSummary, error) {
	path := fmt.Sprintf("%s/%d", endpointModels, number)

	data, err := c.doRequest(ctx, path)
	if err != nil {
		return domain.CarSummary{}, e.Wrap("can't fetch car", err)
	}

	var carSum domain.CarSummary

	if err := json.Unmarshal(data, &carSum); err != nil {
		return domain.CarSummary{}, e.Wrap("can't decode response for car", err)
	}

	return carSum, nil
}

func (c *Client) FetchCar(ctx context.Context, id int) (domain.Car, error) {

	path := fmt.Sprintf("%s/%d", endpointModels, id)

	data, err := c.doRequest(ctx, path)
	if err != nil {
		return domain.Car{}, e.Wrap("can't fetch car", err)
	}

	var car domain.Car

	if err := json.Unmarshal(data, &car); err != nil {
		return domain.Car{}, e.Wrap("can't decode response for car", err)
	}

	return car, nil
}

func (c *Client) FetchCars(ctx context.Context, id int) (domain.Car, error) {

	path := fmt.Sprintf("%s/%d", endpointModels, id)

	data, err := c.doRequest(ctx, path)
	if err != nil {
		return domain.Car{}, e.Wrap("can't fetch car", err)
	}

	var car domain.Car

	if err := json.Unmarshal(data, &car); err != nil {
		return domain.Car{}, e.Wrap("can't decode response for car", err)
	}

	return car, nil
}

func (c *Client) FetchManufacturer(ctx context.Context, id int) (domain.Manufacturer, error) {

	path := fmt.Sprintf("%s/%d", endpointManufacturers, id)

	data, err := c.doRequest(ctx, path)
	if err != nil {
		return domain.Manufacturer{}, e.Wrap("can't fetch manufacturer", err)
	}

	var vendor domain.Manufacturer

	if err := json.Unmarshal(data, &vendor); err != nil {
		return domain.Manufacturer{}, e.Wrap("can't decode response for manufacturer", err)
	}

	return vendor, nil
}

func (c *Client) doRequest(ctx context.Context, path string) (data []byte, err error) {

	URL, err := url.JoinPath(c.host, path)
	if err != nil {
		return nil, e.Wrap("invalid url", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	if err != nil {
		return nil, e.Wrap("can't create request", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, e.Wrap("can't do request", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, e.Wrap("can't read body", err)
	}

	return body, nil
}
