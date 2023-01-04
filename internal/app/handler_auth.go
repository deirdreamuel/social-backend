package app

import (
	"net/http"
	"speakeasy/internal/pkg/authentication"

	"github.com/gin-gonic/gin"
)

// Login Gin handler function to login and get access token
func (s *Server) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		// read and validate request body
		var request authentication.LoginRequest
		if err := c.Bind(&request); err != nil {
			response := map[string]any{
				"status":  http.StatusBadRequest,
				"message": "Bad Request",
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

// Signup Gin handler function to signup user
func (s *Server) Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		// read and validate request body
		var request authentication.SignupRequest
		if err := c.Bind(&request); err != nil {
			response := map[string]any{
				"status":  http.StatusBadRequest,
				"message": "Bad Request",
			}

			c.JSON(http.StatusBadRequest, response)
			return
		}

		signup, err := s.authenticationService.Signup(request)
		if err != nil || !signup.Status {
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
