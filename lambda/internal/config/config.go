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

// Config holds the resolved credentials and settings.
type Config struct {
	SelectedProvider string // Active provider ("newrelic" or "datadog")
	LicenseKey       string // NR Ingest License Key
	AccountID        string // NR Account ID
	APIKey           string // NR User API Key (NRAK-xxx)
	Region           string // "us" | "eu"
	DDApiKey         string // Datadog API Key
	DDSite           string // Datadog Site (e.g. datadoghq.com)
	Source           string // "env_vars" | "request" | "none"
}

// RequestHeaders abstracts the headers from an API Gateway event
// so this package doesn't depend on the Lambda runtime types.
type RequestHeaders map[string]string

// getHeader retrieves a header value case-insensitively.
func getHeader(headers RequestHeaders, key string) string {
	lowerKey := strings.ToLower(key)
	for k, v := range headers {
		if strings.ToLower(k) == lowerKey {
			return v
		}
	}
	return ""
}

// Load resolves configuration from env vars with fallback to request headers.
func Load(headers RequestHeaders) (*Config, error) {
	cfg := &Config{
		LicenseKey:       os.Getenv("NEW_RELIC_LICENSE_KEY"),
		AccountID:        os.Getenv("NEW_RELIC_ACCOUNT_ID"),
		APIKey:           os.Getenv("NEW_RELIC_API_KEY"),
		Region:           os.Getenv("NEW_RELIC_REGION"),
		SelectedProvider: os.Getenv("SELECTED_PROVIDER"),
		DDApiKey:         os.Getenv("DATADOG_API_KEY"),
		DDSite:           os.Getenv("DATADOG_SITE"),
		Source:           "env_vars",
	}

	// Fallback to headers
	if cfg.SelectedProvider == "" {
		cfg.SelectedProvider = getHeader(headers, "x-selected-provider")
	}
	if cfg.SelectedProvider == "" {
		cfg.SelectedProvider = "newrelic" // default
	}

	if cfg.LicenseKey == "" {
		cfg.LicenseKey = getHeader(headers, "x-nr-license-key")
		if cfg.LicenseKey != "" {
			cfg.Source = "request"
		}
	}
	if cfg.AccountID == "" {
		cfg.AccountID = getHeader(headers, "x-nr-account-id")
		if cfg.AccountID != "" {
			cfg.Source = "request"
		}
	}
	if cfg.APIKey == "" {
		cfg.APIKey = getHeader(headers, "x-nr-api-key")
	}
	if cfg.Region == "" {
		cfg.Region = getHeader(headers, "x-nr-region")
	}

	// Default region to US
	if cfg.Region == "" {
		cfg.Region = "us"
	}
	cfg.Region = strings.ToLower(cfg.Region)

	if cfg.DDApiKey == "" {
		cfg.DDApiKey = getHeader(headers, "x-dd-api-key")
		if cfg.DDApiKey != "" {
			cfg.Source = "request"
		}
	}
	if cfg.DDSite == "" {
		cfg.DDSite = getHeader(headers, "x-dd-site")
	}
	if cfg.DDSite == "" {
		cfg.DDSite = "datadoghq.com" // default site
	}
	cfg.DDSite = strings.ToLower(cfg.DDSite)

	return cfg, nil
}

// Validate checks that the minimum required credentials are present.
func (c *Config) Validate() error {
	if c.SelectedProvider == "datadog" {
		if c.DDApiKey == "" {
			return errors.New("DATADOG_API_KEY not configured (set via env var or x-dd-api-key header)")
		}
		return nil
	}
	if c.LicenseKey == "" {
		return errors.New("NEW_RELIC_LICENSE_KEY not configured (set via env var or x-nr-license-key header)")
	}
	if c.AccountID == "" {
		return errors.New("NEW_RELIC_ACCOUNT_ID not configured (set via env var or x-nr-account-id header)")
	}
	return nil
}

// IsFullyConfigured returns true if all required keys are set via env vars.
func (c *Config) IsFullyConfigured() bool {
	if c.SelectedProvider == "datadog" {
		return c.DDApiKey != ""
	}
	return c.LicenseKey != "" && c.AccountID != ""
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

// HasDDApiKey returns true if a Datadog API key is available.
func (c *Config) HasDDApiKey() bool {
	return c.DDApiKey != ""
}

// HasDDSite returns true if a Datadog site is available.
func (c *Config) HasDDSite() bool {
	return c.DDSite != ""
}
