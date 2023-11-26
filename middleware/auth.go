package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quocphan74/gingo.git/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, _ := c.Cookie("jwt")
		fmt.Println(utils.Parsejwt(cookie))
		if _, err := utils.Parsejwt(cookie); err != nil {
			// If not authenticated, return a 401 Unauthorized response
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			// Abort the request processing
			c.Abort()
		}
		c.Next()
	}

}
