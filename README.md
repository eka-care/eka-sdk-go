# ABDM SDK for Go

A comprehensive and extensible Go SDK for the ABDM (Ayushman Bharat Digital Mission) API, following best practices from AWS and New Relic SDKs.

## Installation

```bash
go get github.com/eka-care/eka-sdk-go
```

## Quick Start

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/eka-care/eka-sdk-go/abdm"
)

func main() {
    // Create SDK with configuration
    sdk := abdm.New(
        abdm.WithBaseURL("https://api.eka.care"),
        abdm.WithAPIKey("your-api-key"),
        abdm.WithTimeout(30*time.Second),
        abdm.WithMaxRetries(3),
        abdm.WithLogLevel(abdm.LogLevelInfo),
    )

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

## Configuration

The SDK uses a functional options pattern for configuration, making it easy to customize behavior:

```go
sdk := abdm.New(
    abdm.WithBaseURL("https://api.eka.care"),
    abdm.WithAPIKey("your-api-key"),
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
    // Create SDK
    sdk := abdm.New(
        abdm.WithBaseURL("https://api.eka.care"),
        abdm.WithAPIKey("your-api-key"),
        abdm.WithTimeout(30*time.Second),
        abdm.WithMaxRetries(3),
    )

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

For support and questions, please open an issue on GitHub or contact the development team.