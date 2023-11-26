package services

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/harshvsinghme/socio-backend.git/dbUtils"
	globals "github.com/harshvsinghme/socio-backend.git/global"
	"github.com/harshvsinghme/socio-backend.git/models"
	"github.com/harshvsinghme/socio-backend.git/utils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthService struct {
}

func (service AuthService) Signup(ctx *gin.Context) {

	//Expected Body Structure
	type ReqBody struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	//Variables
	var reqBody ReqBody
	var err error
	var accessToken, refreshToken string
	var firstInsert *mongo.InsertOneResult
	var validation []string
	var userId primitive.ObjectID
	var result models.User

	//Validate Input
	if err = ctx.BindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save info", "errors": []string{"Invalid request body"}})
		return
	}

	//Validate Input data
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

	//Make sure all validations met
	if len(validation) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid details", "errors": validation})
		return
	}

	//Check if email already exists
	err = dbUtils.MongoClient.Database(globals.SECRETS.MONGO_DATABASE).Collection("users").FindOne(context.TODO(), bson.M{"email": reqBody.Email}).Decode(&result)
	if err == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Email already in use", "errors": []string{"Try another email or remove this from existing account"}})
		return
	}

	//insert new user data
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

	//insert new user data in userdetails collection against previously saved user's _id
	_, err = dbUtils.MongoClient.Database(globals.SECRETS.MONGO_DATABASE).Collection("userdetails").InsertOne(context.TODO(), models.UserDetail{UserId: userId})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save user details", "errors": []string{err.Error()}})
		return
	}

	//generate  access token
	accessToken, err = utils.GenerateAccessToken(userId.Hex())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not log you in, try again", "errors": []string{err.Error()}})
		return
	}

	//generate refresh token
	refreshToken, err = utils.GenerateRefreshToken(userId.Hex())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not log you in, try again", "errors": []string{err.Error()}})
		return
	}

	//save access token and refresh token in redis as a mapping(access <-> refresh)
	err = utils.SetKeyValueWithExpiryToRedis(accessToken, refreshToken, time.Hour*24*30)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not log you in, try again", "errors": []string{err.Error()}})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Registered Successfully", "accessToken": accessToken})
}

func (service AuthService) Login(ctx *gin.Context) {

	//Expected Body Structure
	type ReqBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	//Variables
	var reqBody ReqBody
	var err error
	var validation []string
	var accessToken, refreshToken string
	var result models.User

	//Validate Input
	if err = ctx.BindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to sign in", "errors": []string{"Invalid request body"}})
	}

	//Validate Input data
	if reqBody.Email == "" {
		validation = append(validation, "Email can't be empty")
	}
	if reqBody.Password == "" {
		validation = append(validation, "Password can't be empty")
	}

	//Make sure all validations met
	if len(validation) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid details", "errors": validation})
		return
	}

	//Validate User existence
	err = dbUtils.MongoClient.Database(globals.SECRETS.MONGO_DATABASE).Collection("users").FindOne(context.TODO(), bson.M{"email": reqBody.Email}).Decode(&result)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "No such user found", "errors": []string{"Either email or password is wrong"}})
		return
	}

	//Verify password
	err = utils.VerifyPassword([]byte(result.Password), reqBody.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "No such user found", "errors": []string{"Either email or password is wrong"}})
		return
	}

	//generate  access token
	accessToken, err = utils.GenerateAccessToken(result.DocId.Hex())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not log you in, try again", "errors": []string{err.Error()}})
		return
	}

	//generate refresh token
	refreshToken, err = utils.GenerateRefreshToken(result.DocId.Hex())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not log you in, try again", "errors": []string{err.Error()}})
		return
	}

	//save access token and refresh token in redis as a mapping(access <-> refresh)
	err = utils.SetKeyValueWithExpiryToRedis(accessToken, refreshToken, time.Hour*24*30)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not log you in, try again", "errors": []string{err.Error()}})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Logged you in successfully", "accessToken": accessToken})

}

func (service AuthService) Account(ctx *gin.Context) {
	userId := utils.GetValueFromContext(ctx, "userId")
	type UserAccDetails struct {
		UserId    primitive.ObjectID `bson:"userId,omitempty" json:"userId,omitempty"`
		Name      string             `bson:"name" json:"name"`
		Email     string             `bson:"email" json:"email"`
		Following int32              `bson:"following" json:"following"`
		Followers int32              `bson:"followers" json:"followers"`
	}
	var user UserAccDetails
	var cursor *mongo.Cursor
	var err error
	var userObjId primitive.ObjectID
	userObjId, _ = primitive.ObjectIDFromHex(userId)

	cursor, err = dbUtils.MongoClient.Database(globals.SECRETS.MONGO_DATABASE).Collection("users").Aggregate(context.TODO(), utils.PipelineUserAccount(userObjId))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"UserId": userId,
		})
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Session expired, please log in again", "errors": []string{"Please log in again"}})
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		if err := cursor.Decode(&user); err != nil {
			logrus.Error("error while processing cursor call", fmt.Sprintf("%+v", errors.Wrap(err, err.Error())))
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"userData": user})
}
