package app

import "github.com/gin-gonic/gin"

// Routes Gin function which contains api routes
func (s *Server) Routes() *gin.Engine {
	router := s.router

	// health check endpoint
	router.GET("/", s.HealthCheck())

	// version 1 apis
	v1 := router.Group("/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("signup", s.Signup())
			auth.POST("login", s.Login())
		}
	}

	return router
}
