package services

import (
	"context"
	"testing"
)

func TestAuthService_Register(t *testing.T) {
	userRepo := NewMockUserRepo()
	pasienRepo := NewMockPasienRepo()
	waGateway := &StubWAGateway{}
	service := NewAuthService(userRepo, pasienRepo, waGateway)

	ctx := context.Background()

	// Test Successful Registration
	user, err := service.Register(ctx, "test@example.com", "password123", "Test Patient", "1990-01-01", "perempuan", "Jl. Test", "08123456789", "O")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.Email != "test@example.com" {
		t.Errorf("expected email to be test@example.com, got %s", user.Email)
	}

	// Verify Pasien record created
	pasien, err := pasienRepo.FindByUserID(ctx, user.ID)
	if err != nil {
		t.Fatalf("expected no error finding patient, got %v", err)
	}
	if pasien == nil {
		t.Fatal("expected patient record to exist")
	}
	if pasien.NamaLengkap != "Test Patient" {
		t.Errorf("expected name to be Test Patient, got %s", pasien.NamaLengkap)
	}

	// Test Duplicate Email
	_, err = service.Register(ctx, "test@example.com", "newpassword", "Another Name", "1995-05-05", "laki-laki", "Alamat", "08987654321", "A")
	if err == nil {
		t.Error("expected error for duplicate email, got nil")
	}
}

func TestAuthService_Login(t *testing.T) {
	userRepo := NewMockUserRepo()
	pasienRepo := NewMockPasienRepo()
	waGateway := &StubWAGateway{}
	service := NewAuthService(userRepo, pasienRepo, waGateway)

	ctx := context.Background()

	// Setup user
	_, err := service.Register(ctx, "login@example.com", "securepass", "User", "1992-02-02", "laki-laki", "Alamat", "081111", "")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	// Test Successful Login
	token, user, err := service.Login(ctx, "login@example.com", "securepass")
	if err != nil {
		t.Fatalf("expected successful login, got %v", err)
	}
	if token == "" {
		t.Error("expected non-empty token")
	}
	if user.Email != "login@example.com" {
		t.Errorf("expected user email login@example.com, got %s", user.Email)
	}

	// Test Failed Login - Incorrect Password
	_, _, err = service.Login(ctx, "login@example.com", "wrongpass")
	if err == nil {
		t.Error("expected error for wrong password, got nil")
	}

	// Test Failed Login - User Not Found
	_, _, err = service.Login(ctx, "unknown@example.com", "somepass")
	if err == nil {
		t.Error("expected error for non-existent user, got nil")
	}
}

func TestAuthService_ForgotPasswordAndReset(t *testing.T) {
	userRepo := NewMockUserRepo()
	pasienRepo := NewMockPasienRepo()
	waGateway := &StubWAGateway{}
	service := NewAuthService(userRepo, pasienRepo, waGateway)

	ctx := context.Background()

	// Setup User & Patient
	user, err := service.Register(ctx, "forgot@example.com", "oldpassword", "Forgot User", "1988-08-08", "perempuan", "Address", "082222", "")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	// Test Forgot Password
	err = service.ForgotPassword(ctx, "082222")
	if err != nil {
		t.Fatalf("expected no error for forgot password, got %v", err)
	}

	// Retrieve OTP from Mock Store
	otp, ok := userRepo.OTPStore[user.ID]
	if !ok {
		t.Fatal("expected OTP to be stored")
	}
	if len(otp.Code) != 6 {
		t.Errorf("expected OTP code of length 6, got %s", otp.Code)
	}

	// Test Reset Password with Wrong OTP
	err = service.ResetPassword(ctx, "082222", "000000", "newpassword123")
	if err == nil {
		t.Error("expected error for wrong OTP, got nil")
	}

	// Test Reset Password with Correct OTP
	err = service.ResetPassword(ctx, "082222", otp.Code, "newpassword123")
	if err != nil {
		t.Fatalf("expected successful password reset, got %v", err)
	}

	// Verify Login with New Password
	_, _, err = service.Login(ctx, "forgot@example.com", "newpassword123")
	if err != nil {
		t.Errorf("expected successful login with new password, got %v", err)
	}
}
