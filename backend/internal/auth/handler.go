package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/vectorsight/serverless-tool/backend/internal/db"
	"github.com/vectorsight/serverless-tool/backend/internal/types"
)

type credentialsRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Email        string `json:"email"`
}

type refreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

// HandleRegister returns a handler that registers a new user.
func HandleRegister(dbClient *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		var req credentialsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		email := strings.TrimSpace(req.Email)
		password := req.Password

		if email == "" || password == "" {
			respondWithError(w, http.StatusBadRequest, "Email and password are required")
			return
		}

		if len(password) < 6 {
			respondWithError(w, http.StatusBadRequest, "Password must be at least 6 characters long")
			return
		}

		// Hash password
		hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to hash password")
			return
		}

		// Save user to DB
		userID, err := dbClient.CreateUser(email, string(hashedBytes))
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE") || strings.Contains(err.Error(), "duplicate key") {
				respondWithError(w, http.StatusConflict, "User with this email already exists")
				return
			}
			respondWithError(w, http.StatusInternalServerError, "Failed to create user account: "+err.Error())
			return
		}

		// Create an empty preferences row for this user
		emptyPrefs := &types.UserPreferences{
			NRRegion: "us",
		}
		_ = dbClient.SaveUserPreferences(userID, emptyPrefs)

		// Audit Log
		_ = dbClient.CreateAuditLog(userID, "register", email, "success", "User registered successfully")

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "Registration successful"})
	}
}

// HandleLogin returns a handler that authenticates credentials and returns JWTs.
func HandleLogin(dbClient *db.DB, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		var req credentialsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		email := strings.TrimSpace(req.Email)
		password := req.Password

		if email == "" || password == "" {
			respondWithError(w, http.StatusBadRequest, "Email and password are required")
			return
		}

		// Retrieve user details
		userID, passwordHash, err := dbClient.GetUserByEmail(email)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}

		// Compare hashes
		if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}

		// Generate tokens
		access, refresh, err := GenerateTokens(userID, email, jwtSecret)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to generate tokens")
			return
		}

		// Audit Log
		_ = dbClient.CreateAuditLog(userID, "login", email, "success", "User logged in successfully")

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(loginResponse{
			AccessToken:  access,
			RefreshToken: refresh,
			Email:        email,
		})
	}
}

// HandleRefresh returns a handler that issues new token sets using refresh tokens.
func HandleRefresh(jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		var req refreshRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if req.RefreshToken == "" {
			respondWithError(w, http.StatusBadRequest, "Refresh token is required")
			return
		}

		// Validate refresh token
		claims, err := ValidateToken(req.RefreshToken, jwtSecret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid or expired refresh token")
			return
		}

		// Generate new access and refresh tokens
		access, refresh, err := GenerateTokens(claims.UserID, claims.Email, jwtSecret)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to generate tokens")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(loginResponse{
			AccessToken:  access,
			RefreshToken: refresh,
			Email:        claims.Email,
		})
	}
}
