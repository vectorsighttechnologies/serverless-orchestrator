// Package handler implements the API Gateway event handlers.
//
// Each handler follows the Single Responsibility Principle:
//   - Parse request → Call service → Format response
//   - No business logic lives here — it's all in the newrelic package
//
// Response helpers (jsonResponse, errorResponse) follow DRY.
package handler

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"

	"github.com/vectorsight/serverless-tool/lambda/internal/config"
	"github.com/vectorsight/serverless-tool/lambda/internal/types"
)

// jsonResponse creates a successful API Gateway response with JSON body.
// No CORS headers — API Gateway handles CORS configuration.
func jsonResponse(statusCode int, body interface{}) events.APIGatewayProxyResponse {
	data, err := json.Marshal(body)
	if err != nil {
		return errorResponse(500, "Failed to marshal response")
	}
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(data),
	}
}

// errorResponse creates an error API Gateway response.
func errorResponse(statusCode int, message string) events.APIGatewayProxyResponse {
	body := types.ErrorResponse{
		Error:   statusText(statusCode),
		Message: message,
	}
	data, _ := json.Marshal(body)
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(data),
	}
}

// statusText returns a human-readable status text for common HTTP codes.
func statusText(code int) string {
	switch code {
	case 400:
		return "Bad Request"
	case 401:
		return "Unauthorized"
	case 404:
		return "Not Found"
	case 500:
		return "Internal Server Error"
	default:
		return "Error"
	}
}

// loadConfig creates a config from the request headers.
func loadConfig(headers map[string]string) (*config.Config, error) {
	return config.Load(config.RequestHeaders(headers))
}
