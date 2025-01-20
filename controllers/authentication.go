package controllers

import (
	"GoJWT/database"
	"GoJWT/models"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	// Getting the collection
	coll := database.Client.Database("gojwt").Collection("user")

	// Getting User Data from the request body
	var body models.User
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	// Hash the Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to hash password"})

	}

	// Inserting in DB
	user := models.User{Name: body.Name, Email: body.Email, Password: string(hashedPassword)}
	result, err := coll.InsertOne(context.TODO(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user"})
	}
	fmt.Printf("User Created %v", result)

	// Returning response
	c.JSON(http.StatusOK, gin.H{"message": "User Created"})

}
