// Package newrelic contains the core business logic for Lambda instrumentation.
//
// This file defines the runtime configuration map and environment variable constants.


//
// DRY: All runtime metadata lives here. Both install and uninstall
// operations reference this single source of truth.
package newrelic

import "github.com/vectorsighttechnologies/serverless-orchestrator/lambda/internal/types"

// NR Layer ARN prefix template. The account 451483290750 is New Relic's
// official layer publishing account.
const LayerARNPrefix = "arn:aws:lambda:%s:451483290750"

// LayerRegistryURL is the endpoint for discovering NR Lambda layers.
const LayerRegistryURL = "https://%s.layers.newrelic-external.com/get-layers?CompatibleRuntime=%s"

// NREnvVars are the environment variables managed by New Relic instrumentation.
// During uninstall, only these keys are removed — preserving customer env vars.
var NREnvVars = []string{
	"NEW_RELIC_ACCOUNT_ID",
	"NEW_RELIC_EXTENSION_SEND_EXTENSION_LOGS",
	"NEW_RELIC_EXTENSION_SEND_FUNCTION_LOGS",
	"NEW_RELIC_LAMBDA_EXTENSION_ENABLED",
	"NEW_RELIC_LAMBDA_HANDLER",
	"NEW_RELIC_LICENSE_KEY",
	"NEW_RELIC_LOG_ENDPOINT",
	"NEW_RELIC_TELEMETRY_ENDPOINT",
	"NEW_RELIC_APM_LAMBDA_MODE",
	"NR_TAGS",
	"NR_ENV_DELIMITER",
}

// DotNetEnvVars are the extra env vars set for .NET runtimes.
var DotNetEnvVars = map[string]string{
	"CORECLR_ENABLE_PROFILING": "1",
	"CORECLR_PROFILER":         "{36032161-FFC0-4B61-B559-F6C5D41BAE5A}",
	"CORECLR_NEWRELIC_HOME":    "/opt/lib/newrelic-dotnet-agent",
	"CORECLR_PROFILER_PATH":    "/opt/lib/newrelic-dotnet-agent/libNewRelicProfiler.so",
}

// DotNetEnvVarKeys lists the extra keys to remove on uninstall for .NET.
var DotNetEnvVarKeys = []string{
	"CORECLR_ENABLE_PROFILING",
	"CORECLR_PROFILER",
	"CORECLR_NEWRELIC_HOME",
	"CORECLR_PROFILER_PATH",
}

// RuntimeConfig maps Lambda runtimes to their NR handler wrapper and extension support.
//
// Rules:
//   - Handler empty → no handler wrapping needed (e.g., .NET, provided runtimes)
//   - Java runtimes use "com.newrelic.java.HandlerWrapper::" + method name
//   - Node.js / Python / Ruby use their respective wrapper handlers
var RuntimeConfig = map[string]types.RuntimeInfo{
	// .NET
	"dotnetcore3.1": {Handler: "", HasExtension: true},
	"dotnet6":       {Handler: "", HasExtension: true},
	"dotnet8":       {Handler: "", HasExtension: true},

	// Java
	"java8.al2": {Handler: "com.newrelic.java.HandlerWrapper::", HasExtension: true},
	"java11":    {Handler: "com.newrelic.java.HandlerWrapper::", HasExtension: true},
	"java17":    {Handler: "com.newrelic.java.HandlerWrapper::", HasExtension: true},
	"java21":    {Handler: "com.newrelic.java.HandlerWrapper::", HasExtension: true},

	// Node.js
	"nodejs16.x": {Handler: "newrelic-lambda-wrapper.handler", HasExtension: true},
	"nodejs18.x": {Handler: "newrelic-lambda-wrapper.handler", HasExtension: true},
	"nodejs20.x": {Handler: "newrelic-lambda-wrapper.handler", HasExtension: true},
	"nodejs22.x": {Handler: "newrelic-lambda-wrapper.handler", HasExtension: true},
	"nodejs24.x": {Handler: "newrelic-lambda-wrapper.handler", HasExtension: true},

	// Custom / provided runtimes (extension only, no handler wrapping)
	"provided":        {Handler: "", HasExtension: true},
	"provided.al2":    {Handler: "", HasExtension: true},
	"provided.al2023": {Handler: "", HasExtension: true},

	// Python
	"python3.7":  {Handler: "newrelic_lambda_wrapper.handler", HasExtension: true},
	"python3.8":  {Handler: "newrelic_lambda_wrapper.handler", HasExtension: true},
	"python3.9":  {Handler: "newrelic_lambda_wrapper.handler", HasExtension: true},
	"python3.10": {Handler: "newrelic_lambda_wrapper.handler", HasExtension: true},
	"python3.11": {Handler: "newrelic_lambda_wrapper.handler", HasExtension: true},
	"python3.12": {Handler: "newrelic_lambda_wrapper.handler", HasExtension: true},
	"python3.13": {Handler: "newrelic_lambda_wrapper.handler", HasExtension: true},
	"python3.14": {Handler: "newrelic_lambda_wrapper.handler", HasExtension: true},

	// Ruby
	"ruby3.2": {Handler: "newrelic_lambda_wrapper.handler", HasExtension: true},
	"ruby3.3": {Handler: "newrelic_lambda_wrapper.handler", HasExtension: true},
	"ruby3.4": {Handler: "newrelic_lambda_wrapper.handler", HasExtension: true},
}

// IsSupportedRuntime checks if a runtime is supported for NR instrumentation.
func IsSupportedRuntime(runtime string) bool {
	_, ok := RuntimeConfig[runtime]
	return ok
}

// GetHandlerForRuntime returns the NR wrapper handler for the given runtime.
func GetHandlerForRuntime(runtime string) string {
	info, ok := RuntimeConfig[runtime]
	if !ok {
		return ""
	}
	return info.Handler
}
