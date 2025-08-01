# Environment Configuration

The Eka SDK supports two environments with automatic configuration.

## Quick Setup

Set these environment variables:

```bash
export EKA_CLIENT_ID="your-client-id"
export EKA_CLIENT_SECRET="your-client-secret"
export EKA_ENVIRONMENT="production"  # or "development"
```

Then use the SDK:

```go
client := ekasdk.NewFromEnv()
err := client.Login(ctx)
```

## Environments

| Environment | Base URL | Use For |
|-------------|----------|---------|
| `production` | `https://api.eka.care` | Live applications |
| `development` | `https://api.dev.eka.care` | Testing |

## Required Variables

- `EKA_CLIENT_ID`: Your client ID from the developer portal
- `EKA_CLIENT_SECRET`: Your client secret from the developer portal
- `EKA_ENVIRONMENT`: Either "production" or "development"

That's it! The SDK handles everything else automatically.

## Configuration Priority

1. Environment variables (recommended)
2. Explicit options via `WithEnvironment()`
3. Default production environment (fallback)
