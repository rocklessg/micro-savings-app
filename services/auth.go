package services

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

// GenerateJWT generates a JWT token for a given user ID
func GenerateJWT(userID string) (string, error) {
	secret := []byte(getJWTSecret())

	// Define the token claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret
	return token.SignedString(secret)
}

// ValidateJWT validates the provided token and extracts claims
func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	secret := []byte(getJWTSecret())

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	// Extract and return the claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func getJWTSecret() string {
	// Load environment variables
	err := godotenv.Load("../.env") // Remove the path if the .env file is in the same directory (after unit test)
	if err != nil {
		panic("Failed to load .env file")
	}
	return os.Getenv("JWT_SECRET") // Replace with the value from your .env file
}
