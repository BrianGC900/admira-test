package utils

import (
	"context"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type RetryableHTTPClient struct {
	client         HTTPClient
	maxRetries     int
	retryBackoffMs int
	logger         *logrus.Logger
}

func NewRetryableHTTPClient(logger *logrus.Logger, maxRetries, retryBackoffMs int) *RetryableHTTPClient {
	return &RetryableHTTPClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		maxRetries:     maxRetries,
		retryBackoffMs: retryBackoffMs,
		logger:         logger,
	}
}

func (c *RetryableHTTPClient) Do(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for i := 0; i < c.maxRetries; i++ {
		resp, err = c.client.Do(req)
		if err == nil && resp.StatusCode < 500 {
			return resp, nil
		}

		if err != nil {
			c.logger.Warnf("Attempt %d failed: %v", i+1, err)
		} else {
			c.logger.Warnf("Attempt %d failed with status: %d", i+1, resp.StatusCode)
			resp.Body.Close()
		}

		if i < c.maxRetries-1 {
			backoff := time.Duration(math.Pow(2, float64(i))) * time.Duration(c.retryBackoffMs) * time.Millisecond
			c.logger.Debugf("Waiting %v before retry", backoff)
			time.Sleep(backoff)
		}
	}

	return nil, fmt.Errorf("request failed after %d attempts: %v", c.maxRetries, err)
}

func FetchData(ctx context.Context, client HTTPClient, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
