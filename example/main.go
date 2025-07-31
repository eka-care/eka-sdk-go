package main

import (
	"context"
	"fmt"
	"github.com/eka-care/eka-sdk-go/internal/abha/login"
	"github.com/eka-care/eka-sdk-go/internal/abha/registration"
	"log"
	"time"

	"github.com/eka-care/eka-sdk-go/abdm"
	"github.com/eka-care/eka-sdk-go/internal/interfaces"
)

func main() {
	// Create a new ABDM client with configuration
	client := abdm.New(
		abdm.WithBaseURL("https://api.dev.eka.care"),
		abdm.WithAuthorizationToken("your-auth-token"),
		abdm.WithTimeout(30*time.Second),
		abdm.WithMaxRetries(3),
		abdm.WithUserAgent("my-app/1.0"),
		abdm.WithLogLevel(abdm.LogLevelInfo),
	)

	// Add custom middleware (optional)
	// logger := &CustomLogger{}
	// metrics := &CustomMetrics{}
	// client.AddMiddleware(middleware.LoggingMiddleware(logger))
	// client.AddMiddleware(middleware.MetricsMiddleware(metrics))

	// Test the connection
	ctx := context.Background()
	if err := client.Ping(ctx); err != nil {
		log.Printf("Failed to ping API: %v", err)
	} else {
		log.Println("API is reachable")
	}

	//// Example: Aadhaar Registration Flow
	//exampleAadhaarRegistration(client)

	// Example: Mobile Registration Flow
	//exampleMobileRegistration(client)
	exampleMobileLogin(client)

	//// Example: Utility Functions
	//exampleUtilityFunctions(client)
}

func exampleMobileLogin(client *abdm.Client) {
	ctx := context.Background()
	headers := interfaces.Headers{}

	// Step 1: Initialize mobile registration
	//initResp, err := client.Login().LoginInit(ctx, headers, &login.InitLoginRequest{
	//	Identifier: "8590680079",
	//	Method:     login.LoginMethodMobile,
	//})
	//if err != nil {
	//	log.Printf("Failed to initialize mobile login: %v", err)
	//	return
	//}
	//
	//fmt.Printf("Mobile login initiated. Transaction ID: %v\n", initResp)

	//Step 2: Verify mobile OTP
	//verifyResp, err := client.Login().LoginVerify(ctx, headers, &login.VerifyLoginOTPRequest{
	//	OTP:   "906153",
	//	TxnID: "1e16b852-d54c-4ed3-bb13-c905df2cacfe",
	//})
	//if err != nil {
	//	log.Printf("Failed to verify mobile OTP: %v", err)
	//	return
	//}
	//
	//fmt.Printf("Mobile verification completed. Response: %v\n", verifyResp)

	// Step 3: Login with ABHA address
	loginResp, err := client.Login().LoginWithPHRAddress(ctx, headers, &login.PhrAddressLoginRequest{
		PhrAddress: "doc.test@sbx",
		TxnID:      "1e16b852-d54c-4ed3-bb13-c905df2cacfe",
	})
	if err != nil {
		log.Printf("Failed to login with address via mobile: %v", err)
		return
	}

	fmt.Printf("ABHA login completed. Response: %v\n", loginResp)
}

func exampleAadhaarRegistration(client *abdm.Client) {
	ctx := context.Background()
	headers := interfaces.Headers{
		UserID: "user123",
		HipID:  "hip456",
	}

	// Step 1: Initialize Aadhaar registration
	initResp, err := client.Registration().AadhaarInit(ctx, headers, registration.InitRequest{
		AadhaarNumber: "123456789012",
	})
	if err != nil {
		log.Printf("Failed to initialize Aadhaar registration: %v", err)
		return
	}

	fmt.Printf("Aadhaar registration initiated. Transaction ID: %s\n", initResp.TxnID)

	// Step 2: Verify Aadhaar OTP (in real scenario, user would provide OTP)
	verifyResp, err := client.Registration().AadhaarVerify(ctx, headers, registration.VerifyRequest{
		TxnID:  initResp.TxnID,
		OTP:    "123456",
		Mobile: "9876543210",
	})
	if err != nil {
		log.Printf("Failed to verify Aadhaar OTP: %v", err)
		return
	}

	fmt.Printf("Aadhaar verification completed. Skip state: %s\n", verifyResp.SkipState)

	// Step 3: Create ABHA address (if skip_state is abha_create)
	if verifyResp.SkipState == "abha_create" {
		createResp, err := client.Registration().AadhaarCreatePHR(ctx, headers, registration.CreateRequest{
			TxnID:       verifyResp.TxnID,
			AbhaAddress: "john.doe@abdm",
		})
		if err != nil {
			log.Printf("Failed to create ABHA address: %v", err)
			return
		}

		if createResp.Profile != nil {
			fmt.Printf("ABHA created successfully: %s\n", createResp.Profile.AbhaAddress)
		}
	}
}

func exampleMobileRegistration(client *abdm.Client) {
	ctx := context.Background()
	headers := interfaces.Headers{}

	// Step 1: Initialize mobile registration
	initResp, err := client.Registration().MobileInit(ctx, headers, registration.MobileInitRequest{
		MobileNumber: "8590680079",
	})
	if err != nil {
		log.Printf("Failed to initialize mobile registration: %v", err)
		return
	}

	fmt.Printf("Mobile registration initiated. Transaction ID: %s\n", initResp.TxnID)

	//// Step 2: Verify mobile OTP
	//verifyResp, err := client.Registration().MobileVerify(ctx, headers, registration.MobileVerifyOTPRequest{
	//	TxnID: initResp.TxnID,
	//	OTP:   "123456",
	//})
	//if err != nil {
	//	log.Printf("Failed to verify mobile OTP: %v", err)
	//	return
	//}
	//
	//fmt.Printf("Mobile verification completed. Skip state: %s\n", verifyResp.SkipState)
	//
	//// Step 3: Create ABHA address with profile details
	//createResp, err := client.Registration().MobileCreatePHR(ctx, headers, registration.MobileCreateRequest{
	//	TxnID:       verifyResp.TxnID,
	//	AbhaAddress: "jane.doe@abdm",
	//	Profile: registration.ProfileDetailsRequest{
	//		FirstName:    "Jane",
	//		LastName:     abdm.StringPtr("Doe"),
	//		Gender:       "F",
	//		YearOfBirth:  1990,
	//		MonthOfBirth: 5,
	//		DayOfBirth:   15,
	//		Pincode:      "110001",
	//	},
	//})
	//if err != nil {
	//	log.Printf("Failed to create ABHA address via mobile: %v", err)
	//	return
	//}
	//
	//if createResp.Success && createResp.Profile != nil {
	//	fmt.Printf("ABHA created successfully via mobile: %s\n", createResp.Profile.AbhaAddress)
	//}
}

func exampleUtilityFunctions(client *abdm.Client) {
	// Generate transaction ID
	txnID := client.Utils().GenerateTransactionID()
	fmt.Printf("Generated transaction ID: %s\n", txnID)

	// Validate Aadhaar number
	if err := client.Utils().ValidateAadhaarNumber("123456789012"); err != nil {
		fmt.Printf("Invalid Aadhaar number: %v\n", err)
	} else {
		fmt.Println("Valid Aadhaar number")
	}

	// Validate mobile number
	if err := client.Utils().ValidateMobileNumber("9876543210"); err != nil {
		fmt.Printf("Invalid mobile number: %v\n", err)
	} else {
		fmt.Println("Valid mobile number")
	}

	// Validate ABHA address
	if err := client.Utils().ValidateABHAAddress("john.doe@abdm"); err != nil {
		fmt.Printf("Invalid ABHA address: %v\n", err)
	} else {
		fmt.Println("Valid ABHA address")
	}

	// Format date
	formattedDate := client.Utils().FormatDate(1990, 5, 15)
	fmt.Printf("Formatted date: %s\n", formattedDate)

	// Retry with backoff example
	ctx := context.Background()
	err := client.Utils().RetryWithBackoff(ctx, func() error {
		// Simulate some operation that might fail
		return nil
	}, 3, time.Second)
	if err != nil {
		fmt.Printf("Retry failed: %v\n", err)
	} else {
		fmt.Println("Retry succeeded")
	}
}
