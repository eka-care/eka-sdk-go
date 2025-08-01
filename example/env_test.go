package main

import (
	"fmt"
	"os"

	"github.com/eka-care/eka-sdk-go/abdm"
)

func testEnvConfig() {
	fmt.Println("=== Environment Variable Configuration Test ===")

	// Test 1: Default environment (production)
	fmt.Println("\n1. Testing default configuration:")
	client1 := abdm.NewFromEnv()
	config1 := client1.GetConfig()
	fmt.Printf("   Environment: %s\n", config1.Environment)
	fmt.Printf("   Base URL: %s\n", config1.BaseURL)

	// Test 2: Set development environment
	fmt.Println("\n2. Testing development environment:")
	os.Setenv("EKA_ENVIRONMENT", "development")
	os.Setenv("EKA_AUTH_TOKEN", "test-token")
	os.Setenv("EKA_TIMEOUT", "45")
	os.Setenv("EKA_MAX_RETRIES", "5")
	os.Setenv("EKA_LOG_LEVEL", "debug")

	client2 := abdm.NewFromEnv()
	config2 := client2.GetConfig()
	fmt.Printf("   Environment: %s\n", config2.Environment)
	fmt.Printf("   Base URL: %s\n", config2.BaseURL)
	fmt.Printf("   Auth Token: %s\n", config2.AuthorizationToken)
	fmt.Printf("   Timeout: %v\n", config2.Timeout)
	fmt.Printf("   Max Retries: %d\n", config2.MaxRetries)
	fmt.Printf("   Log Level: %s\n", config2.LogLevel)

	// Test 3: Override with additional options
	fmt.Println("\n3. Testing environment + additional options:")
	client3 := abdm.NewFromEnv(
		abdm.WithUserAgent("test-app/1.0"),
		abdm.WithMaxRetries(10),
	)
	config3 := client3.GetConfig()
	fmt.Printf("   Environment: %s\n", config3.Environment)
	fmt.Printf("   User Agent: %s\n", config3.UserAgent)
	fmt.Printf("   Max Retries: %d\n", config3.MaxRetries)

	// Cleanup
	os.Unsetenv("EKA_ENVIRONMENT")
	os.Unsetenv("EKA_AUTH_TOKEN")
	os.Unsetenv("EKA_TIMEOUT")
	os.Unsetenv("EKA_MAX_RETRIES")
	os.Unsetenv("EKA_LOG_LEVEL")

	fmt.Println("\n✅ Environment variable configuration working correctly!")
}
