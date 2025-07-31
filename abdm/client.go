package abdm

import (
	"context"
	"github.com/eka-care/eka-sdk-go/internal/abha/login"
	"github.com/eka-care/eka-sdk-go/internal/abha/registration"
	"net/http"
	"time"

	"github.com/eka-care/eka-sdk-go/internal/config"
	"github.com/eka-care/eka-sdk-go/internal/interfaces"
	"github.com/eka-care/eka-sdk-go/internal/utils"
)

// Client represents the ABDM SDK client
type Client struct {
	config       interfaces.Config
	httpClient   *http.Client
	registration *registration.Service
	login        *login.Service
	utils        *utils.Service
	middleware   []interfaces.Middleware
}

// Ensure Client implements interfaces.Config
var _ interfaces.Config = (*Client)(nil)

// Config holds the configuration for the ABDM client
type Config struct {
	BaseURL            string        `json:"base_url"`
	APIKey             string        `json:"api_key"`
	AuthorizationToken string        `json:"authorization_token"`
	Timeout            time.Duration `json:"timeout"`
	MaxRetries         int           `json:"max_retries"`
	UserAgent          string        `json:"user_agent"`
	LogLevel           LogLevel      `json:"log_level"`
	HTTPClient         *http.Client  `json:"-"`
	DisableSSL         bool          `json:"disable_ssl"`
	Region             string        `json:"region"`
	Credentials        *Credentials  `json:"credentials,omitempty"`
	RetryMode          RetryMode     `json:"retry_mode"`
	MaxBackoffDelay    time.Duration `json:"max_backoff_delay"`
	RequestTimeout     time.Duration `json:"request_timeout"`
	ResponseTimeout    time.Duration `json:"response_timeout"`
	ConnectionTimeout  time.Duration `json:"connection_timeout"`
}

// Credentials represents authentication credentials
type Credentials struct {
	APIKey    string `json:"api_key"`
	SecretKey string `json:"secret_key,omitempty"`
	Token     string `json:"token,omitempty"`
}

// LogLevel represents the logging level
type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

// RetryMode represents the retry strategy
type RetryMode string

const (
	RetryModeStandard RetryMode = "standard"
	RetryModeAdaptive RetryMode = "adaptive"
)

// Option represents a configuration option
type Option func(*Config)

// WithBaseURL sets the base URL
func WithBaseURL(url string) Option {
	return func(c *Config) {
		c.BaseURL = url
	}
}

// WithAPIKey sets the API key
func WithAPIKey(key string) Option {
	return func(c *Config) {
		c.APIKey = key
	}
}

// WithAPIKey sets the Authorization token
func WithAuthorizationToken(token string) Option {
	return func(c *Config) {
		c.AuthorizationToken = token
	}
}

// WithTimeout sets the timeout
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

// WithMaxRetries sets the maximum number of retries
func WithMaxRetries(retries int) Option {
	return func(c *Config) {
		c.MaxRetries = retries
	}
}

// WithUserAgent sets the user agent
func WithUserAgent(userAgent string) Option {
	return func(c *Config) {
		c.UserAgent = userAgent
	}
}

// WithLogLevel sets the log level
func WithLogLevel(level LogLevel) Option {
	return func(c *Config) {
		c.LogLevel = level
	}
}

// WithHTTPClient sets the HTTP client
func WithHTTPClient(client *http.Client) Option {
	return func(c *Config) {
		c.HTTPClient = client
	}
}

// WithRegion sets the region
func WithRegion(region string) Option {
	return func(c *Config) {
		c.Region = region
	}
}

// WithCredentials sets the credentials
func WithCredentials(creds *Credentials) Option {
	return func(c *Config) {
		c.Credentials = creds
	}
}

// WithRetryMode sets the retry mode
func WithRetryMode(mode RetryMode) Option {
	return func(c *Config) {
		c.RetryMode = mode
	}
}

// WithMaxBackoffDelay sets the maximum backoff delay
func WithMaxBackoffDelay(delay time.Duration) Option {
	return func(c *Config) {
		c.MaxBackoffDelay = delay
	}
}

// WithRequestTimeout sets the request timeout
func WithRequestTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.RequestTimeout = timeout
	}
}

// WithResponseTimeout sets the response timeout
func WithResponseTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.ResponseTimeout = timeout
	}
}

// WithConnectionTimeout sets the connection timeout
func WithConnectionTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.ConnectionTimeout = timeout
	}
}

// New creates a new ABDM client with the given options
func New(opts ...Option) *Client {
	cfg := &Config{}

	// Apply options
	for _, opt := range opts {
		opt(cfg)
	}

	// Set defaults
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.eka.care"
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}
	if cfg.MaxRetries == 0 {
		cfg.MaxRetries = 3
	}
	if cfg.UserAgent == "" {
		cfg.UserAgent = "eka-sdk-go/1.0"
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = LogLevelInfo
	}
	if cfg.RetryMode == "" {
		cfg.RetryMode = RetryModeStandard
	}
	if cfg.MaxBackoffDelay == 0 {
		cfg.MaxBackoffDelay = 20 * time.Second
	}
	if cfg.RequestTimeout == 0 {
		cfg.RequestTimeout = 30 * time.Second
	}
	if cfg.ResponseTimeout == 0 {
		cfg.ResponseTimeout = 30 * time.Second
	}
	if cfg.ConnectionTimeout == 0 {
		cfg.ConnectionTimeout = 10 * time.Second
	}

	// Create HTTP client
	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: cfg.Timeout,
		}
	}

	// Create internal config
	internalConfig := &config.Config{
		BaseURL:            cfg.BaseURL,
		APIKey:             cfg.APIKey,
		AuthorizationToken: cfg.AuthorizationToken,
		Timeout:            cfg.Timeout,
		MaxRetries:         cfg.MaxRetries,
		UserAgent:          cfg.UserAgent,
		LogLevel:           string(cfg.LogLevel),
		HTTPClient:         httpClient,
		DisableSSL:         cfg.DisableSSL,
		Region:             cfg.Region,
		RetryMode:          string(cfg.RetryMode),
		MaxBackoffDelay:    cfg.MaxBackoffDelay,
		RequestTimeout:     cfg.RequestTimeout,
		ResponseTimeout:    cfg.ResponseTimeout,
		ConnectionTimeout:  cfg.ConnectionTimeout,
	}

	client := &Client{
		config:     internalConfig,
		httpClient: httpClient,
	}

	// Initialize services
	client.registration = registration.NewService(client)
	client.utils = utils.NewService(client)
	client.login = login.NewService(client)

	return client
}

// AddMiddleware adds middleware to the client
func (c *Client) AddMiddleware(middleware interfaces.Middleware) {
	c.middleware = append(c.middleware, middleware)
}

// GetHTTPClient returns the HTTP client
func (c *Client) GetHTTPClient() *http.Client {
	return c.httpClient
}

// GetMiddleware returns the middleware chain
func (c *Client) GetMiddleware() []interfaces.Middleware {
	return c.middleware
}

// SetHTTPClient allows customizing the HTTP client
func (c *Client) SetHTTPClient(client *http.Client) {
	c.httpClient = client
	c.config.(*config.Config).HTTPClient = client
}

// GetConfig returns the current configuration
func (c *Client) GetConfig() *Config {
	internalConfig := c.config.(*config.Config)
	return &Config{
		BaseURL:            internalConfig.BaseURL,
		APIKey:             internalConfig.APIKey,
		AuthorizationToken: internalConfig.AuthorizationToken,
		Timeout:            internalConfig.Timeout,
		MaxRetries:         internalConfig.MaxRetries,
		UserAgent:          internalConfig.UserAgent,
		LogLevel:           LogLevel(internalConfig.LogLevel),
		HTTPClient:         c.httpClient,
		DisableSSL:         internalConfig.DisableSSL,
		Region:             internalConfig.Region,
		Credentials:        nil,
		RetryMode:          RetryMode(internalConfig.RetryMode),
		MaxBackoffDelay:    internalConfig.MaxBackoffDelay,
		RequestTimeout:     internalConfig.RequestTimeout,
		ResponseTimeout:    internalConfig.ResponseTimeout,
		ConnectionTimeout:  internalConfig.ConnectionTimeout,
	}
}

// Registration returns the registration service
func (c *Client) Registration() *registration.Service {
	return c.registration
}

// Utils returns the utilities service
func (c *Client) Utils() *utils.Service {
	return c.utils
}

// Login returns the login service
func (c *Client) Login() *login.Service {
	return c.login
}

// Ping checks if the API is reachable
func (c *Client) Ping(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", c.config.GetBaseURL()+"/health", nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return &APIError{
			Code:    resp.StatusCode,
			Message: "API is not reachable",
		}
	}

	return nil
}

// APIError represents an API error
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return e.Message
}

// Interface implementation methods for Config
func (c *Client) GetBaseURL() string                  { return c.config.GetBaseURL() }
func (c *Client) GetAPIKey() string                   { return c.config.GetAPIKey() }
func (c *Client) GetTimeout() time.Duration           { return c.config.GetTimeout() }
func (c *Client) GetMaxRetries() int                  { return c.config.GetMaxRetries() }
func (c *Client) GetUserAgent() string                { return c.config.GetUserAgent() }
func (c *Client) GetLogLevel() string                 { return c.config.GetLogLevel() }
func (c *Client) GetDisableSSL() bool                 { return c.config.GetDisableSSL() }
func (c *Client) GetRegion() string                   { return c.config.GetRegion() }
func (c *Client) GetRetryMode() string                { return c.config.GetRetryMode() }
func (c *Client) GetMaxBackoffDelay() time.Duration   { return c.config.GetMaxBackoffDelay() }
func (c *Client) GetRequestTimeout() time.Duration    { return c.config.GetRequestTimeout() }
func (c *Client) GetResponseTimeout() time.Duration   { return c.config.GetResponseTimeout() }
func (c *Client) GetConnectionTimeout() time.Duration { return c.config.GetConnectionTimeout() }
