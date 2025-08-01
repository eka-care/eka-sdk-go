# Eka Care SDK for Go

A Go SDK for integrating with [Eka Care's healthcare APIs](https://developer.eka.care), including ABDM (Ayushman Bharat Digital Mission) services.

## Getting Started

### 1. Get Your Credentials

1. Register at [Eka Care Developer Portal](https://developer.eka.care)
2. Get your `client_id` and `client_secret` from the portal

### 2. Install the SDK

```bash
go get github.com/eka-care/eka-sdk-go
```

### 3. Quick Start Example

```go
package main

import (
    "context"
    "log"
    
    ekasdk "github.com/eka-care/eka-sdk-go"
    "github.com/eka-care/eka-sdk-go/services/abdm/abha/login"
    "github.com/eka-care/eka-sdk-go/internal/interfaces"
)

func main() {
    ctx := context.Background()

    // Create SDK client from environment variables
    client := ekasdk.NewFromEnv()

    // Authenticate with Eka platform
    if err := client.Login(ctx); err != nil {
        log.Fatalf("Authentication failed: %v", err)
    }
    
    // Now you can use any Eka Care API
    headers := interfaces.Headers{
        UserID: "your-user-id",
        HipID:  "your-hip-id",
    }
    
    // Example: ABDM login
    otpReq := &login.InitLoginRequest{
        Identifier: "demo@abdm",
        Method:     login.LoginMethodPhrAddress,
    }
    
    otpResp, err := client.ABDM.Login().LoginInit(ctx, headers, otpReq)
    if err != nil {
        log.Printf("ABDM login failed: %v", err)
        return
    }
    
    log.Printf("Success! Transaction ID: %s", otpResp.TxnID)
}
```

## Environment Setup

Set these environment variables:

```bash
export EKA_CLIENT_ID="your-client-id"
export EKA_CLIENT_SECRET="your-client-secret"
export EKA_ENVIRONMENT="production"  # or "development"
```

That's it! The SDK handles authentication automatically.

## Available Services

Once authenticated, you can access:

- **ABDM Services**: `client.ABDM.Login()`, `client.ABDM.Registration()`, `client.ABDM.Profile()`
- **More services** will be added as they become available

## Need Help?

- **Documentation**: [developer.eka.care](https://developer.eka.care)
- **API Reference**: [API Docs](https://developer.eka.care/api-reference)
- **Issues**: [GitHub Issues](https://github.com/eka-care/eka-sdk-go/issues)