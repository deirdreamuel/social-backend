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
			auth.POST("refresh", s.Refresh())
		}

		trip := v1.Group("/trip")
		{
			trip.POST("", s.CreateTrip())
			trip.GET("/:tripid", s.GetTrip())
			trip.GET("/:tripid/participants", s.GetTripParticipants())
		}

		user := v1.Group("/trips")
		{
			user.GET("/user/:userid", s.GetUserTrips())
			user.GET("/user/me", s.GetMyTrips())
		}

		profile := v1.Group("/profile")
		{
			profile.GET("", s.GetMyProfile())
			profile.POST("", s.CreateProfile())
			profile.POST("/picture", s.UploadProfilePicture())
		}
	}

	return router
}
