// Package handler — integration.go
//
// Handlers for AWS Integration management (Metric Streams + API Polling).
//
// Endpoints:
//   - GET  /integration/status → Check CF stack status
//   - POST /integration/setup  → Deploy integration CF stack
//   - POST /integration/remove → Delete integration CF stack
package handler

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"

	"github.com/vectorsighttechnologies/serverless-orchestrator/lambda/internal/awsclient"
	"github.com/vectorsighttechnologies/serverless-orchestrator/lambda/internal/datadog"
	"github.com/vectorsighttechnologies/serverless-orchestrator/lambda/internal/newrelic"
	"github.com/vectorsighttechnologies/serverless-orchestrator/lambda/internal/types"
)

// HandleIntegrationStatus returns the current NR/DD integration status.
//
// GET /integration/status
func HandleIntegrationStatus(
	ctx context.Context,
	request events.APIGatewayProxyRequest,
	clients *awsclient.Factory,
) events.APIGatewayProxyResponse {
	cfg, err := loadConfig(request.Headers)
	if err != nil {
		return errorResponse(500, "Failed to load config: "+err.Error())
	}
	if err := cfg.Validate(); err != nil {
		return errorResponse(400, err.Error())
	}

	var status *types.IntegrationStatusResponse
	if cfg.SelectedProvider == "datadog" {
		status, err = datadog.GetIntegrationStatus(ctx, clients, cfg)
	} else {
		status, err = newrelic.GetIntegrationStatus(ctx, clients, cfg)
	}

	if err != nil {
		return errorResponse(500, "Failed to check integration status: "+err.Error())
	}

	return jsonResponse(200, status)
}

// HandleIntegrationSetup deploys the NR/DD integration CF stack.
//
// POST /integration/setup
// Body: { "method": "metric_streams"|"api_polling", "includeLogs": true|false }
func HandleIntegrationSetup(
	ctx context.Context,
	request events.APIGatewayProxyRequest,
	clients *awsclient.Factory,
) events.APIGatewayProxyResponse {
	var req types.IntegrationSetupRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return errorResponse(400, "Invalid request body: "+err.Error())
	}

	if req.Method == "" {
		return errorResponse(400, "Integration method is required ('metric_streams' or 'api_polling')")
	}

	cfg, err := loadConfig(request.Headers)
	if err != nil {
		return errorResponse(500, "Failed to load config: "+err.Error())
	}
	if err := cfg.Validate(); err != nil {
		return errorResponse(400, err.Error())
	}

	var resp *types.IntegrationSetupResponse
	if cfg.SelectedProvider == "datadog" {
		resp, err = datadog.SetupIntegration(ctx, clients, cfg, &req)
	} else {
		resp, err = newrelic.SetupIntegration(ctx, clients, cfg, &req)
	}

	if err != nil {
		return errorResponse(500, "Failed to setup integration: "+err.Error())
	}

	return jsonResponse(200, resp)
}

// HandleIntegrationRemove deletes the NR/DD integration CF stack.
//
// POST /integration/remove
func HandleIntegrationRemove(
	ctx context.Context,
	request events.APIGatewayProxyRequest,
	clients *awsclient.Factory,
) events.APIGatewayProxyResponse {
	cfg, err := loadConfig(request.Headers)
	if err != nil {
		return errorResponse(500, "Failed to load config: "+err.Error())
	}
	if err := cfg.Validate(); err != nil {
		return errorResponse(400, err.Error())
	}

	if cfg.SelectedProvider == "datadog" {
		err = datadog.RemoveIntegration(ctx, clients, cfg)
	} else {
		err = newrelic.RemoveIntegration(ctx, clients, cfg)
	}

	if err != nil {
		return errorResponse(500, "Failed to remove integration: "+err.Error())
	}

	return jsonResponse(200, map[string]string{
		"status":  "deleting",
		"message": "Integration stack deletion initiated",
	})
}
