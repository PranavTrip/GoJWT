package middleware

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
)

func RequireAuth(c *gin.Context) {
	// Get cookie from request
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token Expired"})
		}

		// Getting the collection
		coll := database.Client.Database("gojwt").Collection("user")
		// Search the email in the DB
		var res models.User
		err := coll.FindOne(context.TODO(), bson.D{{Key: "email", Value: claims["sub"]}}).Decode(&res)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No user found"})
			return
		}
		c.Set("user", res.Name)
		fmt.Println(res)
		c.Next()

	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Expired Token"})
		return
	}
}
