// Package helpers provides utility functions for the accounts service.
package helpers

import (
	"fmt"
	"time"

	"api.workzen.odoo/constants"

	"github.com/dgrijalva/jwt-go"
)

var (
	ErrExpiredToken     = fmt.Errorf("token has expired")
	ErrInvalidToken     = fmt.Errorf("invalid token")
	ErrTokenBeforeIssue = fmt.Errorf("token used before issue time")
	ErrInvalidUserID    = fmt.Errorf("invalid user ID in token")

	JWTExpireDuration = time.Duration(constants.EncryptionJWTExpire) * time.Hour * 24
	JWTKey            = []byte(constants.EncryptionJWTSecret)
)

func GenerateJWT(payload map[string]any, expire time.Time) (string, error) {
	if expire.Before(time.Now()) {
		return "", fmt.Errorf("expiration time must be in the future")
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"exp": expire.Unix(),
		"iat": now.Unix(),
		"nbf": now.Unix(),
	}

	// Merge user payload
	for k, v := range payload {
		claims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString(JWTKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func VerifyJWT(tokenString string) (bool, map[string]any, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return JWTKey, nil
	})

	if err != nil {
		return false, nil, fmt.Errorf("parse error: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return false, nil, ErrInvalidToken
	}

	// Verify expiry and issued-at time
	if !verifyExpiry(claims) {
		return false, nil, ErrExpiredToken
	}

	return true, claims, nil
}

func verifyExpiry(claims jwt.MapClaims) bool {
	exp, expOk := claims["exp"].(float64)
	iat, iatOk := claims["iat"].(float64)

	if !expOk || !iatOk {
		return false
	}

	now := time.Now().Unix()
	return now <= int64(exp) && now >= int64(iat)
}
