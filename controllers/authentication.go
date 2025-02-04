package controllers

import (
	"GoJWT/database"
	"GoJWT/models"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
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

func Login(c *gin.Context) {
	// Getting the collection
	coll := database.Client.Database("gojwt").Collection("user")

	// Getting User Data from the request body
	var body models.User
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	// Search the email in the DB
	var res models.User
	err := coll.FindOne(context.TODO(), bson.D{{Key: "email", Value: body.Email}}).Decode(&res)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No user found"})
		return
	}

	// Compare the password and the hash
	errors := bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(body.Password))
	if errors != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong Password"})
		return
	}

	// Generating a token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": res.Email,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{})
}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{"message": user})
}
