package app

import (
	"net/http"

	"speakeasy/internal/pkg/authentication"

	"github.com/gin-gonic/gin"
)

// GetUserTrips Gin handler function to get trip by trip id
func (s *Server) GetUserTrips() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		userID := c.Param("userid")

		trips, err := s.tripService.GetTripsByUser(userID)
		if err != nil {
			response := map[string]any{
				"status":  err.Code,
				"message": err.Reason,
			}

			c.JSON(err.Code, response)
			return
		}

		c.JSON(http.StatusOK, trips)
	}
}

// GetUserTrips Gin handler function to get trip by trip id
func (s *Server) GetMyTrips() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		// Get UserID from JWT
		claims, errauth := authentication.GetTokenClaims(c.Request)
		if errauth != nil {
			response := map[string]any{
				"status":  401,
				"message": "Token expired",
			}

			c.JSON(401, response)
			return
		}

		userID := ((*claims)["user_id"]).(string)

		trips, err := s.tripService.GetTripsByUser(userID)
		if err != nil {
			response := map[string]any{
				"status":  err.Code,
				"message": err.Reason,
			}

			c.JSON(err.Code, response)
			return
		}

		c.JSON(http.StatusOK, trips)
	}
}
