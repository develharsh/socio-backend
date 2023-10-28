package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/harshvsinghme/socio-backend.git/dbUtils"
	globals "github.com/harshvsinghme/socio-backend.git/global"
	"golang.org/x/crypto/bcrypt"
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

func HashPassword(password string) ([]byte, error) {
	// Generate a salt with a cost factor of 14
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func VerifyPassword(hashedPassword []byte, inputPassword string) error {
	// Compare the hashed password with the input password
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(inputPassword))
	return err
}

func GetValueFromContext(c *gin.Context, key string) string {
	value, exists := c.Get(key)
	if !exists {
		return ""
	}

	if strValue, ok := value.(string); ok {
		return strValue
	}

	return ""
}
