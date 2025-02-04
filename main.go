package main

import (
	"GoJWT/controllers"
	"GoJWT/database"
	"GoJWT/initializers"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
}

func main() {
	database.ConnectToDatabase()
	r := gin.Default()
	r.POST("/signup", func(c *gin.Context) {
		controllers.Signup(c)
	})
	r.POST("/login", func(c *gin.Context) {
		controllers.Login(c)
	})
	r.Run()
}
