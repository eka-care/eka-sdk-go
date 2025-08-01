# Authentication

The SDK handles authentication automatically using your client credentials.

## Quick Start

1. Set your credentials:
```bash
export EKA_CLIENT_ID="your-client-id"
export EKA_CLIENT_SECRET="your-client-secret"
```

2. Use the SDK:
```go
client := ekasdk.NewFromEnv()
err := client.Login(ctx)  // Automatically authenticates
```

That's it! The SDK manages tokens, refreshing, and authentication for you.

## How It Works

- The SDK automatically logs in using your client credentials
- It refreshes tokens before they expire
- All API calls use the current valid token
- You don't need to manage tokens manually

This works for all Eka Care APIs (ABDM and future services).
