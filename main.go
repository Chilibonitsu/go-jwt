package main

import (
	"go-jwt/controllers"
	"go-jwt/initializers"
	"go-jwt/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDB()
}

func main() {
	r := gin.Default()
	r.POST("/signup", controllers.SignUp)
	r.POST("/tokenize", controllers.Tokenize)

	r.POST("/refreshTokens", middleware.RequireAuth, controllers.RefreshTokens)
	r.Run()
}
