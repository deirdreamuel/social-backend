package app

import (
	"net/http"
	"speakeasy/internal/pkg/trip"

	"github.com/gin-gonic/gin"
)

// CreateTrip Gin handler function to create trip
func (s *Server) CreateTrip() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		// read and validate request body
		var request trip.Trip
		if err := c.Bind(&request); err != nil {
			response := map[string]any{
				"status":  http.StatusBadRequest,
				"message": "Bad Request",
			}

			c.JSON(http.StatusBadRequest, response)
			return
		}

		err := s.tripService.CreateTrip(&request)
		if err != nil {
			response := map[string]any{
				"status":  err.Code,
				"message": err.Reason,
			}

			c.JSON(err.Code, response)
			return
		}

		response := map[string]any{
			"status":  http.StatusCreated,
			"message": "Created",
		}

		c.JSON(http.StatusCreated, response)
	}
}

// GetTrip Gin handler function to get trip by trip id
func (s *Server) GetTrip() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		tripid := c.Param("tripid")

		trip, err := s.tripService.GetTrip(tripid)
		if err != nil {
			response := map[string]any{
				"status":  err.Code,
				"message": err.Reason,
			}

			c.JSON(err.Code, response)
			return
		}

		c.JSON(http.StatusOK, trip)
	}
}