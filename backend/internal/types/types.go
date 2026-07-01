package types

import "time"

// ─────────────────────────────────────────────────────────────
// Shared API Response Types
// ─────────────────────────────────────────────────────────────

// HealthResponse represents GET /api/health output.
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Database  string `json:"database_status"` // "connected" or "error"
}

// ErrorResponse is the standard shape returned for any HTTP error.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// ─────────────────────────────────────────────────────────────
// AWS Lambda Function Entities (Mirrors Lambda Orchestrator)
// ─────────────────────────────────────────────────────────────

// FunctionInfo describes a single Lambda function and its instrumentation state.
type FunctionInfo struct {
	Name         string            `json:"name"`
	Arn          string            `json:"arn"`
	Runtime      string            `json:"runtime"`
	Architecture string            `json:"architecture"`
	Handler      string            `json:"handler"`
	Status       string            `json:"status"` // "instrumented" | "not_instrumented"
	Mode         string            `json:"mode"`   // "serverless" | "apm" | "log_ingestion" | "none"
	LayerVersion string            `json:"layerVersion"`
	LastModified string            `json:"lastModified"`
	MemorySize   int32             `json:"memorySize"`
	Timeout      int32             `json:"timeout"`
	CodeSize     int64             `json:"codeSize"`
	Invocations  int64             `json:"invocations,omitempty"`
	ErrorRate    float64           `json:"errorRate,omitempty"`
	Tags         map[string]string `json:"tags,omitempty"`
}

// FunctionsResponse represents GET /functions payload.
type FunctionsResponse struct {
	Functions []FunctionInfo `json:"functions"`
}

// InstallRequest is the payload structure for POST /api/functions/install.
type InstallRequest struct {
	FunctionArns []string `json:"functionArns"`
	Method       string   `json:"method"` // "layer" | "log_ingestion"
	Mode         string   `json:"mode"`   // "serverless" | "apm"
}

// UninstallRequest is the payload structure for POST /api/functions/uninstall.
type UninstallRequest struct {
	FunctionArns []string `json:"functionArns"`
}

// OperationResult is the response returned for a single function operation.
type OperationResult struct {
	Arn          string `json:"arn"`
	FunctionName string `json:"functionName"`
	Success      bool   `json:"success"`
	Error        string `json:"error,omitempty"`
}

// BatchOperationResponse represents bulk action output.
type BatchOperationResponse struct {
	Results []OperationResult `json:"results"`
}

// ─────────────────────────────────────────────────────────────
// Integration Status
// ─────────────────────────────────────────────────────────────

// IntegrationStatusResponse maps CloudFormation status endpoints.
type IntegrationStatusResponse struct {
	Status      string `json:"status"` // "active" | "not_setup" | "in_progress" | "error"
	Method      string `json:"method,omitempty"`
	StackName   string `json:"stackName,omitempty"`
	StackStatus string `json:"stackStatus,omitempty"`
	LastChecked string `json:"lastChecked"`
	Error       string `json:"error,omitempty"`
}

// IntegrationSetupRequest defines settings to build Kinesis/CloudWatch integrations.
type IntegrationSetupRequest struct {
	Method      string `json:"method"` // "metric_streams" | "api_polling"
	IncludeLogs bool   `json:"includeLogs"`
}

// IntegrationSetupResponse tracks status of CloudFormation launch.
type IntegrationSetupResponse struct {
	Status  string `json:"status"`
	StackID string `json:"stackId,omitempty"`
}

// ─────────────────────────────────────────────────────────────
// User Entities
// ─────────────────────────────────────────────────────────────

// UserPreferences holds credential and Lambda target configuration.
type UserPreferences struct {
	SelectedProvider string `json:"selectedProvider,omitempty"`
	NRAccountID      string `json:"nrAccountId,omitempty"`
	NRApiKey     string `json:"nrApiKey,omitempty"`
	NRLicenseKey string `json:"nrLicenseKey,omitempty"`
	NRRegion     string `json:"nrRegion,omitempty"`     // "us" or "eu"
	DDApiKey     string `json:"ddApiKey,omitempty"`
	DDSite       string `json:"ddSite,omitempty"`
	LambdaAPIURL string `json:"lambdaApiUrl,omitempty"` // URL of the deployed Lambda
	LambdaAPIKey string `json:"lambdaApiKey,omitempty"` // API Gateway API key
}

// AuditEntry tracks actions performed on the platform.
type AuditEntry struct {
	UserID    string    `json:"userId"`
	Action    string    `json:"action"` // "register", "login", "install", "uninstall", "update_preferences"
	Target    string    `json:"target"`
	Status    string    `json:"status"`
	Details   string    `json:"details"`
	Timestamp time.Time `json:"timestamp"`
}

// UserConnection holds credentials and target configuration for a single AWS orchestrator.
type UserConnection struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	AWSRegion    string `json:"awsRegion"`
	LambdaAPIURL string `json:"lambdaApiUrl"`
	LambdaAPIKey string `json:"lambdaApiKey,omitempty"`
	HasAPIKey    bool   `json:"hasApiKey"`
}

