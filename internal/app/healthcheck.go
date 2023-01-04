package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck Gin handler function to check api health
func (s *Server) HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		response := map[string]any{
			"status":  http.StatusOK,
			"message": "health check successful",
		}

		c.JSON(http.StatusOK, response)
	}
}
