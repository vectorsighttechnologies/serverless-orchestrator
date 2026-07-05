// Lambda Instrumentation Platform — Orchestrator
//
// Pure AWS orchestrator. This Lambda ONLY handles AWS SDK operations.
// Security is handled by API Gateway (API Key required).
// Auth, caching, and user state live in the Backend Gateway (Zoho AppSail).
//
// Architecture:
//
//	Backend Gateway → API Gateway (API Key) → Lambda (this) → AWS SDK
//
// Build:
//
//	GOOS=linux GOARCH=arm64 go build -o bootstrap ./cmd/lambda/
package main

import (
	"context"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/vectorsighttechnologies/serverless-orchestrator/lambda/internal/awsclient"
	"github.com/vectorsighttechnologies/serverless-orchestrator/lambda/internal/handler"
)

// clients is the shared AWS client factory — initialised once, reused across invocations.
var clients = awsclient.New()

func main() {
	lambda.Start(routeRequest)
}

// routeRequest routes API Gateway events to the appropriate handler.
// No CORS, no auth — API Gateway handles both.
func routeRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	path := normalisePath(request.Path)
	method := request.HTTPMethod

	switch {
	// ── Health ──
	case method == "GET" && path == "/health":
		return handler.HandleHealth(request), nil

	// ── Functions ──
	case method == "GET" && path == "/functions":
		return handler.HandleListFunctions(ctx, request, clients), nil

	case method == "POST" && path == "/functions/install":
		return handler.HandleInstallFunctions(ctx, request, clients), nil

	case method == "POST" && path == "/functions/uninstall":
		return handler.HandleUninstallFunctions(ctx, request, clients), nil

	// ── Integration ──
	case method == "GET" && path == "/integration/status":
		return handler.HandleIntegrationStatus(ctx, request, clients), nil

	case method == "POST" && path == "/integration/setup":
		return handler.HandleIntegrationSetup(ctx, request, clients), nil

	case method == "POST" && path == "/integration/remove":
		return handler.HandleIntegrationRemove(ctx, request, clients), nil

	// ── Not Found ──
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       `{"error":"Not Found","message":"Unknown endpoint: ` + method + " " + path + `"}`,
		}, nil
	}
}

// normalisePath cleans the request path for consistent routing.
func normalisePath(path string) string {
	for _, prefix := range []string{"/Prod", "/Stage", "/Dev"} {
		if strings.HasPrefix(path, prefix) {
			path = strings.TrimPrefix(path, prefix)
		}
	}
	path = strings.TrimRight(path, "/")
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return path
}
