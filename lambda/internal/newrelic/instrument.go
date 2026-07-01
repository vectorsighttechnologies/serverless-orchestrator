// Package newrelic — instrument.go
//
// Layer-based instrumentation: install and uninstall.
// This is the heart of the platform — ported from:
//   - newrelic_lambda_cli/layers.py :: _add_new_relic(), install(), _remove_new_relic(), uninstall()
//
// SOLID principles applied:
//   - Single Responsibility: Only handles layer install/uninstall, not routing or config
//   - Open/Closed: New runtimes can be added to runtimes.go without modifying this file
//   - Dependency Inversion: Accepts AWS clients via interface (awsclient.Factory)
package newrelic

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	lambdasvc "github.com/aws/aws-sdk-go-v2/service/lambda"
	lambdatypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"

	"github.com/vectorsight/serverless-tool/lambda/internal/awsclient"
	"github.com/vectorsight/serverless-tool/lambda/internal/config"
)

// InstallLayer instruments a single Lambda function with the NR layer.
//
// Logic ported from: layers.py :: install() + _add_new_relic()
//
// Steps:
//  1. Get current function configuration
//  2. Discover compatible NR layer from registry
//  3. Build the UpdateFunctionConfiguration params:
//     - Add NR layer, preserve existing non-NR layers
//     - Set NR env vars (account ID, license key, extension enabled, etc.)
//     - Wrap the handler (runtime-specific)
//  4. Call UpdateFunctionConfiguration
//  5. If APM mode, tag the function
//  6. Attach Secrets Manager policy to the function's IAM role
func InstallLayer(
	ctx context.Context,
	clients *awsclient.Factory,
	cfg *config.Config,
	functionARN string,
	mode string, // "serverless" | "apm"
) error {
	// Prevent instrumenting the orchestrator or internal helper functions
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

	// Step 2: Discover compatible NR layer
	availableLayers, err := DiscoverLayers(region, runtime, architecture)
	if err != nil {
		return fmt.Errorf("failed to discover NR layers: %w", err)
	}

	// Find existing NR layer (for upgrade detection)
	var existingNRLayer string
	for _, layer := range funcConfig.Configuration.Layers {
		arn := aws.ToString(layer.Arn)
		if IsNRLayer(arn, region) {
			existingNRLayer = arn
			break
		}
	}

	nrLayerARN, err := SelectLayer(availableLayers, existingNRLayer, existingNRLayer != "")
	if err != nil {
		return err
	}

	// Step 3: Build update params
	// Preserve existing non-NR layers
	var existingLayers []string
	for _, layer := range funcConfig.Configuration.Layers {
		arn := aws.ToString(layer.Arn)
		if !IsNRLayer(arn, region) {
			existingLayers = append(existingLayers, arn)
		}
	}

	allLayers := append([]string{nrLayerARN}, existingLayers...)

	// Build env vars — start with existing
	envVars := make(map[string]string)
	if funcConfig.Configuration.Environment != nil {
		for k, v := range funcConfig.Configuration.Environment.Variables {
			envVars[k] = v
		}
	}

	// Set NR env vars
	envVars["NEW_RELIC_ACCOUNT_ID"] = cfg.AccountID
	envVars["NEW_RELIC_LICENSE_KEY"] = cfg.LicenseKey
	envVars["NEW_RELIC_LAMBDA_EXTENSION_ENABLED"] = "true"
	envVars["NEW_RELIC_EXTENSION_SEND_FUNCTION_LOGS"] = "false"

	if mode == "apm" {
		envVars["NEW_RELIC_APM_LAMBDA_MODE"] = "True"
	}

	// .NET runtimes need extra profiler env vars
	if strings.HasPrefix(runtime, "dotnet") {
		for k, v := range DotNetEnvVars {
			envVars[k] = v
		}
	}

	// Build the handler wrapping
	originalHandler := aws.ToString(funcConfig.Configuration.Handler)
	runtimeHandler := GetHandlerForRuntime(runtime)

	// For Java, append the default method name
	if strings.HasPrefix(runtime, "java") && runtimeHandler != "" {
		runtimeHandler += "handleRequest"
	}

	updateInput := &lambdasvc.UpdateFunctionConfigurationInput{
		FunctionName: aws.String(functionARN),
		Layers:       allLayers,
		Environment: &lambdatypes.Environment{
			Variables: envVars,
		},
	}

	// Wrap the handler (only if the runtime requires it AND not already wrapped)
	// Skip wrapping for NewRelicLambdaExtension-only layers
	isExtensionOnlyLayer := strings.Contains(nrLayerARN, "NewRelicLambdaExtension")

	if runtimeHandler != "" && !isExtensionOnlyLayer {
		// Save original handler so NR wrapper can call it
		if originalHandler != runtimeHandler {
			// Don't re-save if already a NR handler (upgrade case)
			if !isNRHandler(runtime, originalHandler) {
				envVars["NEW_RELIC_LAMBDA_HANDLER"] = originalHandler
			}
		}
		updateInput.Handler = aws.String(runtimeHandler)
	}

	// Step 4: Apply changes
	_, err = lambdaClient.UpdateFunctionConfiguration(ctx, updateInput)
	if err != nil {
		return fmt.Errorf("failed to update function configuration for %s: %w", functionARN, err)
	}

	// Step 5: APM mode — tag the function
	if mode == "apm" {
		_, err = lambdaClient.TagResource(ctx, &lambdasvc.TagResourceInput{
			Resource: aws.String(functionARN),
			Tags: map[string]string{
				"NR.Apm.Lambda.Mode": "true",
			},
		})
		if err != nil {
			// Non-fatal: tagging failure shouldn't block instrumentation
			fmt.Printf("warning: failed to add APM tag to %s: %v\n", functionARN, err)
		}
	}

	// Step 6: Attach Secrets Manager policy to function's execution role
	if funcConfig.Configuration.Role != nil {
		err = attachSecretsManagerPolicy(ctx, clients, aws.ToString(funcConfig.Configuration.Role))
		if err != nil {
			// Non-fatal
			fmt.Printf("warning: failed to attach SM policy to role: %v\n", err)
		}
	}

	return nil
}

// UninstallLayer removes NR instrumentation from a Lambda function.
//
// Logic ported from: layers.py :: uninstall() + _remove_new_relic()
//
// Steps:
//  1. Get current function config
//  2. Remove NR layers, restore non-NR layers
//  3. Restore original handler from NEW_RELIC_LAMBDA_HANDLER
//  4. Remove all NR env vars
//  5. Call UpdateFunctionConfiguration
//  6. Remove NR tags
//  7. Detach Secrets Manager policy
func UninstallLayer(
	ctx context.Context,
	clients *awsclient.Factory,
	functionARN string,
) error {
	// Prevent uninstrumenting the orchestrator or internal helper functions
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

	// Step 1: Get current config
	funcConfig, err := lambdaClient.GetFunction(ctx, &lambdasvc.GetFunctionInput{
		FunctionName: aws.String(functionARN),
	})
	if err != nil {
		return fmt.Errorf("failed to get function %s: %w", functionARN, err)
	}

	// Step 2: Remove NR layers
	nonNRLayers := []string{}
	for _, layer := range funcConfig.Configuration.Layers {
		arn := aws.ToString(layer.Arn)
		if !IsNRLayer(arn, region) {
			nonNRLayers = append(nonNRLayers, arn)
		}
	}

	// Step 3: Restore original handler
	envVars := make(map[string]string)
	if funcConfig.Configuration.Environment != nil {
		for k, v := range funcConfig.Configuration.Environment.Variables {
			envVars[k] = v
		}
	}

	originalHandler := envVars["NEW_RELIC_LAMBDA_HANDLER"]
	if originalHandler == "" {
		originalHandler = aws.ToString(funcConfig.Configuration.Handler)
	}

	// Step 4: Remove NR env vars
	for _, key := range NREnvVars {
		delete(envVars, key)
	}

	// Remove .NET-specific env vars if applicable
	runtime := string(funcConfig.Configuration.Runtime)
	if strings.HasPrefix(runtime, "dotnet") {
		for _, key := range DotNetEnvVarKeys {
			delete(envVars, key)
		}
	}

	// Step 5: Apply changes
	updateInput := &lambdasvc.UpdateFunctionConfigurationInput{
		FunctionName: aws.String(functionARN),
		Handler:      aws.String(originalHandler),
		Layers:       nonNRLayers,
		Environment: &lambdatypes.Environment{
			Variables: envVars,
		},
	}

	_, err = lambdaClient.UpdateFunctionConfiguration(ctx, updateInput)
	if err != nil {
		return fmt.Errorf("failed to restore function configuration for %s: %w", functionARN, err)
	}

	// Step 6: Remove NR tags
	_, err = lambdaClient.UntagResource(ctx, &lambdasvc.UntagResourceInput{
		Resource: aws.String(functionARN),
		TagKeys:  []string{"NR.Apm.Lambda.Mode"},
	})
	if err != nil {
		// Non-fatal
		fmt.Printf("warning: failed to remove NR tags from %s: %v\n", functionARN, err)
	}

	// Step 7: Detach Secrets Manager policy
	if funcConfig.Configuration.Role != nil {
		err = detachSecretsManagerPolicy(ctx, clients, aws.ToString(funcConfig.Configuration.Role))
		if err != nil {
			// Non-fatal
			fmt.Printf("warning: failed to detach SM policy from role: %v\n", err)
		}
	}

	return nil
}

// ─────────────────────────────────────────────────────────────
// Helper functions
// ─────────────────────────────────────────────────────────────

// isNRHandler checks if the handler is already a NR wrapper handler.
// Prevents double-wrapping during upgrades.
func isNRHandler(runtime, handler string) bool {
	if strings.HasPrefix(runtime, "nodejs") {
		return handler == "newrelic-lambda-wrapper.handler" ||
			handler == "/opt/nodejs/node_modules/newrelic-esm-lambda-wrapper/index.handler"
	}
	runtimeHandler := GetHandlerForRuntime(runtime)
	if runtimeHandler == "" {
		return false
	}
	// For Java, handler starts with the wrapper prefix
	if strings.HasPrefix(runtime, "java") {
		return strings.HasPrefix(handler, "com.newrelic.java.HandlerWrapper::")
	}
	return handler == runtimeHandler
}

// attachSecretsManagerPolicy attaches the AWSLambdaBasicExecutionRole
// policy for Secrets Manager access to the function's IAM role.
// Ported from: layers.py :: _attach_license_key_policy()
func attachSecretsManagerPolicy(ctx context.Context, clients *awsclient.Factory, roleARN string) error {
	iamClient, err := clients.IAM(ctx)
	if err != nil {
		return err
	}

	roleName := extractRoleName(roleARN)
	if roleName == "" {
		return fmt.Errorf("could not extract role name from ARN: %s", roleARN)
	}

	_, err = iamClient.AttachRolePolicy(ctx, &iam.AttachRolePolicyInput{
		RoleName:  aws.String(roleName),
		PolicyArn: aws.String("arn:aws:iam::aws:policy/SecretsManagerReadWrite"),
	})
	return err
}

// detachSecretsManagerPolicy removes the Secrets Manager policy from the role.
// Ported from: layers.py :: _detach_license_key_policy()
func detachSecretsManagerPolicy(ctx context.Context, clients *awsclient.Factory, roleARN string) error {
	iamClient, err := clients.IAM(ctx)
	if err != nil {
		return err
	}

	roleName := extractRoleName(roleARN)
	if roleName == "" {
		return fmt.Errorf("could not extract role name from ARN: %s", roleARN)
	}

	_, err = iamClient.DetachRolePolicy(ctx, &iam.DetachRolePolicyInput{
		RoleName:  aws.String(roleName),
		PolicyArn: aws.String("arn:aws:iam::aws:policy/SecretsManagerReadWrite"),
	})
	return err
}

// extractRoleName extracts the role name from an IAM role ARN.
// e.g., "arn:aws:iam::123456789012:role/my-role" → "my-role"
func extractRoleName(roleARN string) string {
	parts := strings.Split(roleARN, "/")
	if len(parts) < 2 {
		return ""
	}
	return parts[len(parts)-1]
}
