package config

import (
	"fmt"
	"os"
)

// Config holds all backend configurations.
type Config struct {
	Port          string
	DBDriver      string // "sqlite" or "postgres"
	DatabaseURL   string // connection string or file path
	EncryptionKey string // 32-byte hex string
	JWTSecret     string // JWT secret key
}

// Load loads all settings from environment variables.
func Load() (*Config, error) {
	cfg := &Config{
		Port:          envOrDefault("X_ZOHO_CATALYST_LISTEN_PORT", "9000"),
		DBDriver:      envOrDefault("DB_DRIVER", "sqlite"),
		DatabaseURL:   envOrDefault("DATABASE_URL", "serverless_orchestrator.db"),
		EncryptionKey: os.Getenv("ENCRYPTION_KEY"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
	}

	// Set a fallback JWT secret for easier local development/open-source usage
	if cfg.JWTSecret == "" {
		cfg.JWTSecret = "development-fallback-secret-key-please-change-in-prod"
	}

	// Validate encryption key length or hex encoding if provided
	if cfg.EncryptionKey != "" && len(cfg.EncryptionKey) != 64 {
		return nil, fmt.Errorf("ENCRYPTION_KEY must be a 64-character hex string (32 bytes)")
	}

	// If ENCRYPTION_KEY is empty, we fall back to a derived key from JWTSecret
	if cfg.EncryptionKey == "" {
		// Use a 32-byte representation derived from JWT secret
		derived := make([]byte, 32)
		copy(derived, []byte(cfg.JWTSecret))
		// Format as a 64-character hex string for internal consistency
		cfg.EncryptionKey = fmt.Sprintf("%x", derived)
	}

	return cfg, nil
}

// envOrDefault reads an environment variable or returns the default value if empty.
func envOrDefault(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}
