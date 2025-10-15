package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)

	if err != nil {
		return "", err
	}

	return hashedPassword, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	valid, err := argon2id.ComparePasswordAndHash(password, hash)
	fmt.Printf("%v", valid)

	if err != nil {
		return valid, err
	}

	return valid, nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claim := &jwt.RegisteredClaims{
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Issuer:    "chirp",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	tokenString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	fmt.Printf("ssss%v", token)

	if err != nil {
		fmt.Printf("JWT parsing error: %v\n", err)
		return uuid.Nil, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		fmt.Printf("JWT validated successfully for user: %s\n", claims.Subject)
		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			return uuid.Nil, err
		}
		return userID, nil
	} else {
		fmt.Printf("Token validation failed - ok: %v, valid: %v\n", ok, token.Valid)
		return uuid.Nil, fmt.Errorf("invalid token")
	}
}

func GetBearerToken(headers http.Header) (string, error) {
	authorization := headers.Get("Authorization")

	if authorization == "" {
		return "", fmt.Errorf("authorization is invalid")
	}

	tokenSplit := strings.Split(authorization, "Bearer ")
	if (len(tokenSplit)) < 2 {
		return "", fmt.Errorf("authorization is invalid")
	}
	token := tokenSplit[1]

	return token, nil
}

func MakeRefreshToken() (string, error) {
	key := make([]byte, 32)
	rand.Read(key)
	rf := hex.EncodeToString(key)
	return rf, nil
}

func GetAPIKey(headers http.Header) (string, error) {
	authorization := headers.Get("Authorization")

	if authorization == "" {
		return "", fmt.Errorf("authorization is invalid")
	}

	apiKeySplit := strings.Split(authorization, "ApiKey ")
	if len(apiKeySplit) < 2 {
		return "", fmt.Errorf("authorization is invalid")
	}
	apiKey := strings.TrimSpace(apiKeySplit[1])

	return apiKey, nil
}
