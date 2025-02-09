package utils

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const secretKey = "kjfhdkjhGHDFGgghfjgHJGgdhghjFGUYtgGDFY87687GUD74574RFUGFJHghiuY87557UGIUGDFYT757576GHGJHGJHddfdfher"

// A simple in-memory token blacklist (use Redis for production)
var (
	tokenBlacklist = make(map[string]bool)
	blacklistMutex sync.Mutex
)

// GenerateToken creates a JWT token for authentication
func GenerateToken(email string, userID int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 2).Unix(),
	})

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// VerifyToken parses and validates a JWT token
func VerifyToken(token string) (int64, error) {
	// Parse and validate the token
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		// Check if the signing method is HMAC
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("Unexpected signing method.")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return 0, errors.New("Could not parse token.")
	}

	// Ensure that the token is valid
	if !parsedToken.Valid {
		return 0, errors.New("Invalid token.")
	}

	// Extract the claims
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("Invalid token claims.")
	}

	// Check the expiration date
	expirationTime, ok := claims["exp"].(float64)
	if !ok {
		return 0, errors.New("Invalid expiration time.")
	}

	// If token is expired, reject it immediately
	if time.Unix(int64(expirationTime), 0).Before(time.Now()) {
		return 0, errors.New("Token has expired.")
	}

	// Check if the token is blacklisted
	blacklistMutex.Lock()
	if tokenBlacklist[token] {
		blacklistMutex.Unlock()
		return 0, errors.New("Token is revoked.")
	}
	blacklistMutex.Unlock()

	// If token is valid and not blacklisted, return the user ID
	userID := int64(claims["userID"].(float64))
	return userID, nil
}

// InvalidateToken adds a token to the blacklist
func InvalidateToken(token string) {
	blacklistMutex.Lock()
	tokenBlacklist[token] = true
	blacklistMutex.Unlock()
}