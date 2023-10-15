package services

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/harshvsinghme/socio-backend.git/dbUtils"
	globals "github.com/harshvsinghme/socio-backend.git/global"
	"github.com/harshvsinghme/socio-backend.git/models"
	"github.com/harshvsinghme/socio-backend.git/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthService struct {
}

func (service AuthService) Signup(ctx *gin.Context) {

	type ReqBody struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var reqBody ReqBody
	var err error
	var accessToken, refreshToken string
	var firstInsert *mongo.InsertOneResult
	var validation []string
	var userId primitive.ObjectID
	var result bson.M

	if err = ctx.BindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save info", "errors": []string{"Invalid request body"}})
		return
	}

	if reqBody.Name == "" {
		validation = append(validation, "Name can't be empty")
	}
	if reqBody.Email == "" {
		validation = append(validation, "Email can't be empty")
	}
	if reqBody.Password == "" {
		validation = append(validation, "Password can't be empty")
	} else {
		var hashed []byte
		hashed, err = utils.HashPassword(reqBody.Password)
		if err != nil {
			validation = append(validation, "Password is invalid")
		}
		reqBody.Password = string(hashed)
	}

	if len(validation) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid details", "errors": validation})
		return
	}

	err = dbUtils.MongoClient.Database(globals.SECRETS.MONGO_DATABASE).Collection("users").FindOne(context.TODO(), bson.M{"email": reqBody.Email}).Decode(&result)
	if err == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Email already in use", "errors": []string{"Try another email or remove this from existing account"}})
		return
	}

	firstInsert, err = dbUtils.MongoClient.Database(globals.SECRETS.MONGO_DATABASE).Collection("users").InsertOne(context.TODO(), reqBody)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save user", "errors": []string{err.Error()}})
		return
	}
	userId = firstInsert.InsertedID.(primitive.ObjectID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save user details", "errors": []string{err.Error()}})
		return
	}
	_, err = dbUtils.MongoClient.Database(globals.SECRETS.MONGO_DATABASE).Collection("userdetails").InsertOne(context.TODO(), models.UserDetail{UserId: userId})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save user details", "errors": []string{err.Error()}})
		return
	}
	//generate  access token
	accessToken, err = utils.GenerateAccessToken(userId.Hex())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not sign you in, try to login now", "errors": []string{err.Error()}})
		return
	}

	//generate refresh token
	refreshToken, err = utils.GenerateRefreshToken(userId.Hex())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not sign you in, try to login now", "errors": []string{err.Error()}})
		return
	}
	//save access token and refresh token
	err = utils.SetKeyValueWithExpiryToRedis(accessToken, refreshToken, time.Hour*24*30)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not sign you in, try to login now", "errors": []string{err.Error()}})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Registered Successfully", "accessToken": accessToken, "errors": []string{}})
}

func (service AuthService) Login(ctx *gin.Context) {
}
