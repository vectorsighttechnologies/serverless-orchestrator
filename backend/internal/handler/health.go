package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/vectorsight/serverless-tool/backend/internal/db"
	"github.com/vectorsight/serverless-tool/backend/internal/types"
)

// HandleHealth check checks the database ping and config availability.
func HandleHealth(dbClient *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dbStatus := "connected"
		if err := dbClient.SQLDB.Ping(); err != nil {
			dbStatus = "error: " + err.Error()
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(types.HealthResponse{
			Status:    "ok",
			Timestamp: time.Now().Format(time.RFC3339),
			Database:  dbStatus,
		})
	}
}
