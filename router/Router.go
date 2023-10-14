package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

	return router
}
