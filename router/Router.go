package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/harshvsinghme/socio-backend.git/middleware"
	"github.com/harshvsinghme/socio-backend.git/services"
)

func InitRouter() *gin.Engine {

	router := gin.Default()

	router.GET("/status", func(ctx *gin.Context) {
		// Create a JSON message
		message := gin.H{
			"message": "Hello, Application is running!",
		}

		// Respond with JSON
		ctx.JSON(http.StatusOK, message)
	})

	router.POST("/v1/auth-service/auth/signup", services.AuthService{}.Signup)

	router.POST("/v1/auth-service/auth/login", services.AuthService{}.Login)

	router.GET("/v1/user-service/account", middleware.SessionMiddleware(),
		services.AuthService{}.Account)

	return router
}
