package main

import (
	"log"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/vectorsighttechnologies/serverless-orchestrator/backend/internal/auth"
	"github.com/vectorsighttechnologies/serverless-orchestrator/backend/internal/config"
	"github.com/vectorsighttechnologies/serverless-orchestrator/backend/internal/db"
	"github.com/vectorsighttechnologies/serverless-orchestrator/backend/internal/handler"
	"github.com/vectorsighttechnologies/serverless-orchestrator/backend/internal/lambda"
	"github.com/vectorsighttechnologies/serverless-orchestrator/backend/internal/middleware"
)

func main() {
	log.Println("Starting Lambda Monitor Backend Gateway...")

	// 1. Load Configurations
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Configuration load failure: %v", err)
	}

	// 2. Initialize DB Client
	dbClient, err := db.NewDB(cfg.DBDriver, cfg.DatabaseURL, cfg.EncryptionKey)
	if err != nil {
		log.Fatalf("Database initialization failure: %v", err)
	}
	defer dbClient.Close()
	log.Printf("Database connection established using %s driver", cfg.DBDriver)

	// 3. Initialize In-Memory TTL Cache (5m default expiration, 10m cleanup interval)
	memCache := cache.New(5*time.Minute, 10*time.Minute)

	// 4. Initialize API Gateway Lambda Invoker
	lambdaClient := lambda.NewClient()

	// 5. Create ServeMux Router
	mux := http.NewServeMux()

	// 6. Setup Auth Middleware wrapper
	authMiddleware := auth.AuthMiddleware(cfg.JWTSecret)

	// ─── Public Endpoints ───
	mux.HandleFunc("/api/health", handler.HandleHealth(dbClient))
	mux.HandleFunc("/api/auth/register", auth.HandleRegister(dbClient))
	mux.HandleFunc("/api/auth/login", auth.HandleLogin(dbClient, cfg.JWTSecret))
	mux.HandleFunc("/api/auth/refresh", auth.HandleRefresh(cfg.JWTSecret))

	// ─── Protected Preferences Endpoints ───
	mux.Handle("/api/user/preferences", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.HandleGetPreferences(dbClient)(w, r)
		case http.MethodPut, http.MethodPost:
			handler.HandleSavePreferences(dbClient)(w, r)
		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = w.Write([]byte(`{"error":"Method Not Allowed"}`))
		}
	})))

	// ─── Protected Connections Endpoints ───
	mux.Handle("/api/user/connections", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.HandleGetConnections(dbClient)(w, r)
		case http.MethodPost:
			handler.HandleSaveConnection(dbClient)(w, r)
		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = w.Write([]byte(`{"error":"Method Not Allowed"}`))
		}
	})))

	mux.Handle("/api/user/connections/", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			handler.HandleDeleteConnection(dbClient)(w, r)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = w.Write([]byte(`{"error":"Method Not Allowed"}`))
		}
	})))


	// ─── Protected Lambda Orchestrator Proxy Endpoints ───
	mux.Handle("/api/functions", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.HandleListFunctions(dbClient, lambdaClient, memCache)(w, r)
	})))

	mux.Handle("/api/functions/install", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.HandleInstallFunctions(dbClient, lambdaClient, memCache)(w, r)
	})))

	mux.Handle("/api/functions/uninstall", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.HandleUninstallFunctions(dbClient, lambdaClient, memCache)(w, r)
	})))

	mux.Handle("/api/integration/status", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.HandleIntegrationStatus(dbClient, lambdaClient, memCache)(w, r)
	})))

	mux.Handle("/api/integration/setup", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.HandleIntegrationSetup(dbClient, lambdaClient, memCache)(w, r)
	})))

	mux.Handle("/api/integration/remove", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.HandleIntegrationRemove(dbClient, lambdaClient, memCache)(w, r)
	})))

	// 7. Apply CORS Middleware globally and start Server
	serverAddress := ":" + cfg.Port
	log.Printf("Server listening on http://localhost%s", serverAddress)
	if err := http.ListenAndServe(serverAddress, middleware.CORS(mux)); err != nil {
		log.Fatalf("Server listener failed: %v", err)
	}
}
