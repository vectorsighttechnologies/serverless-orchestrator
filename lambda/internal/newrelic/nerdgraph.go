package newrelic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/vectorsighttechnologies/serverless-orchestrator/lambda/internal/config"
)

// NerdGraphClient handles interactions with New Relic's NerdGraph GraphQL API.
type NerdGraphClient struct {
	httpClient *http.Client
	apiKey     string
	endpoint   string
}

// NewNerdGraphClient creates a new NerdGraph GraphQL client.
func NewNerdGraphClient(cfg *config.Config) *NerdGraphClient {
	endpoint := "https://api.newrelic.com/graphql"
	if cfg.Region == "eu" {
		endpoint = "https://api.eu.newrelic.com/graphql"
	}
	return &NerdGraphClient{
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		apiKey:     cfg.APIKey,
		endpoint:   endpoint,
	}
}

// GraphQLRequest defines a standard GraphQL request payload.
type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// GraphQLResponse represents the standard GraphQL response structure.
type GraphQLResponse struct {
	Data   interface{} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// Query executes a GraphQL query or mutation.
func (c *NerdGraphClient) Query(ctx context.Context, query string, variables map[string]interface{}, responseTarget interface{}) error {
	if c.apiKey == "" {
		return fmt.Errorf("New Relic User API key is not configured")
	}

	payload := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal GraphQL request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create http request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("API-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("GraphQL HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GraphQL request failed with status %d: %s", resp.StatusCode, string(respBytes))
	}

	// Unmarshal standard envelope to extract errors
	var envelope struct {
		Data   json.RawMessage `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := json.Unmarshal(respBytes, &envelope); err != nil {
		return fmt.Errorf("failed to parse GraphQL response envelope: %w. raw: %s", err, string(respBytes))
	}

	if len(envelope.Errors) > 0 {
		var errMsg []string
		for _, e := range envelope.Errors {
			errMsg = append(errMsg, e.Message)
		}
		return fmt.Errorf("GraphQL errors: %s", strings.Join(errMsg, "; "))
	}

	if responseTarget != nil {
		if err := json.Unmarshal(envelope.Data, responseTarget); err != nil {
			return fmt.Errorf("failed to parse GraphQL data payload: %w. raw: %s", err, string(respBytes))
		}
	}

	return nil
}

// LinkedAccount represents the New Relic linked account info.
type LinkedAccount struct {
	ID         interface{} `json:"id"`
	Name       string      `json:"name"`
	ExternalID string      `json:"externalId"`
}

// CloudQueryResponse represents the NerdGraph response to query linked cloud accounts.
type CloudQueryResponse struct {
	Actor struct {
		Account struct {
			Cloud struct {
				LinkedAccounts []LinkedAccount `json:"linkedAccounts"`
			} `json:"cloud"`
		} `json:"account"`
	} `json:"actor"`
}
