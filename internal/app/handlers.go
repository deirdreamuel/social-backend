package app

import (
	"fmt"
	"net/http"
	"speakeasy/internal/pkg/authentication"

	"github.com/gin-gonic/gin"
)

func (s *Server) ApiStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		response := map[string]any{
			"status":  http.StatusOK,
			"message": "status check successful",
		}

		c.JSON(http.StatusOK, response)
	}
}

// func (s *Server) PutProfile() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Header("Content-Type", "application/json")

// 		// read and validate request body
// 		var request profile.Profile
// 		if err := c.Bind(&request); err != nil {
// 			response := map[string]any{
// 				"status":  http.StatusBadRequest,
// 				"message": fmt.Sprintf("Error: %s", err),
// 			}

// 			c.JSON(http.StatusBadRequest, response)
// 			return
// 		}

// 		// put profile request
// 		err := s.profileService.PutProfile(request)
// 		if err != nil {
// 			response := map[string]any{
// 				"status":  http.StatusBadRequest,
// 				"message": fmt.Sprintf("Error: %s", err),
// 			}

// 			c.JSON(http.StatusBadRequest, response)
// 			return
// 		}

// 		response := map[string]any{
// 			"status":  http.StatusOK,
// 			"message": "profile created",
// 		}

// 		c.JSON(http.StatusOK, response)
// 	}
// }

// func (s *Server) ReadProfile() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Header("Content-Type", "application/json")

// 		response := map[string]any{
// 			"status":  http.StatusOK,
// 			"message": "profile read successfully",
// 		}

// 		c.JSON(http.StatusOK, response)
// 	}
// }

func (s *Server) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		// read and validate request body
		var request authentication.LoginRequest
		if err := c.Bind(&request); err != nil {
			response := map[string]any{
				"status":  http.StatusBadRequest,
				"message": fmt.Sprintf("Bad Request"),
			}

			c.JSON(http.StatusBadRequest, response)
			return
		}

		login, err := s.authenticationService.Login(request)
		if err != nil {
			response := map[string]any{
				"status":  err.Code,
				"message": err.Reason,
			}

			c.JSON(err.Code, response)
			return
		}

		c.JSON(http.StatusOK, login)
	}
}

func (s *Server) Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		// read and validate request body
		var request authentication.SignupRequest
		if err := c.Bind(&request); err != nil {
			response := map[string]any{
				"status":  http.StatusBadRequest,
				"message": fmt.Sprintf("Bad Request"),
			}

			c.JSON(http.StatusBadRequest, response)
			return
		}

		signup, err := s.authenticationService.Signup(request)
		if err != nil || signup.Status != true {
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

		c.JSON(http.StatusOK, response)
	}
}
