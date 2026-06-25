// Package config implements flexible credential resolution.
//
// Priority order: Environment variables → Request headers → Defaults.
// This follows the Single Responsibility Principle — config loading
// is isolated from business logic and handlers.
package config

import (
	"errors"
	"os"
	"strings"
)

// Config holds the resolved New Relic credentials and settings.
type Config struct {
	LicenseKey string // NR Ingest License Key
	AccountID  string // NR Account ID
	APIKey     string // NR User API Key (NRAK-xxx)
	Region     string // "us" | "eu"
	Source     string // "env_vars" | "request" | "none"
}

// RequestHeaders abstracts the headers from an API Gateway event
// so this package doesn't depend on the Lambda runtime types.
type RequestHeaders map[string]string

// Load resolves configuration from env vars with fallback to request headers.
//
// A "source" field is added to help client applications identify how credentials
// were resolved (e.g. from environment variables vs request parameters) to customize UI display.
func Load(headers RequestHeaders) (*Config, error) {
	cfg := &Config{
		LicenseKey: os.Getenv("NEW_RELIC_LICENSE_KEY"),
		AccountID:  os.Getenv("NEW_RELIC_ACCOUNT_ID"),
		APIKey:     os.Getenv("NEW_RELIC_API_KEY"),
		Region:     os.Getenv("NEW_RELIC_REGION"),
		Source:     "env_vars",
	}

	// Fallback: read from request headers if env vars not set
	if cfg.LicenseKey == "" {
		cfg.LicenseKey = headers["x-nr-license-key"]
		if cfg.LicenseKey != "" {
			cfg.Source = "request"
		}
	}
	if cfg.AccountID == "" {
		cfg.AccountID = headers["x-nr-account-id"]
		if cfg.AccountID != "" {
			cfg.Source = "request"
		}
	}
	if cfg.APIKey == "" {
		cfg.APIKey = headers["x-nr-api-key"]
	}
	if cfg.Region == "" {
		cfg.Region = headers["x-nr-region"]
	}

	// Default region to US
	if cfg.Region == "" {
		cfg.Region = "us"
	}

	// Normalise to lowercase
	cfg.Region = strings.ToLower(cfg.Region)

	return cfg, nil
}

// Validate checks that the minimum required credentials are present.
func (c *Config) Validate() error {
	if c.LicenseKey == "" {
		return errors.New("NEW_RELIC_LICENSE_KEY not configured (set via env var or x-nr-license-key header)")
	}
	if c.AccountID == "" {
		return errors.New("NEW_RELIC_ACCOUNT_ID not configured (set via env var or x-nr-account-id header)")
	}
	return nil
}

// IsFullyConfigured returns true if all three keys are set via env vars.
func (c *Config) IsFullyConfigured() bool {
	return c.LicenseKey != "" && c.AccountID != "" && c.APIKey != ""
}

// HasLicenseKey returns true if a license key is available.
func (c *Config) HasLicenseKey() bool {
	return c.LicenseKey != ""
}

// HasAccountID returns true if an account ID is available.
func (c *Config) HasAccountID() bool {
	return c.AccountID != ""
}

// HasAPIKey returns true if an API key is available.
func (c *Config) HasAPIKey() bool {
	return c.APIKey != ""
}
