package lambda

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/vectorsighttechnologies/serverless-orchestrator/backend/internal/types"
)

// Client handles outgoing HTTP requests to the AWS Lambda Orchestrator API Gateway.
type Client struct {
	httpClient *http.Client
}

// NewClient creates a new Lambda Orchestrator client.
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 45 * time.Second, // AWS API gateway integrations timeout after 29s anyway
		},
	}
}

// Invoke sends a request to the configured Lambda API Gateway, passing authentication and NR headers.
func (c *Client) Invoke(
	ctx context.Context,
	method string,
	path string,
	body []byte,
	lambdaUrl string,
	lambdaApiKey string,
	prefs *types.UserPreferences,
) ([]byte, int, error) {
	if lambdaUrl == "" {
		return nil, 0, errors.New("lambda API endpoint is not configured")
	}

	url := fmt.Sprintf("%s/%s", strings.TrimSuffix(lambdaUrl, "/"), strings.TrimPrefix(path, "/"))

	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	// 1. Set request headers
	req.Header.Set("Content-Type", "application/json")
	if lambdaApiKey != "" {
		req.Header.Set("x-api-key", lambdaApiKey)
	}


	// 2. Set New Relic credentials headers (Orchestrator config fallbacks)
	if prefs.NRLicenseKey != "" {
		req.Header.Set("x-nr-license-key", prefs.NRLicenseKey)
	}
	if prefs.NRAccountID != "" {
		req.Header.Set("x-nr-account-id", prefs.NRAccountID)
	}
	if prefs.NRApiKey != "" {
		req.Header.Set("x-nr-api-key", prefs.NRApiKey)
	}
	if prefs.NRRegion != "" {
		req.Header.Set("x-nr-region", prefs.NRRegion)
	}

	// 3. Perform request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to dispatch request to Lambda: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read Lambda response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return respBody, resp.StatusCode, fmt.Errorf("lambda API returned failure status code %d", resp.StatusCode)
	}

	return respBody, resp.StatusCode, nil
}

// A simple local strings wrapper to trim prefix/suffix to avoid separate import issues
func stringsTrim(val, cut string) string {
	val = strings.TrimSuffix(val, cut)
	return strings.TrimPrefix(val, cut)
}
