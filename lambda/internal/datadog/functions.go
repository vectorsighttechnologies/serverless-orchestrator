package datadog

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
	lambdatypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"

	"github.com/vectorsighttechnologies/serverless-orchestrator/lambda/internal/awsclient"
	"github.com/vectorsighttechnologies/serverless-orchestrator/lambda/internal/types"
)

// ListFunctions returns all Lambda functions in the current region
// with their Datadog instrumentation status.
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

			// Check for Datadog and New Relic instrumentation by scanning layers
			isDatadog := false
			isNewRelic := false
			var ddLayerArn string

			for _, layer := range fn.Layers {
				layerArn := deref(layer.Arn)
				if IsDatadogLayer(layerArn, region) {
					isDatadog = true
					ddLayerArn = layerArn
				} else if strings.Contains(strings.ToLower(layerArn), "newrelic") {
					isNewRelic = true
				}
			}

			if isDatadog {
				info.Status = "instrumented"
				parts := strings.Split(ddLayerArn, ":")
				if len(parts) > 0 {
					info.LayerVersion = parts[len(parts)-1]
				}
				info.Mode = detectMode(fn.Environment)
			} else if isNewRelic {
				info.Status = "instrumented_newrelic"
				info.Mode = "serverless"
			}

			// Skip internal orchestrator helper functions
			if isInternalFunction(info.Name) {
				continue
			}

			functions = append(functions, info)
		}
	}

	return functions, nil
}

// detectMode determines the instrumentation mode from the environment variables.
func detectMode(env *lambdatypes.EnvironmentResponse) string {
	if env == nil || env.Variables == nil {
		return "none"
	}
	vars := env.Variables
	if _, ok := vars["DD_API_KEY"]; ok {
		return "serverless"
	}
	return "none"
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
