package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/vectorsighttechnologies/serverless-orchestrator/backend/internal/auth"
	"github.com/vectorsighttechnologies/serverless-orchestrator/backend/internal/db"
	"github.com/vectorsighttechnologies/serverless-orchestrator/backend/internal/lambda"
	"github.com/vectorsighttechnologies/serverless-orchestrator/backend/internal/types"
)

// HandleIntegrationStatus retrieves current integration status from orchestrator, using TTL cache.
func HandleIntegrationStatus(
	dbClient *db.DB,
	lambdaClient *lambda.Client,
	memCache *cache.Cache,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		userID, ok := auth.GetUserIDFromContext(r.Context())
		if !ok {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized context")
			return
		}

		// Retrieve user preferences
		prefs, err := dbClient.GetUserPreferences(userID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to load preferences: "+err.Error())
			return
		}

		lambdaUrl, lambdaApiKey, err := resolveConnection(r, dbClient, userID, prefs)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Check cache (60s TTL for integration status)
		cacheKey := "integration-status-" + userID + "-" + lambdaUrl
		if cachedData, found := memCache.Get(cacheKey); found {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Cache", "HIT")
			_, _ = w.Write(cachedData.([]byte))
			return
		}

		// Call orchestrator Lambda
		respBytes, statusCode, err := lambdaClient.Invoke(r.Context(), http.MethodGet, "/integration/status", nil, lambdaUrl, lambdaApiKey, prefs)
		if err != nil {
			if len(respBytes) > 0 {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(statusCode)
				_, _ = w.Write(respBytes)
				return
			}
			respondWithError(w, statusCode, "Failed to fetch integration status: "+err.Error())
			return
		}

		// Cache response
		memCache.Set(cacheKey, respBytes, 60*time.Second)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Cache", "MISS")
		w.WriteHeader(statusCode)
		_, _ = w.Write(respBytes)
	}
}

// HandleIntegrationSetup deploys the CloudFormation integration stack on AWS.
func HandleIntegrationSetup(
	dbClient *db.DB,
	lambdaClient *lambda.Client,
	memCache *cache.Cache,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		userID, ok := auth.GetUserIDFromContext(r.Context())
		if !ok {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized context")
			return
		}

		// Retrieve user preferences
		prefs, err := dbClient.GetUserPreferences(userID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to load preferences: "+err.Error())
			return
		}

		lambdaUrl, lambdaApiKey, err := resolveConnection(r, dbClient, userID, prefs)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Parse request body
		var req types.IntegrationSetupRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		if req.Method == "" {
			respondWithError(w, http.StatusBadRequest, "Integration method is required")
			return
		}

		bodyBytes, err := json.Marshal(req)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to serialize request")
			return
		}

		// Call orchestrator Lambda
		respBytes, statusCode, err := lambdaClient.Invoke(r.Context(), http.MethodPost, "/integration/setup", bodyBytes, lambdaUrl, lambdaApiKey, prefs)
		if err != nil {
			_ = dbClient.CreateAuditLog(userID, "setup_integration", req.Method, "failure", err.Error())
			if len(respBytes) > 0 {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(statusCode)
				_, _ = w.Write(respBytes)
				return
			}
			respondWithError(w, statusCode, "Failed to setup integration: "+err.Error())
			return
		}

		// Invalidate status cache
		memCache.Delete("integration-status-" + userID)

		// Audit Log
		_ = dbClient.CreateAuditLog(userID, "setup_integration", req.Method, "success", fmt.Sprintf("Setup initiated: logs=%v", req.IncludeLogs))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_, _ = w.Write(respBytes)
	}
}

// HandleIntegrationRemove deletes the CloudFormation integration stack on AWS.
func HandleIntegrationRemove(
	dbClient *db.DB,
	lambdaClient *lambda.Client,
	memCache *cache.Cache,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		userID, ok := auth.GetUserIDFromContext(r.Context())
		if !ok {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized context")
			return
		}

		// Retrieve user preferences
		prefs, err := dbClient.GetUserPreferences(userID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to load preferences: "+err.Error())
			return
		}

		lambdaUrl, lambdaApiKey, err := resolveConnection(r, dbClient, userID, prefs)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Call orchestrator Lambda
		respBytes, statusCode, err := lambdaClient.Invoke(r.Context(), http.MethodPost, "/integration/remove", nil, lambdaUrl, lambdaApiKey, prefs)
		if err != nil {
			_ = dbClient.CreateAuditLog(userID, "remove_integration", "stack", "failure", err.Error())
			if len(respBytes) > 0 {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(statusCode)
				_, _ = w.Write(respBytes)
				return
			}
			respondWithError(w, statusCode, "Failed to remove integration: "+err.Error())
			return
		}

		// Invalidate status cache
		memCache.Delete("integration-status-" + userID)

		// Audit Log
		_ = dbClient.CreateAuditLog(userID, "remove_integration", "stack", "success", "Removal initiated")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_, _ = w.Write(respBytes)
	}
}
