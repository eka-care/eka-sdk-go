package main

import (
	"context"
	"fmt"

	"github.com/eka-care/eka-sdk-go/abdm"
	"github.com/eka-care/eka-sdk-go/abdm/abha/login"
	"github.com/eka-care/eka-sdk-go/abdm/abha/registration"
	"github.com/eka-care/eka-sdk-go/internal/interfaces"
)

func demoAPIs() {
	// Create ABDM client using environment variables (recommended)
	client := abdm.NewFromEnv()

	fmt.Println("=== ABDM API Implementation Examples ===")

	// Demo context and headers
	ctx := context.Background()
	headers := interfaces.Headers{
		UserID: "demo-user-123",
		HipID:  "demo-hip-456",
	}

	// Example 1: Login Flow
	fmt.Println("1. Login API Examples:")
	demoLoginAPIs(client, ctx, headers)

	// Example 2: Registration Flow
	fmt.Println("\n2. Registration API Examples:")
	demoRegistrationAPIs(client, ctx, headers)

	// Example 3: Utility APIs
	fmt.Println("\n3. Utility API Examples:")
	demoUtilityAPIs(client, ctx, headers)

	fmt.Println("\n=== All APIs demonstrated successfully! ===")
}

func demoLoginAPIs(client *abdm.Client, ctx context.Context, headers interfaces.Headers) {
	fmt.Println("  Login Init API:")
	initReq := &login.InitLoginRequest{
		Identifier: "demo@example.com",
		Method:     login.LoginMethodPhrAddress,
	}
	fmt.Printf("    Request: %+v\n", initReq)

	fmt.Println("  Login Verify API:")
	verifyReq := &login.VerifyLoginOTPRequest{
		TxnID: "demo-txn-id",
		OTP:   "123456",
	}
	fmt.Printf("    Request: %+v\n", verifyReq)

	fmt.Println("  PHR Address Login API:")
	phrReq := &login.PhrAddressLoginRequest{
		PhrAddress: "demo@abdm",
		TxnID:      "demo-txn-id",
	}
	fmt.Printf("    Request: %+v\n", phrReq)
}

func demoRegistrationAPIs(client *abdm.Client, ctx context.Context, headers interfaces.Headers) {
	fmt.Println("  Aadhaar Registration APIs:")

	// Aadhaar Init
	aadhaarInitReq := registration.InitRequest{
		AadhaarNumber: "123456789012",
	}
	fmt.Printf("    AadhaarInit Request: %+v\n", aadhaarInitReq)

	// Aadhaar Verify
	aadhaarVerifyReq := registration.VerifyRequest{
		TxnID:  "demo-txn-id",
		OTP:    "123456",
		Mobile: "9876543210",
	}
	fmt.Printf("    AadhaarVerify Request: %+v\n", aadhaarVerifyReq)

	// Aadhaar Create PHR
	aadhaarCreateReq := registration.CreateRequest{
		TxnID:       "demo-txn-id",
		AbhaAddress: "demo@abdm",
	}
	fmt.Printf("    AadhaarCreatePHR Request: %+v\n", aadhaarCreateReq)

	fmt.Println("  Mobile Registration APIs:")

	// Mobile Init
	mobileInitReq := registration.MobileInitRequest{
		MobileNumber: "9876543210",
	}
	fmt.Printf("    MobileInit Request: %+v\n", mobileInitReq)

	// Mobile Verify
	mobileVerifyReq := registration.MobileVerifyOTPRequest{
		TxnID: "demo-txn-id",
		OTP:   "123456",
	}
	fmt.Printf("    MobileVerify Request: %+v\n", mobileVerifyReq)

	// Mobile Create PHR
	mobileCreateReq := registration.MobileCreateRequest{
		TxnID:       "demo-txn-id",
		AbhaAddress: "demo@abdm",
		Profile: registration.ProfileDetailsRequest{
			FirstName:    "John",
			LastName:     stringPtr("Doe"),
			Gender:       "M",
			YearOfBirth:  1990,
			MonthOfBirth: 1,
			DayOfBirth:   15,
			Pincode:      "110001",
		},
	}
	fmt.Printf("    MobileCreatePHR Request: %+v\n", mobileCreateReq)
}

func demoUtilityAPIs(client *abdm.Client, ctx context.Context, headers interfaces.Headers) {
	fmt.Println("  Check ABHA Address Exists:")
	checkReq := registration.DoesHealthIdExistRequest{
		AbhaAddress: "demo@abdm",
	}
	fmt.Printf("    Request: %+v\n", checkReq)

	fmt.Println("  Suggest ABHA Address:")
	fmt.Printf("    Parameters: firstName=John, lastName=Doe, dob=1990-01-15, txnID=demo-txn-id\n")

	fmt.Println("  Available Service Methods:")
	fmt.Println("    Login Service:")
	fmt.Println("      - LoginInit()")
	fmt.Println("      - LoginVerify()")
	fmt.Println("      - LoginWithPHRAddress()")

	fmt.Println("    Registration Service:")
	fmt.Println("      Aadhaar APIs:")
	fmt.Println("        - AadhaarInit()")
	fmt.Println("        - AadhaarVerify()")
	fmt.Println("        - AadhaarResend()")
	fmt.Println("        - AadhaarMobileVerify()")
	fmt.Println("        - AadhaarMobileResend()")
	fmt.Println("        - AadhaarCreatePHR()")
	fmt.Println("      Mobile APIs:")
	fmt.Println("        - MobileInit()")
	fmt.Println("        - MobileVerify()")
	fmt.Println("        - MobileResend()")
	fmt.Println("        - MobileCreatePHR()")
	fmt.Println("      Utility APIs:")
	fmt.Println("        - CheckAbhaAddressExists()")
	fmt.Println("        - SuggestAbhaAddress()")
}

func stringPtr(s string) *string {
	return &s
}

// Example of how to use the APIs in practice (commented out to avoid network calls)
func exampleActualAPICall(client *abdm.Client, ctx context.Context, headers interfaces.Headers) {
	/*
		// Real API call example:
		initResp, err := client.Registration().AadhaarInit(ctx, headers, registration.InitRequest{
			AadhaarNumber: "123456789012",
		})
		if err != nil {
			log.Printf("Error: %v", err)
			return
		}

		fmt.Printf("Transaction ID: %s\n", initResp.TxnID)

		// Continue with verify...
		verifyResp, err := client.Registration().AadhaarVerify(ctx, headers, registration.VerifyRequest{
			TxnID:  initResp.TxnID,
			OTP:    "123456", // User-provided OTP
			Mobile: "9876543210",
		})
		if err != nil {
			log.Printf("Error: %v", err)
			return
		}

		fmt.Printf("Skip State: %s\n", verifyResp.SkipState)
	*/
}
