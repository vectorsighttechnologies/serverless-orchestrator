// Package newrelic — functions.go
//
// Listing Lambda functions and detecting their instrumentation status.
package newrelic

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
	lambdatypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"

	"github.com/vectorsighttechnologies/serverless-orchestrator/lambda/internal/awsclient"
	"github.com/vectorsighttechnologies/serverless-orchestrator/lambda/internal/types"
)

// ListFunctions returns all Lambda functions in the current region
// with their New Relic instrumentation status.
//
// Improvement: Returns structured FunctionInfo instead of raw dicts.
func ListFunctions(ctx context.Context, clients *awsclient.Factory) ([]types.FunctionInfo, error) {
	lambdaClient, err := clients.Lambda(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create Lambda client: %w", err)
	}

	region, err := clients.Region(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get region: %w", err)
	}

	var functions []types.FunctionInfo
	paginator := lambda.NewListFunctionsPaginator(lambdaClient, &lambda.ListFunctionsInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list Lambda functions: %w", err)
		}

		for _, fn := range page.Functions {
			info := types.FunctionInfo{
				Name:         deref(fn.FunctionName),
				Arn:          deref(fn.FunctionArn),
				Runtime:      string(fn.Runtime),
				Handler:      deref(fn.Handler),
				LastModified: deref(fn.LastModified),
				MemorySize:   derefInt32(fn.MemorySize),
				Timeout:      derefInt32(fn.Timeout),
				CodeSize:     fn.CodeSize,
				Status:       "not_instrumented",
				Mode:         "none",
				LayerVersion: "",
			}

			// Determine architecture
			if len(fn.Architectures) > 0 {
				info.Architecture = string(fn.Architectures[0])
			} else {
				info.Architecture = "x86_64"
			}

			// Check for NR instrumentation
			for _, layer := range fn.Layers {
				layerArn := deref(layer.Arn)
				if IsNRLayer(layerArn, region) {
					info.Status = "instrumented"
					// Extract layer version from ARN (last segment after ":")
					parts := strings.Split(layerArn, ":")
					if len(parts) > 0 {
						info.LayerVersion = parts[len(parts)-1]
					}

					// Detect mode from env vars
					info.Mode = detectMode(fn.Environment)
					break
				}
			}

			// Skip the internal helper / orchestrator functions
			if isInternalFunction(info.Name) {
				continue
			}

			functions = append(functions, info)
		}
	}

	return functions, nil
}

// detectMode determines the instrumentation mode from the function's env vars.
func detectMode(env *lambdatypes.EnvironmentResponse) string {
	if env == nil {
		return "none"
	}
	vars := env.Variables
	if vars == nil {
		return "none"
	}

	// Check for APM mode
	if apm, ok := vars["NEW_RELIC_APM_LAMBDA_MODE"]; ok {
		if strings.EqualFold(apm, "true") {
			return "apm"
		}
	}

	// Check for extension (= serverless mode)
	if ext, ok := vars["NEW_RELIC_LAMBDA_EXTENSION_ENABLED"]; ok {
		if strings.EqualFold(ext, "true") {
			return "serverless"
		}
	}

	return "log_ingestion"
}

// deref safely dereferences a *string, returning "" for nil.
func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// derefInt32 safely dereferences a *int32, returning 0 for nil.
func derefInt32(i *int32) int32 {
	if i == nil {
		return 0
	}
	return *i
}

// isInternalFunction checks if a Lambda function name belongs to the orchestrator
// or is an internal helper function deployed by New Relic templates.
func isInternalFunction(name string) bool {
	nameLower := strings.ToLower(name)
	selfName := strings.ToLower(os.Getenv("AWS_LAMBDA_FUNCTION_NAME"))

	// Skip orchestrator itself
	if selfName != "" && nameLower == selfName {
		return true
	}
	if strings.Contains(nameLower, "serverless-orchestrator") {
		return true
	}

	// Skip New Relic integration helper functions
	if strings.Contains(nameLower, "newrelic-log-ingestion") {
		return true
	}
	if strings.Contains(nameLower, "graphqlapicall") {
		return true
	}
	if strings.Contains(nameLower, "graphqlconfigureapicall") {
		return true
	}
	if strings.Contains(nameLower, "loggroupmanagerfunction") {
		return true
	}

	return false
}
