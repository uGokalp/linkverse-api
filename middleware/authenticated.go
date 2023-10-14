package middleware

import (
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Authenticated(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	if tokenString == "" {
		c.AbortWithStatusJSON(401, gin.H{"message": "Unauthorized"})
		return
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{"message": "Unauthorized"})
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatusJSON(401, gin.H{"message": "Token expired"})
		}
		c.Set("userId", claims["sub"])
		c.Next()
	} else {
		c.AbortWithStatusJSON(401, gin.H{"message": "Unauthorized"})
		return

	}

}
