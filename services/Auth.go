package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/harshvsinghme/socio-backend.git/dbUtils"
	globals "github.com/harshvsinghme/socio-backend.git/global"
	"github.com/harshvsinghme/socio-backend.git/models"
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
	var accessToken string
	var firstInsert *mongo.InsertOneResult

	if err = ctx.BindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	firstInsert, err = dbUtils.MongoClient.Database(globals.SECRETS.MONGO_DATABASE).Collection("users").InsertOne(context.TODO(), reqBody)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save user", "error": err.Error()})
		return
	}
	fmt.Println(firstInsert.InsertedID)
	_, err = dbUtils.MongoClient.Database(globals.SECRETS.MONGO_DATABASE).Collection("userdetails").InsertOne(context.TODO(), models.UserDetail{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save user details", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Registered Successfully", "accessToken": accessToken})
}

func (service AuthService) Login(ctx *gin.Context) {
}
