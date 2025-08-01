package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/eka-care/eka-sdk-go/abdm"
)

// demonstrateEnvironmentConfiguration shows how to use the SDK with environment variables
// This is the recommended approach for production applications
func demonstrateEnvironmentConfiguration() {
	fmt.Println("=== Environment-First Configuration (Recommended) ===")

	// Method 1: Set environment variables programmatically (for demo purposes)
	// In real applications, set these in your shell, Docker, Kubernetes, etc.
	os.Setenv("EKA_ENVIRONMENT", "development")
	os.Setenv("EKA_AUTH_TOKEN", "your-development-token")
	os.Setenv("EKA_TIMEOUT", "45")
	os.Setenv("EKA_MAX_RETRIES", "5")
	os.Setenv("EKA_USER_AGENT", "my-healthcare-app/2.1")
	os.Setenv("EKA_LOG_LEVEL", "debug")

	// Create client using environment variables
	client := abdm.NewFromEnv()

	// Display configuration
	config := client.GetConfig()
	fmt.Printf("Environment: %s\n", config.Environment)
	fmt.Printf("Base URL: %s\n", config.BaseURL)
	fmt.Printf("Timeout: %v\n", config.Timeout)
	fmt.Printf("Max Retries: %d\n", config.MaxRetries)
	fmt.Printf("User Agent: %s\n", config.UserAgent)
	fmt.Printf("Log Level: %s\n", config.LogLevel)

	// Test the client
	ctx := context.Background()
	fmt.Println("\nTesting API connectivity...")

	// This would fail without proper auth token, but demonstrates the pattern
	if err := client.Ping(ctx); err != nil {
		fmt.Printf("Expected ping failure (no real token): %v\n", err)
	}

	// Clean up environment for demo
	os.Unsetenv("EKA_ENVIRONMENT")
	os.Unsetenv("EKA_AUTH_TOKEN")
	os.Unsetenv("EKA_TIMEOUT")
	os.Unsetenv("EKA_MAX_RETRIES")
	os.Unsetenv("EKA_USER_AGENT")
	os.Unsetenv("EKA_LOG_LEVEL")
}

// demonstrateDeploymentPatterns shows common deployment patterns
func demonstrateDeploymentPatterns() {
	fmt.Println("\n=== Common Deployment Patterns ===")

	// Pattern 1: Docker environment
	fmt.Println("Docker environment variables:")
	fmt.Println("docker run -e EKA_ENVIRONMENT=production \\")
	fmt.Println("           -e EKA_AUTH_TOKEN=your-prod-token \\")
	fmt.Println("           -e EKA_TIMEOUT=30 \\")
	fmt.Println("           your-app:latest")

	// Pattern 2: Kubernetes ConfigMap/Secret
	fmt.Println("\nKubernetes deployment:")
	fmt.Println("apiVersion: apps/v1")
	fmt.Println("kind: Deployment")
	fmt.Println("spec:")
	fmt.Println("  template:")
	fmt.Println("    spec:")
	fmt.Println("      containers:")
	fmt.Println("      - name: app")
	fmt.Println("        env:")
	fmt.Println("        - name: EKA_ENVIRONMENT")
	fmt.Println("          value: \"production\"")
	fmt.Println("        - name: EKA_AUTH_TOKEN")
	fmt.Println("          valueFrom:")
	fmt.Println("            secretKeyRef:")
	fmt.Println("              name: eka-secrets")
	fmt.Println("              key: auth-token")

	// Pattern 3: Shell environment
	fmt.Println("\nShell environment:")
	fmt.Println("export EKA_ENVIRONMENT=production")
	fmt.Println("export EKA_AUTH_TOKEN=your-token")
	fmt.Println("go run main.go")
}

// demonstrateCompatibilityMode shows the available configuration approaches
func demonstrateCompatibilityMode() {
	fmt.Println("\n=== Available Configuration Approaches ===")

	// Environment-first approach (recommended)
	os.Setenv("EKA_ENVIRONMENT", "production")
	os.Setenv("EKA_AUTH_TOKEN", "your-token")
	envClient := abdm.NewFromEnv()
	fmt.Printf("Environment-based client: %s\n", envClient.GetConfig().BaseURL)

	// Explicit configuration approach
	explicitClient := abdm.New(
		abdm.WithEnvironment(abdm.EnvironmentDevelopment),
		abdm.WithAuthorizationToken("your-token"),
	)
	fmt.Printf("Explicit configuration client: %s\n", explicitClient.GetConfig().BaseURL)

	fmt.Println("\nRecommendation: Use NewFromEnv() for modern cloud-native applications.")

	// Clean up
	os.Unsetenv("EKA_ENVIRONMENT")
	os.Unsetenv("EKA_AUTH_TOKEN")
}

// demonstrateFlexibleConfiguration shows how to combine environment variables with options
func demonstrateFlexibleConfiguration() {
	fmt.Println("\n=== Flexible Configuration ===")

	// Set base environment
	os.Setenv("EKA_ENVIRONMENT", "production")
	os.Setenv("EKA_AUTH_TOKEN", "your-token")

	// Create client with environment + additional options
	client := abdm.NewFromEnv(
		// Override or add specific configurations
		abdm.WithUserAgent("my-special-app/1.0"),
		abdm.WithMaxRetries(10), // Override default
		abdm.WithLogLevel(abdm.LogLevelWarn),
	)

	config := client.GetConfig()
	fmt.Printf("Environment from env var: %s\n", config.Environment)
	fmt.Printf("User agent from option: %s\n", config.UserAgent)
	fmt.Printf("Max retries from option: %d\n", config.MaxRetries)

	// Clean up
	os.Unsetenv("EKA_ENVIRONMENT")
	os.Unsetenv("EKA_AUTH_TOKEN")
}

func runEnvironmentExample() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	demonstrateEnvironmentConfiguration()
	demonstrateDeploymentPatterns()
	demonstrateCompatibilityMode()
	demonstrateFlexibleConfiguration()

	fmt.Println("\n=== Summary ===")
	fmt.Println("✅ Use NewFromEnv() for new applications")
	fmt.Println("✅ Set environment variables in deployment")
	fmt.Println("✅ Follow 12-factor app principles")
	fmt.Println("✅ Use New() with options for explicit configuration")
	fmt.Println("\nFor more information, visit: https://developer.eka.care")
}
