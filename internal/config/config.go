package config

import (
	"net/http"
	"time"

	"github.com/eka-care/eka-sdk-go/internal/interfaces"
)

// Environment represents the deployment environment
type Environment string

const (
	EnvironmentProduction  Environment = "production"
	EnvironmentDevelopment Environment = "development"
)

// GetBaseURL returns the base URL for the environment
func (e Environment) GetBaseURL() string {
	switch e {
	case EnvironmentDevelopment:
		return "https://api.dev.eka.care"
	case EnvironmentProduction:
		return "https://api.eka.care"
	default:
		return "https://api.eka.care" // Default to production
	}
}

// Config holds the internal configuration for the ABDM client
type Config struct {
	Environment        Environment
	BaseURL            string
	APIKey             string
	AuthorizationToken string
	Timeout            time.Duration
	MaxRetries         int
	UserAgent          string
	LogLevel           string
	HTTPClient         *http.Client
	DisableSSL         bool
	Region             string
	RetryMode          string
	MaxBackoffDelay    time.Duration
	RequestTimeout     time.Duration
	ResponseTimeout    time.Duration
	ConnectionTimeout  time.Duration
}

// Ensure Config implements interfaces.Config
var _ interfaces.Config = (*Config)(nil)

// NewConfig creates a new configuration with defaults
func NewConfig() *Config {
	return &Config{
		Environment:       EnvironmentProduction,
		BaseURL:           "https://api.eka.care",
		Timeout:           30 * time.Second,
		MaxRetries:        3,
		UserAgent:         "eka-sdk-go/1.0",
		LogLevel:          "info",
		RetryMode:         "standard",
		MaxBackoffDelay:   20 * time.Second,
		RequestTimeout:    30 * time.Second,
		ResponseTimeout:   30 * time.Second,
		ConnectionTimeout: 10 * time.Second,
	}
}

// Interface implementation methods
func (c *Config) GetEnvironment() Environment { return c.Environment }
func (c *Config) GetBaseURL() string          { return c.BaseURL }
func (c *Config) GetAPIKey() string {
	switch {
	case c.AuthorizationToken != "":
		return c.AuthorizationToken
	default:
		return c.APIKey
	}
}
func (c *Config) GetTimeout() time.Duration           { return c.Timeout }
func (c *Config) GetMaxRetries() int                  { return c.MaxRetries }
func (c *Config) GetUserAgent() string                { return c.UserAgent }
func (c *Config) GetLogLevel() string                 { return c.LogLevel }
func (c *Config) GetHTTPClient() *http.Client         { return c.HTTPClient }
func (c *Config) GetDisableSSL() bool                 { return c.DisableSSL }
func (c *Config) GetRegion() string                   { return c.Region }
func (c *Config) GetRetryMode() string                { return c.RetryMode }
func (c *Config) GetMaxBackoffDelay() time.Duration   { return c.MaxBackoffDelay }
func (c *Config) GetRequestTimeout() time.Duration    { return c.RequestTimeout }
func (c *Config) GetResponseTimeout() time.Duration   { return c.ResponseTimeout }
func (c *Config) GetConnectionTimeout() time.Duration { return c.ConnectionTimeout }
