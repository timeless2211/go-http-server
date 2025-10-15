package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name          string
		password      string
		hash          string
		wantErr       bool
		matchPassword bool
	}{
		{
			name:          "Correct password",
			password:      password1,
			hash:          hash1,
			wantErr:       false,
			matchPassword: true,
		},
		{
			name:          "Incorrect password",
			password:      "wrongPassword",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Password doesn't match different hash",
			password:      password1,
			hash:          hash2,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Empty password",
			password:      "",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Invalid hash",
			password:      password1,
			hash:          "invalidhash",
			wantErr:       true,
			matchPassword: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && match != tt.matchPassword {
				t.Errorf("CheckPasswordHash() expects %v, got %v", tt.matchPassword, match)
			}
		})
	}
}

func TestGenerateToken(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "maxscecret"
	expiresIn := time.Hour

	tokenString, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantErr     bool
		userID      uuid.UUID
	}{
		{
			name:        "Valid token",
			tokenString: tokenString,
			tokenSecret: tokenSecret,
			wantErr:     false,
			userID:      userID,
		},
		{
			name:        "Invalid secret",
			tokenString: tokenString,
			tokenSecret: "wrongsecret",
			wantErr:     true,
		},
		{
			name:        "Empty token",
			tokenString: "",
			tokenSecret: tokenSecret,
			wantErr:     true,
		},
		{
			name:        "Malformed token",
			tokenString: "not.a.valid.jwt",
			tokenSecret: tokenSecret,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && match != tt.userID {
				t.Errorf("ValidateJWT() expects %v, got %v", tt.userID, match)
			}
		})
	}
}

func TestReadBearerToken(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "maxscecret"
	expiresIn := time.Hour

	tokenString, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	test := []struct {
		name    string
		headers http.Header
		wantErr bool
	}{
		{
			name: "Valid token",
			headers: http.Header{
				"Authorization": {"Bearer " + tokenString},
			},
			wantErr: false,
		},
		{
			name: "Invalid token",
			headers: http.Header{
				"Authorization": {tokenString},
			},
			wantErr: true,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GetBearerToken(tt.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && token != tokenString {
				t.Errorf("GetBearerToken() expects %v, got %v", token, tokenString)
			}
		})
	}
}
