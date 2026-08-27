package discord

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

const maxResponseBytes = 1 << 20

type Client struct {
	id          string
	secret      string
	redirectURI string
	baseURL     string
	httpClient  *http.Client
}

func NewClient(id, secret, redirectURI string) *Client {
	return &Client{
		id:          id,
		secret:      secret,
		redirectURI: redirectURI,
		baseURL:     "https://discord.com/api/v10",
		httpClient:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) newRequest(
	ctx context.Context,
	method,
	path string,
	body io.Reader,
) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return nil, err
	}

	// version is a place holder, it should actually later come from
	// somewhere standard.
	req.Header.Set("User-Agent", "arbiter/0.1 (https://github.com/sonolink/arbiter)")

	return req, nil
}

func readBody(resp *http.Response) ([]byte, error) {
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes+1))
	if err != nil {
		return nil, err
	}

	if len(body) > maxResponseBytes {
		return nil, fmt.Errorf("discord: response exceeds %d bytes", maxResponseBytes)
	}

	return body, nil
}
