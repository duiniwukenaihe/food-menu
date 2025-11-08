package auth

import (
	"testing"
	"example.com/app/internal/models"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword123"
	hash, err := HashPassword(password)
	
	if err != nil {
		t.Errorf("HashPassword() returned error: %v", err)
	}
	
	if hash == "" {
		t.Error("HashPassword() returned empty hash")
	}
	
	if hash == password {
		t.Error("HashPassword() should not return the same password")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "testpassword123"
	hash, _ := HashPassword(password)
	
	// Test correct password
	if !CheckPasswordHash(password, hash) {
		t.Error("CheckPasswordHash() should return true for correct password")
	}
	
	// Test incorrect password
	if CheckPasswordHash("wrongpassword", hash) {
		t.Error("CheckPasswordHash() should return false for incorrect password")
	}
	
	// Test empty password
	if CheckPasswordHash("", hash) {
		t.Error("CheckPasswordHash() should return false for empty password")
	}
}

func TestGenerateToken(t *testing.T) {
	user := &models.User{
		ID:       1,
		Username: "testuser",
		Role:     "user",
	}
	
	token, err := GenerateToken(user)
	
	if err != nil {
		t.Errorf("GenerateToken() returned error: %v", err)
	}
	
	if token == "" {
		t.Error("GenerateToken() returned empty token")
	}
}

func TestValidateToken(t *testing.T) {
	user := &models.User{
		ID:       1,
		Username: "testuser",
		Role:     "user",
	}
	
	token, _ := GenerateToken(user)
	
	// Test valid token
	claims, err := ValidateToken(token)
	if err != nil {
		t.Errorf("ValidateToken() returned error for valid token: %v", err)
	}
	
	if claims.UserID != user.ID {
		t.Errorf("Expected UserID %d, got %d", user.ID, claims.UserID)
	}
	
	if claims.Username != user.Username {
		t.Errorf("Expected Username %s, got %s", user.Username, claims.Username)
	}
	
	if claims.Role != user.Role {
		t.Errorf("Expected Role %s, got %s", user.Role, claims.Role)
	}
	
	// Test invalid token
	_, err = ValidateToken("invalid.token.here")
	if err == nil {
		t.Error("ValidateToken() should return error for invalid token")
	}
	
	// Test empty token
	_, err = ValidateToken("")
	if err == nil {
		t.Error("ValidateToken() should return error for empty token")
	}
}
