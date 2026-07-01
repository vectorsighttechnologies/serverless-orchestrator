// Package types defines shared domain types used across the backend.
// Following the DRY principle, all shared data structures live here
// to avoid duplication between handlers, services, and the API layer.
package types

// ─────────────────────────────────────────────────────────────
// API Request / Response Types
// ─────────────────────────────────────────────────────────────

// HealthResponse is returned by GET /health.
type HealthResponse struct {
	Status string       `json:"status"`
	Config ConfigStatus `json:"config"`
}

// ConfigStatus indicates which credentials are configured on the orchestrator.
type ConfigStatus struct {
	SelectedProvider     string `json:"selectedProvider,omitempty"`
	LicenseKeyConfigured bool   `json:"licenseKeyConfigured"`
	AccountIDConfigured  bool   `json:"accountIdConfigured"`
	APIKeyConfigured     bool   `json:"apiKeyConfigured"`
	DDApiKeyConfigured   bool   `json:"ddApiKeyConfigured"`
	DDSiteConfigured     bool   `json:"ddSiteConfigured"`
	Region               string `json:"region"`
	Source               string `json:"source"` // "env_vars" | "request" | "none"
}

// FunctionsResponse is returned by GET /functions.
type FunctionsResponse struct {
	Functions []FunctionInfo `json:"functions"`
}

// FunctionInfo describes a single Lambda function and its instrumentation state.
type FunctionInfo struct {
	Name         string `json:"name"`
	Arn          string `json:"arn"`
	Runtime      string `json:"runtime"`
	Architecture string `json:"architecture"`
	Handler      string `json:"handler"`
	Status       string `json:"status"` // "instrumented" | "not_instrumented"
	Mode         string `json:"mode"`   // "serverless" | "apm" | "log_ingestion" | "none"
	LayerVersion string `json:"layerVersion"`
	LastModified string `json:"lastModified"`
	MemorySize   int32  `json:"memorySize"`
	Timeout      int32  `json:"timeout"`
	CodeSize     int64  `json:"codeSize"`
}

// InstallRequest is the body for POST /functions/install.
type InstallRequest struct {
	FunctionArns []string `json:"functionArns"`
	Method       string   `json:"method"` // "layer" | "log_ingestion"
	Mode         string   `json:"mode"`   // "serverless" | "apm"
}

// UninstallRequest is the body for POST /functions/uninstall.
type UninstallRequest struct {
	FunctionArns []string `json:"functionArns"`
}

// OperationResult describes the outcome of an install/uninstall operation on one function.
type OperationResult struct {
	Arn          string `json:"arn"`
	FunctionName string `json:"functionName"`
	Success      bool   `json:"success"`
	Error        string `json:"error,omitempty"`
}

// BatchOperationResponse is returned by install/uninstall endpoints.
type BatchOperationResponse struct {
	Results []OperationResult `json:"results"`
}

// IntegrationStatusResponse is returned by GET /integration/status.
type IntegrationStatusResponse struct {
	Status      string `json:"status"` // "active" | "not_setup" | "in_progress" | "error"
	Method      string `json:"method,omitempty"`
	StackName   string `json:"stackName,omitempty"`
	StackStatus string `json:"stackStatus,omitempty"`
	LastChecked string `json:"lastChecked"`
	Error       string `json:"error,omitempty"`
}

// IntegrationSetupRequest is the body for POST /integration/setup.
type IntegrationSetupRequest struct {
	Method      string `json:"method"` // "metric_streams" | "api_polling"
	IncludeLogs bool   `json:"includeLogs"`
}

// IntegrationSetupResponse is returned by POST /integration/setup.
type IntegrationSetupResponse struct {
	Status  string `json:"status"`
	StackID string `json:"stackId,omitempty"`
}

// ErrorResponse is a generic API error response.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// ─────────────────────────────────────────────────────────────
// Internal Domain Types
// ─────────────────────────────────────────────────────────────

// LayerInfo describes a New Relic Lambda layer from the layer registry.
type LayerInfo struct {
	LayerName            string `json:"LayerName"`
	LatestMatchingVersion struct {
		LayerVersionArn         string   `json:"LayerVersionArn"`
		Version                 int64    `json:"Version"`
		CompatibleRuntimes      []string `json:"CompatibleRuntimes"`
		CompatibleArchitectures []string `json:"CompatibleArchitectures"`
	} `json:"LatestMatchingVersion"`
}

// LayerRegistryResponse is the response from layers.newrelic-external.com.
type LayerRegistryResponse struct {
	Layers []LayerInfo `json:"Layers"`
}

// RuntimeInfo holds the handler wrapping config for a specific Lambda runtime.
// Ported from Python CLI's utils.RUNTIME_CONFIG.
type RuntimeInfo struct {
	Handler          string // New handler value (empty = no wrapping needed)
	HasExtension     bool   // Whether the NR Lambda Extension is supported
}
