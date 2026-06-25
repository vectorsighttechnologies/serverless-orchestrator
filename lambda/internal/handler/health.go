// Package handler — health.go
//
// GET /health — returns which NR credentials are configured.
// The UI uses this to decide whether to show credential input form.
package handler

import (
	"github.com/aws/aws-lambda-go/events"

	"github.com/vectorsighttechnologies/serverless-orchestrator/lambda/internal/config"
	"github.com/vectorsighttechnologies/serverless-orchestrator/lambda/internal/types"
)

// HandleHealth returns the configuration status of the orchestrator.
// This endpoint is called on first connect from the UI.
func HandleHealth(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	cfg, _ := loadConfig(request.Headers)

	source := "none"
	if cfg != nil {
		source = cfg.Source
	}

	resp := types.HealthResponse{
		Status: "ok",
		Config: types.ConfigStatus{
			LicenseKeyConfigured: cfg != nil && cfg.HasLicenseKey(),
			AccountIDConfigured:  cfg != nil && cfg.HasAccountID(),
			APIKeyConfigured:     cfg != nil && cfg.HasAPIKey(),
			Region:               defaultString(cfg, "us"),
			Source:               source,
		},
	}

	return jsonResponse(200, resp)
}

// defaultString returns the config region or a fallback.
func defaultString(cfg *config.Config, fallback string) string {
	if cfg != nil && cfg.Region != "" {
		return cfg.Region
	}
	return fallback
}
