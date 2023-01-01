package app

import "github.com/gin-gonic/gin"

func (s *Server) Routes() *gin.Engine {
	router := s.router

	// version 1 apis
	v1 := router.Group("/v1")
	{
		v1.GET("/status", s.ApiStatus())

		// prefix the user routes
		// user := v1.Group("/profile")
		// {
		// 	user.POST("", s.PutProfile())
		// 	user.GET("", s.ReadProfile())
		// }

		auth := v1.Group("/auth")
		{
			auth.POST("signup", s.Signup())
			auth.POST("login", s.Login())
		}
	}

	return router
}
