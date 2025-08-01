# Environment Configuration Guide

The Eka SDK Go supports multiple deployment environments with automatic URL resolution and flexible configuration options.

## Supported Environments

| Environment | Base URL | Use Case |
|-------------|----------|----------|
| `EnvironmentProduction` | `https://api.eka.care` | Production deployment |
| `EnvironmentDevelopment` | `https://api.dev.eka.care` | Development/testing |

## Usage Examples

### 1. Environment Variables (Recommended)

Set environment variables for configuration:

```bash
# Required
export EKA_AUTH_TOKEN="your-auth-token"

# Environment (defaults to production)
export EKA_ENVIRONMENT="production"  # or "development"

# Optional configuration
export EKA_TIMEOUT="30"              # timeout in seconds
export EKA_MAX_RETRIES="3"           # number of retries
export EKA_USER_AGENT="my-app/1.0"   # custom user agent
export EKA_LOG_LEVEL="info"          # debug, info, warn, error
```

```go
// Create client from environment variables
client := abdm.NewFromEnv()

// Or with additional options
client := abdm.NewFromEnv(
    abdm.WithTimeout(60*time.Second),
)
```

### 2. Explicit Environment Setting

```go
client := abdm.New(
    abdm.WithEnvironment(abdm.EnvironmentDevelopment),
    abdm.WithAuthorizationToken("your-auth-token"),
    abdm.WithTimeout(30*time.Second),
)
```

## Environment Variables

The SDK supports configuration via environment variables using `NewFromEnv()`:

| Variable | Description | Default |
|----------|-------------|---------|
| `EKA_ENVIRONMENT` | Environment: "production", "development" | "production" |
| `EKA_AUTH_TOKEN` | Authorization token | Required |
| `EKA_API_KEY` | API key (alternative to auth token) | - |
| `EKA_TIMEOUT` | Request timeout in seconds | 30 |
| `EKA_MAX_RETRIES` | Maximum number of retries | 3 |
| `EKA_USER_AGENT` | Custom user agent | "eka-sdk-go/1.0" |
| `EKA_LOG_LEVEL` | Log level: debug, info, warn, error | "info" |

### Example Usage:
```bash
export EKA_ENVIRONMENT=development
export EKA_AUTH_TOKEN=your-auth-token
export EKA_TIMEOUT=60
export EKA_MAX_RETRIES=5
export EKA_LOG_LEVEL=debug
```

```go
// Automatically reads from environment variables
client := abdm.NewFromEnv()

// Can still override with additional options
client := abdm.NewFromEnv(
    abdm.WithUserAgent("my-app/1.0"),
)
```

## Migration from Hardcoded URLs

### Before:
```go
client := abdm.New(
    abdm.WithBaseURL("https://api.dev.eka.care"),
    abdm.WithAuthorizationToken("your-auth-token"),
)
```

### After:
```go
// Recommended approach with environment variables
export EKA_ENVIRONMENT=development
export EKA_AUTH_TOKEN=your-auth-token
client := abdm.NewFromEnv()

// Or explicit environment configuration
client := abdm.New(
    abdm.WithEnvironment(abdm.EnvironmentDevelopment),
    abdm.WithAuthorizationToken("your-auth-token"),
)
```

## Best Practices

1. **Use environment variables** for configuration (recommended)
2. **Set environment via configuration** rather than hardcoding URLs
3. **Configure timeouts and retries** per environment needs
4. **Follow 12-factor app principles** for cloud-native deployment

## Configuration Priority

1. Environment variables (recommended)
2. Explicit options via `WithEnvironment()`
3. Default production environment (fallback)
