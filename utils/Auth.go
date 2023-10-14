package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/harshvsinghme/socio-backend.git/dbUtils"
	globals "github.com/harshvsinghme/socio-backend.git/global"
)

func GenerateAccessToken(userId string) (string, error) {
	// Create the claims for the access token
	claims := jwt.MapClaims{
		"exp":    time.Now().Add(5 * 24 * time.Hour).Unix(), // Access token expiration time
		"userId": userId,
	}

	// Create the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the access token secret
	accessToken, err := token.SignedString([]byte(globals.SECRETS.JWT_SECRET))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func GenerateRefreshToken(userId string) (string, error) {
	// Create the claims for the refresh token
	claims := jwt.MapClaims{
		"exp":    time.Now().Add(30 * 24 * time.Hour).Unix(), // Refresh token expiration time
		"userId": userId,
	}

	// Create the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the refresh token secret
	refreshToken, err := token.SignedString([]byte(globals.SECRETS.JWT_SECRET))
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

func SetKeyValueWithExpiryToRedis(key string, value string, expiration time.Duration) error {
	RedisClient := dbUtils.RedisClient
	err := RedisClient.Set(key, value, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}
