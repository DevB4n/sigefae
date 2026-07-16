package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireRole(role string) gin.HandlerFunc {

	return func(c *gin.Context) {

		value, exists := c.Get("rol")

		if !exists {

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "rol no encontrado",
			})

			c.Abort()
			return
		}

		if value.(string) != role {

			c.JSON(http.StatusForbidden, gin.H{
				"error": "no tiene permisos para realizar esta acción",
			})

			c.Abort()
			return
		}

		c.Next()
	}
}