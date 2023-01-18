package app

import (
	"log"
	"net/http"
	"speakeasy/internal/pkg/authentication"
	"speakeasy/internal/pkg/profile"

	"github.com/gin-gonic/gin"
)

// UploadProfilePicture Gin handler function to create trip
func (s *Server) UploadProfilePicture() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, _, _ := c.Request.FormFile("profile_pic")

		// Get UserID from JWT
		claims, errauth := authentication.GetTokenClaims(c.Request)
		if errauth != nil {
			log.Printf("(hander.UploadProfilePicture) error: %s", errauth)
			response := map[string]any{
				"status":  401,
				"message": "Token expired",
			}

			c.JSON(401, response)
			return
		}
		userID := ((*claims)["user_id"]).(string)

		if err := s.profileService.UploadProfilePicture(userID, file); err != nil {
			log.Printf("(hander.UploadProfilePicture) error: %s", err)
			response := map[string]any{
				"status":  http.StatusBadRequest,
				"message": "Bad Request",
			}

			c.JSON(http.StatusBadRequest, response)
			return
		}

		// File saved successfully. Return proper result
		c.JSON(http.StatusOK, gin.H{
			"message": "Your profile picture has been successfully updated.",
		})
	}
}

// CreateProfile handler function
func (s *Server) CreateProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request profile.Profile

		// Get UserID from JWT
		claims, errauth := authentication.GetTokenClaims(c.Request)
		if errauth != nil {
			log.Printf("(hander.CreateProfile) error: %s", errauth)
			response := map[string]any{
				"status":  401,
				"message": "Token expired",
			}

			c.JSON(401, response)
			return
		}
		userID := ((*claims)["user_id"]).(string)

		if err := c.Bind(&request); err != nil {
			log.Printf("(hander.CreateProfile) error: %s", err)
			response := map[string]any{
				"status":  http.StatusBadRequest,
				"message": "Bad Request",
			}

			c.JSON(http.StatusBadRequest, response)
			return
		}

		if err := s.profileService.PutProfile(userID, &request); err != nil {
			log.Printf("(hander.CreateProfile) error: %v", err)
			response := map[string]any{
				"status":  err.Code,
				"message": err.Reason,
			}

			c.JSON(err.Code, response)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Your profile has been successfully updated.",
		})
	}
}

// GetProfile handler function
func (s *Server) GetMyProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get UserID from JWT
		claims, errauth := authentication.GetTokenClaims(c.Request)
		if errauth != nil {
			log.Printf("(hander.GetProfile) error: %s", errauth)
			response := map[string]any{
				"status":  401,
				"message": "Token expired",
			}

			c.JSON(401, response)
			return
		}
		userID := ((*claims)["user_id"]).(string)

		profile, err := s.profileService.GetProfile(userID)
		if err != nil {
			log.Printf("(hander.GetProfile) error: %v", err)
			response := map[string]any{
				"status":  err.Code,
				"message": err.Reason,
			}

			c.JSON(err.Code, response)
			return
		}

		c.JSON(http.StatusOK, profile)
	}
}
