package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestAuthService_HashPassword(t *testing.T) {
	authService := NewAuthService("test-secret")

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "short password",
			password: "123",
			wantErr:  true,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := authService.HashPassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && hash == "" {
				t.Errorf("HashPassword() returned empty hash for valid password")
			}
		})
	}
}

func TestAuthService_VerifyPassword(t *testing.T) {
	authService := NewAuthService("test-secret")
	password := "password123"
	
	hash, err := authService.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	tests := []struct {
		name           string
		hashedPassword string
		password       string
		wantErr        bool
	}{
		{
			name:           "correct password",
			hashedPassword: hash,
			password:       password,
			wantErr:        false,
		},
		{
			name:           "incorrect password",
			hashedPassword: hash,
			password:       "wrongpassword",
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := authService.VerifyPassword(tt.hashedPassword, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAuthService_GenerateAndValidateToken(t *testing.T) {
	authService := NewAuthService("test-secret")
	userID := 123
	username := "testuser"
	isAdmin := false

	// Generate token
	token, err := authService.GenerateToken(userID, username, isAdmin)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	if token == "" {
		t.Errorf("GenerateToken() returned empty token")
	}

	// Validate token
	claims, err := authService.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("ValidateToken() userID = %v, want %v", claims.UserID, userID)
	}

	if claims.Username != username {
		t.Errorf("ValidateToken() username = %v, want %v", claims.Username, username)
	}

	if claims.IsAdmin != isAdmin {
		t.Errorf("ValidateToken() isAdmin = %v, want %v", claims.IsAdmin, isAdmin)
	}
}

func TestAuthService_ValidateToken_InvalidToken(t *testing.T) {
	authService := NewAuthService("test-secret")

	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "empty token",
			token: "",
		},
		{
			name:  "invalid token",
			token: "invalid.token.here",
		},
		{
			name:  "malformed token",
			token: "not-a-jwt-token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := authService.ValidateToken(tt.token)
			if err == nil {
				t.Errorf("ValidateToken() expected error for invalid token")
			}
		})
	}
}

func TestAuthService_ValidateToken_ExpiredToken(t *testing.T) {
	authService := NewAuthService("test-secret")

	// Create an expired token manually
	claims := Claims{
		UserID:   123,
		Username: "testuser",
		IsAdmin:  false,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expired 1 hour ago
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatalf("Failed to create expired token: %v", err)
	}

	// Try to validate expired token
	_, err = authService.ValidateToken(tokenString)
	if err == nil {
		t.Errorf("ValidateToken() expected error for expired token")
	}
}

func TestAuthService_GenerateAndValidatePasswordResetToken(t *testing.T) {
	authService := NewAuthService("test-secret")
	userID := 123
	email := "test@example.com"

	// Generate reset token
	token, err := authService.GeneratePasswordResetToken(userID, email)
	if err != nil {
		t.Fatalf("GeneratePasswordResetToken() error = %v", err)
	}

	if token == "" {
		t.Errorf("GeneratePasswordResetToken() returned empty token")
	}

	// Validate reset token
	validatedUserID, validatedEmail, err := authService.ValidatePasswordResetToken(token)
	if err != nil {
		t.Fatalf("ValidatePasswordResetToken() error = %v", err)
	}

	if validatedUserID != userID {
		t.Errorf("ValidatePasswordResetToken() userID = %v, want %v", validatedUserID, userID)
	}

	if validatedEmail != email {
		t.Errorf("ValidatePasswordResetToken() email = %v, want %v", validatedEmail, email)
	}
}