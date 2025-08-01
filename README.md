# ABDM SDK for Go

A comprehensive and extensible Go SDK for the ABDM (Ayushman Bharat Digital Mission) API. This SDK is available through the [Eka Care Developer Portal](https://developer.eka.care) for integrating with Eka's healthcare APIs.

## Getting Started

### Prerequisites

Before using this SDK, you need to:

1. **Register at [Eka Care Developer Portal](https://developer.eka.care)**
2. **Get your API credentials** (client_id and client_secret)
3. **Obtain an access token** using the [authentication API](https://developer.eka.care/api-reference/authorization/getting-started)

### Installation

```bash
go get github.com/eka-care/eka-sdk-go
```

## Quick Start

### Authentication

First, authenticate your client to get an access token. Visit the [Eka Care Developer Portal](https://developer.eka.care/api-reference/authorization/getting-started) for detailed authentication steps.

```go
// Example: After obtaining your access token from the authentication API
accessToken := "your-access-token-here"
```

### Basic Usage

```go
package main

import (
    "context"
    "log"

    "github.com/eka-care/eka-sdk-go/abdm"
)

func main() {
    // Recommended: Create SDK using environment variables
    sdk := abdm.NewFromEnv()
    
    // Alternative: Explicit configuration
    // sdk := abdm.New(
    //     abdm.WithEnvironment(abdm.EnvironmentProduction), // or EnvironmentDevelopment
    //     abdm.WithAuthorizationToken("your-access-token"),
    //     abdm.WithTimeout(30*time.Second),
    //     abdm.WithMaxRetries(3),
    // )

    ctx := context.Background()

    // Test connectivity
    if err := sdk.Ping(ctx); err != nil {
        log.Fatalf("Failed to ping API: %v", err)
    }

    // Use registration service
    headers := abdm.Headers{
        UserID: "user123",
        HipID:  "hip456",
    }

    // Initialize Aadhaar registration
    initResp, err := sdk.Registration().AadhaarInit(ctx, headers, abdm.InitRequest{
        AadhaarNumber: "123456789012",
    })
    if err != nil {
        log.Fatalf("Failed to initialize: %v", err)
    }

    fmt.Printf("Transaction ID: %s\n", initResp.TxnID)
}
```

## Environment Configuration

The SDK automatically detects the environment using environment variables (recommended approach):

### Environment Variables

Set these environment variables to configure the SDK:

```bash
# Required
export EKA_AUTH_TOKEN="your-access-token"

# Environment (defaults to production)
export EKA_ENVIRONMENT="production"  # or "development"

# Optional configuration
export EKA_TIMEOUT="30"              # timeout in seconds
export EKA_MAX_RETRIES="3"           # number of retries
export EKA_USER_AGENT="my-app/1.0"   # custom user agent
export EKA_LOG_LEVEL="info"          # debug, info, warn, error
```

### Usage with Environment Variables

```go
// SDK automatically reads configuration from environment variables
sdk := abdm.NewFromEnv()

// You can still override specific settings
sdk := abdm.NewFromEnv(
    abdm.WithTimeout(60*time.Second),  // Override default timeout
    abdm.WithMaxRetries(5),            // Override default retries
)
```

### Manual Configuration (Alternative)

If you prefer explicit configuration instead of environment variables:

```go
// Production
sdk := abdm.New(
    abdm.WithEnvironment(abdm.EnvironmentProduction),
    abdm.WithAuthorizationToken("your-prod-access-token"),
)

// Development  
sdk := abdm.New(
    abdm.WithEnvironment(abdm.EnvironmentDevelopment),
    abdm.WithAuthorizationToken("your-dev-access-token"),
)
```

| Environment | Base URL |
|-------------|----------|
| Production | `https://api.eka.care` |
| Development | `https://api.dev.eka.care` |

> **Note**: Access tokens are environment-specific. Make sure to use the appropriate token for each environment. Get your tokens from the [Eka Care Developer Portal](https://developer.eka.care).

For more details, see [Environment Configuration Guide](./ENVIRONMENT_CONFIG.md).

## Authentication & Security

This SDK integrates with Eka Care's secure authentication system:

- **Access Tokens**: Include your access token in all API requests
- **Token Management**: Implement token refresh logic for expired tokens (401 errors)
- **Environment Separation**: Use different tokens for production and development
- **Security**: Never expose your client credentials in client-side code

For detailed authentication steps, visit the [Getting Started Guide](https://developer.eka.care/api-reference/authorization/getting-started).

## Configuration

The SDK uses a functional options pattern for configuration:

```go
sdk := abdm.New(
    abdm.WithEnvironment(abdm.EnvironmentDevelopment),
    abdm.WithAuthorizationToken("your-access-token"), // Use token from Eka Care auth API
    abdm.WithTimeout(30*time.Second),
    abdm.WithMaxRetries(3),
    abdm.WithUserAgent("my-app/1.0"),
    abdm.WithLogLevel(abdm.LogLevelDebug),
    abdm.WithRetryMode(abdm.RetryModeStandard),
    abdm.WithMaxBackoffDelay(20*time.Second),
    abdm.WithRequestTimeout(30*time.Second),
    abdm.WithResponseTimeout(30*time.Second),
    abdm.WithConnectionTimeout(10*time.Second),
    abdm.WithRegion("us-east-1"),
    abdm.WithCredentials(&abdm.Credentials{
        APIKey: "your-api-key",
    }),
)
```

## Services

The SDK provides the following services for ABDM integration:

### Registration Service
- **Aadhaar Registration**: Initialize and verify ABHA registration using Aadhaar
- **Mobile Registration**: Create ABHA account using mobile number
- **ABHA Creation**: Generate ABHA address and complete registration

### Login Service  
- **Multi-method Login**: Support for Aadhaar, mobile, ABHA number, and ABHA address
- **OTP Verification**: Secure authentication with OTP validation
- **Session Management**: Handle login sessions and tokens

### Profile Service
- **Profile Management**: Retrieve and manage user profile information
- **Digital Assets**: Generate ABHA cards and QR codes
- **Session APIs**: Initialize and verify secure sessions

### Utils Service
- **Validation**: Validate Aadhaar numbers, mobile numbers, and ABHA addresses
- **Utilities**: Generate transaction IDs, format dates, retry mechanisms

## Basic Usage Examples

```go
// Initialize SDK
sdk := abdm.NewFromEnv()

// Registration example
initResp, err := sdk.Registration().AadhaarInit(ctx, headers, request)

// Login example  
loginResp, err := sdk.Login().LoginInit(ctx, headers, request)

// Profile example
profile, err := sdk.Profile().GetProfile(ctx, headers)

// Utils example
err := sdk.Utils().ValidateAadhaarNumber("123456789012")
```

## Error Handling

The SDK provides comprehensive error handling:

```go
resp, err := sdk.Registration().AadhaarInit(ctx, headers, req)
if err != nil {
    if apiErr, ok := err.(*abdm.APIError); ok {
        log.Printf("API Error: %d - %s", apiErr.Code, apiErr.Message)
    } else {
        log.Printf("Network Error: %v", err)
    }
    return
}
```

## Complete Example

```go
package main

import (
    "context"
    "log"

    "github.com/eka-care/eka-sdk-go/abdm"
)

func main() {
    // Create SDK using environment variables (recommended)
    sdk := abdm.NewFromEnv()

    ctx := context.Background()
    headers := abdm.Headers{
        UserID: "user123",
        HipID:  "hip456",
    }

    // Test connectivity
    if err := sdk.Ping(ctx); err != nil {
        log.Fatalf("Failed to ping API: %v", err)
    }

    // Example: Initialize registration
    initResp, err := sdk.Registration().AadhaarInit(ctx, headers, request)
    if err != nil {
        log.Fatalf("Registration failed: %v", err)
    }

    log.Printf("Registration initiated: %s", initResp.TxnID)
}
```

For detailed examples, see the `example/` directory:
- `api_demo.go` - Complete API workflows
- `environment_example.go` - Configuration examples  
- `profile_example.go` - Profile and session APIs

## Error Handling

```go
resp, err := sdk.Registration().AadhaarInit(ctx, headers, req)
if err != nil {
    if apiErr, ok := err.(*abdm.APIError); ok {
        log.Printf("API Error: %d - %s", apiErr.Code, apiErr.Message)
    } else {
        log.Printf("Network Error: %v", err)
    }
    return
}
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support and questions:

- **Documentation**: Visit the [Eka Care Developer Portal](https://developer.eka.care)
- **API Reference**: [https://developer.eka.care/api-reference](https://developer.eka.care/api-reference)
- **Issues**: Open an issue on GitHub
- **Developer Support**: Contact the development team through the developer portal

## Resources

- [Eka Care Developer Portal](https://developer.eka.care)
- [Authentication Guide](https://developer.eka.care/api-reference/authorization/getting-started)
- [API Documentation](https://developer.eka.care/api-reference)
- [Environment Configuration Guide](./ENVIRONMENT_CONFIG.md)