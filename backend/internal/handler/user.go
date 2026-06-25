package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/vectorsighttechnologies/serverless-orchestrator/backend/internal/auth"
	"github.com/vectorsighttechnologies/serverless-orchestrator/backend/internal/db"
	"github.com/vectorsighttechnologies/serverless-orchestrator/backend/internal/types"
)

type preferencesResponse struct {
	SelectedProvider string `json:"selectedProvider"`
	NRAccountID      string `json:"nrAccountId"`
	NRRegion         string `json:"nrRegion"`
	LambdaAPIURL     string `json:"lambdaApiUrl"`
	HasNRApiKey      bool   `json:"hasNrApiKey"`
	HasNRLicenseKey  bool   `json:"hasNrLicenseKey"`
	HasLambdaAPIKey  bool   `json:"hasLambdaApiKey"`
}

// HandleGetPreferences returns a user's preferences without exposing raw secrets.
func HandleGetPreferences(dbClient *db.DB) http.HandlerFunc {
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

		prefs, err := dbClient.GetUserPreferences(userID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to load preferences: "+err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(preferencesResponse{
			SelectedProvider: "newrelic",
			NRAccountID:      prefs.NRAccountID,
			NRRegion:         prefs.NRRegion,
			LambdaAPIURL:     prefs.LambdaAPIURL,
			HasNRApiKey:      prefs.NRApiKey != "",
			HasNRLicenseKey:  prefs.NRLicenseKey != "",
			HasLambdaAPIKey:  prefs.LambdaAPIKey != "",
		})
	}
}

// HandleSavePreferences saves a user's preferences, preserving existing credentials if not updated.
func HandleSavePreferences(dbClient *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut && r.Method != http.MethodPost {
			respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		userID, ok := auth.GetUserIDFromContext(r.Context())
		if !ok {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized context")
			return
		}

		var incoming types.UserPreferences
		if err := json.NewDecoder(r.Body).Decode(&incoming); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Clean inputs
		incoming.LambdaAPIURL = strings.TrimSpace(incoming.LambdaAPIURL)
		incoming.LambdaAPIKey = strings.TrimSpace(incoming.LambdaAPIKey)
		incoming.NRAccountID = strings.TrimSpace(incoming.NRAccountID)
		incoming.NRApiKey = strings.TrimSpace(incoming.NRApiKey)
		incoming.NRLicenseKey = strings.TrimSpace(incoming.NRLicenseKey)
		incoming.NRRegion = strings.ToLower(strings.TrimSpace(incoming.NRRegion))

		if incoming.NRRegion == "" {
			incoming.NRRegion = "us"
		}

		// Retrieve existing preferences to prevent wiping out unprovided secrets
		existing, err := dbClient.GetUserPreferences(userID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to load current preferences: "+err.Error())
			return
		}

		// Standard secure API design: if incoming credential is blank, preserve existing.
		if incoming.NRApiKey == "" {
			incoming.NRApiKey = existing.NRApiKey
		}
		if incoming.NRLicenseKey == "" {
			incoming.NRLicenseKey = existing.NRLicenseKey
		}
		if incoming.LambdaAPIKey == "" {
			incoming.LambdaAPIKey = existing.LambdaAPIKey
		}

		// Save preferences (encodes and encrypts sensitive fields inside SaveUserPreferences)
		if err := dbClient.SaveUserPreferences(userID, &incoming); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to save preferences: "+err.Error())
			return
		}

		// Audit Log
		_ = dbClient.CreateAuditLog(userID, "update_preferences", "preferences", "success", "Saved user configuration preferences")

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "Preferences updated successfully"})
	}
}

// respondWithError formats and returns a standard JSON error response.
func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	errBody := types.ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	}
	_ = json.NewEncoder(w).Encode(errBody)
}

// HandleGetConnections returns all connections for the user.
func HandleGetConnections(dbClient *db.DB) http.HandlerFunc {
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

		conns, err := dbClient.GetUserConnections(userID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to load connections: "+err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(conns)
	}
}

// HandleSaveConnection saves or updates a connection.
func HandleSaveConnection(dbClient *db.DB) http.HandlerFunc {
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

		var conn types.UserConnection
		if err := json.NewDecoder(r.Body).Decode(&conn); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		conn.Name = strings.TrimSpace(conn.Name)
		conn.AWSRegion = strings.TrimSpace(conn.AWSRegion)
		conn.LambdaAPIURL = strings.TrimSpace(conn.LambdaAPIURL)
		conn.LambdaAPIKey = strings.TrimSpace(conn.LambdaAPIKey)

		if conn.Name == "" || conn.AWSRegion == "" || conn.LambdaAPIURL == "" {
			respondWithError(w, http.StatusBadRequest, "Connection name, region, and URL are required")
			return
		}

		id, err := dbClient.SaveUserConnection(userID, &conn)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to save connection: "+err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"id":      id,
			"message": "Connection saved successfully",
		})
	}
}

// HandleDeleteConnection deletes a connection.
func HandleDeleteConnection(dbClient *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		userID, ok := auth.GetUserIDFromContext(r.Context())
		if !ok {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized context")
			return
		}

		// Extract connection ID from URL path (e.g., /api/user/connections/uuid)
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 5 || parts[4] == "" {
			respondWithError(w, http.StatusBadRequest, "Connection ID is required")
			return
		}
		connID := parts[4]

		err := dbClient.DeleteUserConnection(userID, connID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to delete connection: "+err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "Connection deleted successfully"})
	}
}

// resolveConnection extracts X-Connection-ID, fetches connection from DB, and falls back to legacy preferences.
func resolveConnection(r *http.Request, dbClient *db.DB, userID string, prefs *types.UserPreferences) (string, string, error) {
	connID := r.Header.Get("X-Connection-ID")
	if connID != "" {
		conn, err := dbClient.GetUserConnection(userID, connID)
		if err != nil {
			return "", "", fmt.Errorf("invalid AWS connection ID: %w", err)
		}
		return conn.LambdaAPIURL, conn.LambdaAPIKey, nil
	}

	// Fallback to legacy single preferences
	if prefs.LambdaAPIURL != "" {
		return prefs.LambdaAPIURL, prefs.LambdaAPIKey, nil
	}

	// Try loading first connection as fallback
	conns, err := dbClient.GetUserConnections(userID)
	if err == nil && len(conns) > 0 {
		return conns[0].LambdaAPIURL, conns[0].LambdaAPIKey, nil
	}

	return "", "", fmt.Errorf("no AWS connections configured. Please setup an AWS connection in Settings first")
}


