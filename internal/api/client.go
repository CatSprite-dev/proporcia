package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	httpClient http.Client
	baseURL    string
	logger     *slog.Logger
}

func NewClient(baseURL string, logger *slog.Logger) *Client {
	return &Client{
		httpClient: http.Client{
			Timeout: 15 * time.Second,
		},
		baseURL: baseURL,
		logger:  logger,
	}
}

func (c *Client) DoRequest(ctx context.Context, url string, httpMethod string, token string, payload interface{}) ([]byte, error) {
	start := time.Now()
	endpoint := url[strings.LastIndex(url, "/")+1:]

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	c.logger.Debug("sending request", "method", httpMethod, "endpoint", endpoint)

	req, err := http.NewRequestWithContext(ctx, httpMethod, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status %d: %s", res.StatusCode, res.Status)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	c.logger.Debug("request completed", "endpoint", endpoint, "status", res.StatusCode, "duration", time.Since(start))

	return data, nil
}
