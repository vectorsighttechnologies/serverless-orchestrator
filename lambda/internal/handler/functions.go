// Package handler — functions.go
//
// Handlers for Lambda function listing, instrumentation, and uninstrumentation.
//
// Each handler receives a context.Context (from the Lambda runtime) along
// with the API Gateway request and shared AWS clients.
//
// SOLID: Handlers are thin — they parse, delegate to services, and format.
package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/aws/aws-lambda-go/events"

	"github.com/vectorsighttechnologies/serverless-orchestrator/lambda/internal/awsclient"
	"github.com/vectorsighttechnologies/serverless-orchestrator/lambda/internal/datadog"
	"github.com/vectorsighttechnologies/serverless-orchestrator/lambda/internal/newrelic"
	"github.com/vectorsighttechnologies/serverless-orchestrator/lambda/internal/types"
)

// HandleListFunctions returns all Lambda functions with their NR/DD instrumentation status.
//
// GET /functions
func HandleListFunctions(
	ctx context.Context,
	request events.APIGatewayProxyRequest,
	clients *awsclient.Factory,
) events.APIGatewayProxyResponse {
	cfg, _ := loadConfig(request.Headers)
	provider := "newrelic"
	if cfg != nil && cfg.SelectedProvider != "" {
		provider = cfg.SelectedProvider
	}

	var functions []types.FunctionInfo
	var err error

	if provider == "datadog" {
		functions, err = datadog.ListFunctions(ctx, clients)
	} else {
		functions, err = newrelic.ListFunctions(ctx, clients)
	}

	if err != nil {
		return errorResponse(500, "Failed to list functions: "+err.Error())
	}

	return jsonResponse(200, types.FunctionsResponse{Functions: functions})
}

// HandleInstallFunctions instruments the specified Lambda functions.
//
// POST /functions/install
// Body: { "functionArns": [...], "method": "layer"|"log_ingestion", "mode": "serverless"|"apm" }
//
// Supports bulk operations — processes all functions concurrently for speed.
func HandleInstallFunctions(
	ctx context.Context,
	request events.APIGatewayProxyRequest,
	clients *awsclient.Factory,
) events.APIGatewayProxyResponse {
	// Parse request body
	var req types.InstallRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return errorResponse(400, "Invalid request body: "+err.Error())
	}

	if len(req.FunctionArns) == 0 {
		return errorResponse(400, "At least one function ARN is required")
	}

	// Defaults
	if req.Method == "" {
		req.Method = "layer"
	}
	if req.Mode == "" {
		req.Mode = "serverless"
	}

	// Validate method and mode for New Relic (only if not Datadog)
	cfg, err := loadConfig(request.Headers)
	if err != nil {
		return errorResponse(500, "Failed to load config: "+err.Error())
	}
	if err := cfg.Validate(); err != nil {
		return errorResponse(400, err.Error())
	}

	if cfg.SelectedProvider != "datadog" {
		if req.Method != "layer" && req.Method != "log_ingestion" {
			return errorResponse(400, fmt.Sprintf("Invalid method %q. Must be 'layer' or 'log_ingestion'", req.Method))
		}
		if req.Method == "layer" && req.Mode != "serverless" && req.Mode != "apm" {
			return errorResponse(400, fmt.Sprintf("Invalid mode %q. Must be 'serverless' or 'apm'", req.Mode))
		}
	}

	// Process all functions concurrently
	results := processBulk(ctx, req.FunctionArns, func(functionARN string) error {
		if cfg.SelectedProvider == "datadog" {
			return datadog.InstallLayer(ctx, clients, cfg, functionARN)
		}

		switch req.Method {
		case "layer":
			return newrelic.InstallLayer(ctx, clients, cfg, functionARN, req.Mode)
		case "log_ingestion":
			return fmt.Errorf("log_ingestion requires AWS integration setup first (POST /integration/setup)")
		default:
			return fmt.Errorf("unsupported method: %s", req.Method)
		}
	})

	return jsonResponse(200, types.BatchOperationResponse{Results: results})
}

// HandleUninstallFunctions removes NR/DD instrumentation from the specified functions.
//
// POST /functions/uninstall
// Body: { "functionArns": [...] }
func HandleUninstallFunctions(
	ctx context.Context,
	request events.APIGatewayProxyRequest,
	clients *awsclient.Factory,
) events.APIGatewayProxyResponse {
	var req types.UninstallRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return errorResponse(400, "Invalid request body: "+err.Error())
	}

	if len(req.FunctionArns) == 0 {
		return errorResponse(400, "At least one function ARN is required")
	}

	cfg, _ := loadConfig(request.Headers)
	provider := "newrelic"
	if cfg != nil && cfg.SelectedProvider != "" {
		provider = cfg.SelectedProvider
	}

	results := processBulk(ctx, req.FunctionArns, func(functionARN string) error {
		if provider == "datadog" {
			return datadog.UninstallLayer(ctx, clients, functionARN)
		}

		// Remove layer instrumentation
		if err := newrelic.UninstallLayer(ctx, clients, functionARN); err != nil {
			return err
		}
		// Also remove any CW subscription filter (best-effort)
		_ = newrelic.RemoveLogSubscription(ctx, clients, functionARN)
		return nil
	})

	return jsonResponse(200, types.BatchOperationResponse{Results: results})
}

// ─────────────────────────────────────────────────────────────
// Bulk Processing (DRY — shared by install and uninstall)
// ─────────────────────────────────────────────────────────────

// processBulk executes a function on each ARN concurrently and collects results.
// This is the key improvement over the sequential Python CLI.
func processBulk(
	ctx context.Context,
	arns []string,
	operation func(string) error,
) []types.OperationResult {
	var (
		mu      sync.Mutex
		results = make([]types.OperationResult, 0, len(arns))
		wg      sync.WaitGroup
	)

	for _, arn := range arns {
		wg.Add(1)
		go func(functionARN string) {
			defer wg.Done()

			result := types.OperationResult{
				Arn:          functionARN,
				FunctionName: extractFunctionName(functionARN),
				Success:      true,
			}

			if err := operation(functionARN); err != nil {
				result.Success = false
				result.Error = err.Error()
				fmt.Printf("[ORCHESTRATOR] Error processing function %s: %v\n", functionARN, err)
			} else {
				fmt.Printf("[ORCHESTRATOR] Successfully processed function %s\n", functionARN)
			}

			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}(arn)
	}

	wg.Wait()
	return results
}

// extractFunctionName extracts the function name from a Lambda ARN.
// e.g., "arn:aws:lambda:us-east-1:123:function:my-func" → "my-func"
func extractFunctionName(arn string) string {
	parts := strings.Split(arn, ":")
	if len(parts) >= 7 {
		return parts[6]
	}
	return arn
}
