package httpclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
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

func (c *Client) DoRequest(ctx context.Context, path string) (data []byte, err error) {
	const op = "httpclient.DoRequest"

	URL, err := url.JoinPath(c.host, path)
	if err != nil {
		return nil, fmt.Errorf("invalid url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	if err != nil {
		return nil, fmt.Errorf("can't create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't do request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("can't read body: %w", err)
	}

	return body, nil
}
