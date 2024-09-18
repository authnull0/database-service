package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthnzCall() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("X-Authorization")

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Redirect(http.StatusUnauthorized, "")
			return
		}

		result := Authnz(token, c, 0)

		if result == false {
			c.JSON(http.StatusUnauthorized, "Unauthorized")
			c.Redirect(http.StatusUnauthorized, "")
			return
		}

		c.Next()

	}
}
