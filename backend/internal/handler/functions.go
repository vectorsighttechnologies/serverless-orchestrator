package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/vectorsight/serverless-tool/backend/internal/auth"
	"github.com/vectorsight/serverless-tool/backend/internal/db"
	"github.com/vectorsight/serverless-tool/backend/internal/lambda"
	"github.com/vectorsight/serverless-tool/backend/internal/types"
)

// HandleListFunctions retrieves Lambda functions, checking in-memory cache first.
func HandleListFunctions(
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

		// Retrieve credentials
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

		// Cache lookup
		cacheKey := "functions-" + userID
		if cachedData, found := memCache.Get(cacheKey); found {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Cache", "HIT")
			_, _ = w.Write(cachedData.([]byte))
			return
		}

		// Call orchestrator Lambda
		respBytes, statusCode, err := lambdaClient.Invoke(r.Context(), http.MethodGet, "/functions", nil, lambdaUrl, lambdaApiKey, prefs)
		if err != nil {
			// Try to unpack error response if available
			if len(respBytes) > 0 {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(statusCode)
				_, _ = w.Write(respBytes)
				return
			}
			respondWithError(w, statusCode, "Failed to fetch functions from orchestrator: "+err.Error())
			return
		}

		// Set cache (30 second TTL)
		memCache.Set(cacheKey, respBytes, 30*time.Second)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Cache", "MISS")
		w.WriteHeader(statusCode)
		_, _ = w.Write(respBytes)
	}
}

// HandleInstallFunctions handles proxying New Relic instrumentation to the orchestrator.
func HandleInstallFunctions(
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

		// Retrieve credentials
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

		// Parse request
		var req types.InstallRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		if len(req.FunctionArns) == 0 {
			respondWithError(w, http.StatusBadRequest, "At least one function ARN is required")
			return
		}

		bodyBytes, err := json.Marshal(req)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to serialize request")
			return
		}

		// Call orchestrator Lambda
		respBytes, statusCode, err := lambdaClient.Invoke(r.Context(), http.MethodPost, "/functions/install", bodyBytes, lambdaUrl, lambdaApiKey, prefs)
		if err != nil {
			_ = dbClient.CreateAuditLog(userID, "install", fmt.Sprintf("%d functions", len(req.FunctionArns)), "failure", err.Error())
			if len(respBytes) > 0 {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(statusCode)
				_, _ = w.Write(respBytes)
				return
			}
			respondWithError(w, statusCode, "Failed to install layer: "+err.Error())
			return
		}

		// Invalidate cache
		memCache.Delete("functions-" + userID)

		// Audit Log
		_ = dbClient.CreateAuditLog(userID, "install", fmt.Sprintf("%d functions", len(req.FunctionArns)), "success", "")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_, _ = w.Write(respBytes)
	}
}

// HandleUninstallFunctions handles proxying New Relic uninstrumentation to the orchestrator.
func HandleUninstallFunctions(
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

		// Retrieve credentials
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

		// Parse request
		var req types.UninstallRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		if len(req.FunctionArns) == 0 {
			respondWithError(w, http.StatusBadRequest, "At least one function ARN is required")
			return
		}

		bodyBytes, err := json.Marshal(req)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to serialize request")
			return
		}

		// Call orchestrator Lambda
		respBytes, statusCode, err := lambdaClient.Invoke(r.Context(), http.MethodPost, "/functions/uninstall", bodyBytes, lambdaUrl, lambdaApiKey, prefs)
		if err != nil {
			_ = dbClient.CreateAuditLog(userID, "uninstall", fmt.Sprintf("%d functions", len(req.FunctionArns)), "failure", err.Error())
			if len(respBytes) > 0 {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(statusCode)
				_, _ = w.Write(respBytes)
				return
			}
			respondWithError(w, statusCode, "Failed to uninstall layer: "+err.Error())
			return
		}

		// Invalidate cache
		memCache.Delete("functions-" + userID)

		// Audit Log
		_ = dbClient.CreateAuditLog(userID, "uninstall", fmt.Sprintf("%d functions", len(req.FunctionArns)), "success", "")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_, _ = w.Write(respBytes)
	}
}
