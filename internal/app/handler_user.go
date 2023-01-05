package app

import (
	"net/http"

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
