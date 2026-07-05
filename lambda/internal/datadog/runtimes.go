// Package datadog contains the core business logic for Datadog Lambda instrumentation.
package datadog

import (
	"strings"

	"github.com/vectorsighttechnologies/serverless-orchestrator/lambda/internal/types"
)

// Datadog Layer publishing account ID.
const LayerAccountID = "464622532012"
const GovCloudAccountID = "002406178527"

// DatadogEnvVars are the environment variables managed by Datadog instrumentation.
var DatadogEnvVars = []string{
	"DD_API_KEY",
	"DD_API_KEY_SECRET_ARN",
	"DD_SITE",
	"DD_LAMBDA_HANDLER",
	"DD_TRACE_ENABLED",
	"DD_LOGS_ENABLED",
	"DD_MERGE_XRAY_TRACES",
	"DD_FLUSH_TO_LOG",
	"DD_SERVICE",
	"DD_ENV",
	"DD_VERSION",
	"AWS_LAMBDA_EXEC_WRAPPER",
}

// DatadogRuntimeInfo holds Datadog specific runtime info.
type DatadogRuntimeInfo struct {
	types.RuntimeInfo
	LibraryName string // Name of the Datadog library layer (empty if none)
}

// RuntimeConfig maps Lambda runtimes to their Datadog handler wrapper, extension support, and library layer names.
var RuntimeConfig = map[string]DatadogRuntimeInfo{
	// Node.js
	"nodejs16.x": {RuntimeInfo: types.RuntimeInfo{Handler: "/opt/nodejs/node_modules/datadog-lambda-js/handler.handler", HasExtension: true}, LibraryName: "Datadog-Node16-x"},
	"nodejs18.x": {RuntimeInfo: types.RuntimeInfo{Handler: "/opt/nodejs/node_modules/datadog-lambda-js/handler.handler", HasExtension: true}, LibraryName: "Datadog-Node18-x"},
	"nodejs20.x": {RuntimeInfo: types.RuntimeInfo{Handler: "/opt/nodejs/node_modules/datadog-lambda-js/handler.handler", HasExtension: true}, LibraryName: "Datadog-Node20-x"},
	"nodejs22.x": {RuntimeInfo: types.RuntimeInfo{Handler: "/opt/nodejs/node_modules/datadog-lambda-js/handler.handler", HasExtension: true}, LibraryName: "Datadog-Node22-x"},
	"nodejs24.x": {RuntimeInfo: types.RuntimeInfo{Handler: "/opt/nodejs/node_modules/datadog-lambda-js/handler.handler", HasExtension: true}, LibraryName: "Datadog-Node24-x"},

	// Python
	"python3.7":  {RuntimeInfo: types.RuntimeInfo{Handler: "datadog_lambda.handler.handler", HasExtension: true}, LibraryName: "Datadog-Python37"},
	"python3.8":  {RuntimeInfo: types.RuntimeInfo{Handler: "datadog_lambda.handler.handler", HasExtension: true}, LibraryName: "Datadog-Python38"},
	"python3.9":  {RuntimeInfo: types.RuntimeInfo{Handler: "datadog_lambda.handler.handler", HasExtension: true}, LibraryName: "Datadog-Python39"},
	"python3.10": {RuntimeInfo: types.RuntimeInfo{Handler: "datadog_lambda.handler.handler", HasExtension: true}, LibraryName: "Datadog-Python310"},
	"python3.11": {RuntimeInfo: types.RuntimeInfo{Handler: "datadog_lambda.handler.handler", HasExtension: true}, LibraryName: "Datadog-Python311"},
	"python3.12": {RuntimeInfo: types.RuntimeInfo{Handler: "datadog_lambda.handler.handler", HasExtension: true}, LibraryName: "Datadog-Python312"},
	"python3.13": {RuntimeInfo: types.RuntimeInfo{Handler: "datadog_lambda.handler.handler", HasExtension: true}, LibraryName: "Datadog-Python313"},

	// Java (uses Exec Wrapper)
	"java8.al2": {RuntimeInfo: types.RuntimeInfo{Handler: "", HasExtension: true}, LibraryName: ""},
	"java11":    {RuntimeInfo: types.RuntimeInfo{Handler: "", HasExtension: true}, LibraryName: ""},
	"java17":    {RuntimeInfo: types.RuntimeInfo{Handler: "", HasExtension: true}, LibraryName: ""},
	"java21":    {RuntimeInfo: types.RuntimeInfo{Handler: "", HasExtension: true}, LibraryName: ""},

	// .NET (uses Exec Wrapper)
	"dotnet6": {RuntimeInfo: types.RuntimeInfo{Handler: "", HasExtension: true}, LibraryName: ""},
	"dotnet8": {RuntimeInfo: types.RuntimeInfo{Handler: "", HasExtension: true}, LibraryName: ""},

	// Ruby (uses Exec Wrapper)
	"ruby3.2": {RuntimeInfo: types.RuntimeInfo{Handler: "", HasExtension: true}, LibraryName: ""},
	"ruby3.3": {RuntimeInfo: types.RuntimeInfo{Handler: "", HasExtension: true}, LibraryName: ""},

	// Provided / custom (Extension only)
	"provided":        {RuntimeInfo: types.RuntimeInfo{Handler: "", HasExtension: true}, LibraryName: ""},
	"provided.al2":    {RuntimeInfo: types.RuntimeInfo{Handler: "", HasExtension: true}, LibraryName: ""},
	"provided.al2023": {RuntimeInfo: types.RuntimeInfo{Handler: "", HasExtension: true}, LibraryName: ""},
}

// Default layer versions
const (
	DefaultExtensionVersion = 64
	DefaultLibraryVersion   = 115
)

// IsSupportedRuntime checks if a runtime is supported for Datadog instrumentation.
func IsSupportedRuntime(runtime string) bool {
	_, ok := RuntimeConfig[runtime]
	return ok
}

// GetHandlerForRuntime returns the Datadog wrapper handler for Node.js / Python.
func GetHandlerForRuntime(runtime string) string {
	info, ok := RuntimeConfig[runtime]
	if !ok {
		return ""
	}
	return info.Handler
}

// GetLibraryNameForRuntime returns the name of the Datadog library layer for the runtime.
func GetLibraryNameForRuntime(runtime string) string {
	info, ok := RuntimeConfig[runtime]
	if !ok {
		return ""
	}
	return info.LibraryName
}

// UsesExecWrapper checks if runtime uses /opt/datadog_wrapper.
func UsesExecWrapper(runtime string) bool {
	return strings.HasPrefix(runtime, "java") || strings.HasPrefix(runtime, "dotnet") || strings.HasPrefix(runtime, "ruby")
}
