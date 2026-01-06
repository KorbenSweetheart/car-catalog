package webapi

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
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

// The `image` property relates to an image for a `carModel`, and can be found in the `/api/images` directory.

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
		return nil, e.Wrap("unexpected error", fmt.Errorf("unexpected status code: %d", resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, e.Wrap("can't read body", err)
	}

	return body, nil
}
