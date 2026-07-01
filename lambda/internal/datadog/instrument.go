package datadog

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	lambdasvc "github.com/aws/aws-sdk-go-v2/service/lambda"
	lambdatypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"

	"github.com/vectorsight/serverless-tool/lambda/internal/awsclient"
	"github.com/vectorsight/serverless-tool/lambda/internal/config"
)

// InstallLayer instruments a single Lambda function with Datadog.
func InstallLayer(
	ctx context.Context,
	clients *awsclient.Factory,
	cfg *config.Config,
	functionARN string,
) error {
	targetName := functionARN
	if parts := strings.Split(functionARN, ":"); len(parts) >= 7 {
		targetName = parts[6]
	}
	if isInternalFunction(targetName) {
		return fmt.Errorf("cannot instrument the orchestrator or internal helper function (%s)", targetName)
	}

	lambdaClient, err := clients.Lambda(ctx)
	if err != nil {
		return fmt.Errorf("failed to create Lambda client: %w", err)
	}

	region, err := clients.Region(ctx)
	if err != nil {
		return fmt.Errorf("failed to get region: %w", err)
	}

	// Step 1: Get current function config
	funcConfig, err := lambdaClient.GetFunction(ctx, &lambdasvc.GetFunctionInput{
		FunctionName: aws.String(functionARN),
	})
	if err != nil {
		return fmt.Errorf("failed to get function %s: %w", functionARN, err)
	}

	runtime := string(funcConfig.Configuration.Runtime)
	if !IsSupportedRuntime(runtime) {
		return fmt.Errorf("unsupported runtime %q for function %s", runtime, functionARN)
	}

	architectures := funcConfig.Configuration.Architectures
	architecture := "x86_64"
	if len(architectures) > 0 {
		architecture = string(architectures[0])
	}

	// Step 2: Build target layers list
	// Start with the standard extension layer
	extLayerARN := BuildExtensionLayerARN(region, architecture, DefaultExtensionVersion)
	targetLayers := []string{extLayerARN}

	// If runtime needs a library layer, append it
	libName := GetLibraryNameForRuntime(runtime)
	if libName != "" {
		libLayerARN := BuildLibraryLayerARN(region, libName, DefaultLibraryVersion)
		targetLayers = append(targetLayers, libLayerARN)
	}

	// Preserve existing non-Datadog layers
	for _, layer := range funcConfig.Configuration.Layers {
		arn := aws.ToString(layer.Arn)
		if !IsDatadogLayer(arn, region) {
			targetLayers = append(targetLayers, arn)
		}
	}

	// Step 3: Configure environment variables
	envVars := make(map[string]string)
	if funcConfig.Configuration.Environment != nil {
		for k, v := range funcConfig.Configuration.Environment.Variables {
			envVars[k] = v
		}
	}

	envVars["DD_API_KEY"] = cfg.DDApiKey
	envVars["DD_SITE"] = cfg.DDSite
	envVars["DD_TRACE_ENABLED"] = "true"
	envVars["DD_LOGS_ENABLED"] = "true"
	envVars["DD_MERGE_XRAY_TRACES"] = "true"

	originalHandler := aws.ToString(funcConfig.Configuration.Handler)
	runtimeHandler := GetHandlerForRuntime(runtime)

	updateInput := &lambdasvc.UpdateFunctionConfigurationInput{
		FunctionName: aws.String(functionARN),
		Layers:       targetLayers,
	}

	if UsesExecWrapper(runtime) {
		envVars["AWS_LAMBDA_EXEC_WRAPPER"] = "/opt/datadog_wrapper"
	} else if runtimeHandler != "" {
		// Wrap the handler directly (Node.js/Python)
		if originalHandler != runtimeHandler {
			if !isDatadogHandler(runtime, originalHandler) {
				envVars["DD_LAMBDA_HANDLER"] = originalHandler
			}
		}
		updateInput.Handler = aws.String(runtimeHandler)
	}

	updateInput.Environment = &lambdatypes.Environment{
		Variables: envVars,
	}

	// Step 4: Apply updates
	_, err = lambdaClient.UpdateFunctionConfiguration(ctx, updateInput)
	if err != nil {
		return fmt.Errorf("failed to update function configuration for %s: %w", functionARN, err)
	}

	return nil
}

// UninstallLayer removes Datadog instrumentation from a Lambda function.
func UninstallLayer(
	ctx context.Context,
	clients *awsclient.Factory,
	functionARN string,
) error {
	targetName := functionARN
	if parts := strings.Split(functionARN, ":"); len(parts) >= 7 {
		targetName = parts[6]
	}
	if isInternalFunction(targetName) {
		return fmt.Errorf("cannot uninstrument the orchestrator or internal helper function (%s)", targetName)
	}

	lambdaClient, err := clients.Lambda(ctx)
	if err != nil {
		return fmt.Errorf("failed to create Lambda client: %w", err)
	}

	region, err := clients.Region(ctx)
	if err != nil {
		return fmt.Errorf("failed to get region: %w", err)
	}

	// Step 1: Get current configuration
	funcConfig, err := lambdaClient.GetFunction(ctx, &lambdasvc.GetFunctionInput{
		FunctionName: aws.String(functionARN),
	})
	if err != nil {
		return fmt.Errorf("failed to get function %s: %w", functionARN, err)
	}

	// Step 2: Remove Datadog layers, preserving other layers
	remainingLayers := []string{}
	for _, layer := range funcConfig.Configuration.Layers {
		arn := aws.ToString(layer.Arn)
		if !IsDatadogLayer(arn, region) {
			remainingLayers = append(remainingLayers, arn)
		}
	}

	// Step 3: Restore handler and strip environment variables
	envVars := make(map[string]string)
	if funcConfig.Configuration.Environment != nil {
		for k, v := range funcConfig.Configuration.Environment.Variables {
			envVars[k] = v
		}
	}

	originalHandler := envVars["DD_LAMBDA_HANDLER"]
	if originalHandler == "" {
		originalHandler = aws.ToString(funcConfig.Configuration.Handler)
	}

	for _, key := range DatadogEnvVars {
		delete(envVars, key)
	}

	// Step 4: Apply changes
	updateInput := &lambdasvc.UpdateFunctionConfigurationInput{
		FunctionName: aws.String(functionARN),
		Handler:      aws.String(originalHandler),
		Layers:       remainingLayers,
		Environment: &lambdatypes.Environment{
			Variables: envVars,
		},
	}

	_, err = lambdaClient.UpdateFunctionConfiguration(ctx, updateInput)
	if err != nil {
		return fmt.Errorf("failed to restore function configuration for %s: %w", functionARN, err)
	}

	return nil
}

// isDatadogHandler checks if handler is already wrapped.
func isDatadogHandler(runtime, handler string) bool {
	runtimeHandler := GetHandlerForRuntime(runtime)
	if runtimeHandler == "" {
		return false
	}
	return handler == runtimeHandler
}

// isInternalFunction checks if the function belongs to the orchestrator or platform.
func isInternalFunction(name string) bool {
	nameLower := strings.ToLower(name)
	selfName := strings.ToLower(os.Getenv("AWS_LAMBDA_FUNCTION_NAME"))
	if selfName != "" && nameLower == selfName {
		return true
	}
	return strings.Contains(nameLower, "serverless-orchestrator")
}
