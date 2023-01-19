package app

import (
	"log"
	"net/http"
	"speakeasy/internal/pkg/authentication"
	"speakeasy/internal/pkg/profile"
	"speakeasy/pkg"

	"github.com/gin-gonic/gin"
)

// Login Gin handler function to login and get access token
func (s *Server) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		// read and validate request body
		var request authentication.LoginRequest
		if err := c.Bind(&request); err != nil {
			LogAndSendErrorResponse(c, &pkg.Error{
				Code:   http.StatusBadRequest,
				Reason: "Bad Request",
			})
			return
		}

		login, err := s.authenticationService.Login(request)
		if err != nil {
			LogAndSendErrorResponse(c, err)
			return
		}

		c.JSON(http.StatusOK, login)
	}
}

// Signup Gin handler function to signup user
func (s *Server) Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		var request authentication.SignupRequest
		if err := c.Bind(&request); err != nil {
			LogAndSendErrorResponse(c, &pkg.Error{
				Code:   http.StatusBadRequest,
				Reason: "Bad Request",
			})
			return
		}

		signup, err := s.authenticationService.Signup(request)
		if err != nil || !signup.Status {
			LogAndSendErrorResponse(c, err)
			return
		}

		err = s.profileService.PutProfile(&profile.Profile{
			UserID: signup.UserID,
			Name:   request.Name,
		})

		if err != nil {
			LogAndSendErrorResponse(c, err)
			return
		}

		response := map[string]any{
			"status":  http.StatusCreated,
			"message": "Created",
		}

		c.JSON(http.StatusCreated, response)
	}
}

// Signup Gin handler function to signup user
func (s *Server) Refresh() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		var request authentication.RefreshRequest
		if err := c.Bind(&request); err != nil {
			LogAndSendErrorResponse(c, &pkg.Error{
				Code:   http.StatusBadRequest,
				Reason: "Bad Request",
			})
			return
		}

		token, err := s.authenticationService.Refresh(&request)
		if err != nil {
			LogAndSendErrorResponse(c, &pkg.Error{
				Code:   http.StatusUnauthorized,
				Reason: "Unauthorized",
			})
			return
		}

		c.JSON(http.StatusOK, token)
	}
}

func LogAndSendErrorResponse(c *gin.Context, err *pkg.Error) {
	log.Printf("(%s) error: Reason: %s, Code: %d", c.Request.URL, err.Reason, err.Code)
	response := map[string]any{
		"status":  err.Code,
		"message": err.Reason,
	}

	c.JSON(err.Code, response)
}
