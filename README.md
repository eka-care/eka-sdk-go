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

### Registration Service

The registration service provides methods for ABHA registration via Aadhaar and mobile:

```go
// Aadhaar Registration
initResp, err := sdk.Registration().AadhaarInit(ctx, headers, abdm.InitRequest{
    AadhaarNumber: "123456789012",
})

verifyResp, err := sdk.Registration().AadhaarVerify(ctx, headers, abdm.VerifyRequest{
    TxnID:  initResp.TxnID,
    OTP:    "123456",
    Mobile: "9876543210",
})

// Mobile Registration
initResp, err := sdk.Registration().MobileInit(ctx, headers, abdm.MobileInitRequest{
    MobileNumber: "9876543210",
})

verifyResp, err := sdk.Registration().MobileVerify(ctx, headers, abdm.MobileVerifyOTPRequest{
    TxnID: initResp.TxnID,
    OTP:   "123456",
})
```

### Utils Service

The utils service provides utility functions for validation and common operations:

```go
// Validate inputs
err := sdk.Utils().ValidateAadhaarNumber("123456789012")
err := sdk.Utils().ValidateMobileNumber("9876543210")
err := sdk.Utils().ValidateABHAAddress("john.doe@abdm")

// Generate transaction IDs
txnID := sdk.Utils().GenerateTransactionID()

// Format dates
dateStr := sdk.Utils().FormatDate(1990, 5, 15)
```

## Middleware

The SDK supports middleware for logging, metrics, retries, and authentication:

```go
// Add logging middleware
logger := &customLogger{}
sdk.AddMiddleware(middleware.LoggingMiddleware(logger))

// Add metrics middleware
metrics := &customMetrics{}
sdk.AddMiddleware(middleware.MetricsMiddleware(metrics))

// Add retry middleware
sdk.AddMiddleware(middleware.RetryMiddleware(3, 1*time.Second))

// Add authentication middleware
sdk.AddMiddleware(middleware.AuthMiddleware(func(req *http.Request) error {
    // Custom authentication logic
    return nil
}))
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

Here's a complete example showing the Aadhaar registration flow:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

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

    // Step 1: Initialize Aadhaar registration
    initResp, err := sdk.Registration().AadhaarInit(ctx, headers, abdm.InitRequest{
        AadhaarNumber: "123456789012",
    })
    if err != nil {
        log.Fatalf("Failed to initialize: %v", err)
    }

    fmt.Printf("Transaction ID: %s\n", initResp.TxnID)

    // Step 2: Verify Aadhaar OTP
    verifyResp, err := sdk.Registration().AadhaarVerify(ctx, headers, abdm.VerifyRequest{
        TxnID:  initResp.TxnID,
        OTP:    "123456",
        Mobile: "9876543210",
    })
    if err != nil {
        log.Fatalf("Failed to verify: %v", err)
    }

    // Step 3: Create ABHA address (if needed)
    if verifyResp.SkipState == "abha_create" {
        createResp, err := sdk.Registration().AadhaarCreatePHR(ctx, headers, abdm.CreateRequest{
            TxnID:       verifyResp.TxnID,
            AbhaAddress: "john.doe@abdm",
        })
        if err != nil {
            log.Fatalf("Failed to create ABHA: %v", err)
        }

        if createResp.Profile != nil {
            fmt.Printf("ABHA created: %+v\n", createResp.Profile)
        }
    }
}
```

```
eka-sdk-go/
├── abdm/                    # Main SDK package
│   ├── client.go           # Main client implementation
│   └── types.go            # Request/response types
├── internal/               # Internal packages
│   ├── config/            # Configuration management
│   ├── http/              # HTTP client implementation
│   ├── middleware/        # Middleware implementations
│   ├── registration/      # Registration service
│   ├── utils/             # Utility functions
│   └── errors/            # Error handling
└── example/               # Usage examples
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