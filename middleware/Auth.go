package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/harshvsinghme/socio-backend.git/dbUtils"
	globals "github.com/harshvsinghme/socio-backend.git/global"
	"github.com/harshvsinghme/socio-backend.git/models"
	"github.com/harshvsinghme/socio-backend.git/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var RedisClient *redis.Client

func SessionMiddleware() gin.HandlerFunc {
	RedisClient = dbUtils.RedisClient
	return func(ctx *gin.Context) {
		// Extract the bearer token from the Authorization header
		authHeader := ctx.GetHeader("Authorization")
		accessToken := strings.TrimPrefix(authHeader, "Bearer ")

		//print all url header -request here

		var err error
		var userId, refreshToken string
		var userExists bool

		if accessToken == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Session expired, please log in again", "errors": []string{"Please log in again"}})
			ctx.Abort()
			return
		}
		userId, _ = validateAccessToken(accessToken)
		if userId != "" {
			userExists = validateUserExistence(userId)
			if !userExists {
				ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Session expired, please log in again", "errors": []string{"Please log in again"}})
				ctx.Abort()
				return
			}
			ctx.Set("userId", userId)
			ctx.Next()
			return
		}

		refreshToken, _ = getValueFromRedis(accessToken)
		if refreshToken == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Session expired, please log in again", "errors": []string{"Please log in again"}})
			ctx.Abort()
			return
		}

		userId, _ = validateRefreshToken(refreshToken)
		if userId == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Session expired, please log in again", "errors": []string{"Please log in again"}})
			ctx.Abort()
			return
		}

		//remove previous access and refresh token
		deleteKeyFromRedis(accessToken)
		//generate new access token
		accessToken, err = utils.GenerateAccessToken(userId)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Session expired, please log in again", "errors": []string{"Please log in again"}})
			ctx.Abort()
			return
		}
		//generate new refresh token
		refreshToken, err = utils.GenerateRefreshToken(userId)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Session expired, please log in again", "errors": []string{"Please log in again"}})
			ctx.Abort()
			return
		}

		//save new access token and refresh token
		err = utils.SetKeyValueWithExpiryToRedis(accessToken, refreshToken, time.Hour*24*30)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Session expired, please log in again", "errors": []string{"Please log in again"}})
			ctx.Abort()
			return
		}
		//set userId in context
		userExists = validateUserExistence(userId)
		if !userExists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Session expired, please log in again", "errors": []string{"Please log in again"}})
			ctx.Abort()
			return
		}
		ctx.Set("userId", userId)
		ctx.Next()
	}
}

func validateAccessToken(accessToken string) (string, error) {
	// Parse and validate the access token
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method and return the access token secret
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method)
		}
		return []byte(globals.SECRETS.JWT_SECRET), nil
	})

	if err != nil {
		return "", err
	}

	// Check if the token is valid
	if !token.Valid {
		return "", fmt.Errorf("invalid access token")
	}

	// Extract the user ID from the claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("failed to parse token claims")
	}

	userId, ok := claims["userId"].(string)
	if !ok {
		return "", fmt.Errorf("failed to extract user ID from token")
	}

	return userId, nil
}

func getValueFromRedis(key string) (string, error) {
	value, err := RedisClient.Get(key).Result()
	if err == redis.Nil {
		// Handle case when the key does not exist in Redis
		return "", fmt.Errorf("key does not exist in Redis")
	} else if err != nil {
		// Handle other error cases
		return "", err
	}
	return value, nil
}

func validateRefreshToken(refreshToken string) (string, error) {
	// Parse and validate the refresh token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method and return the refresh token secret
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method)
		}
		return []byte(globals.SECRETS.JWT_SECRET), nil
	})

	if err != nil {
		return "", err
	}

	// Check if the token is valid
	if !token.Valid {
		return "", fmt.Errorf("invalid refresh token")
	}

	// Extract the user ID from the claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("failed to parse token claims")
	}

	userId, ok := claims["userId"].(string)
	if !ok {
		return "", fmt.Errorf("failed to extract user ID from token")
	}

	return userId, nil
}

func deleteKeyFromRedis(key string) error {
	err := RedisClient.Del(key).Err()
	if err != nil {
		return err
	}

	return nil
}

func validateUserExistence(userId string) bool {
	var user models.User
	var err error
	var userObjId primitive.ObjectID
	userObjId, err = primitive.ObjectIDFromHex(userId)
	if err != nil {
		return false
	}
	err = dbUtils.MongoClient.Database(globals.SECRETS.MONGO_DATABASE).Collection("users").FindOne(context.TODO(), bson.M{"_id": userObjId}).Decode(&user)
	// fmt.Println("Hey", user.Name, user.Email)
	return err == nil
}
