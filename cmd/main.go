package main

import (
	"speakeasy/internal/app"
	"speakeasy/internal/pkg/authentication"
	"speakeasy/internal/pkg/profile"
	"speakeasy/internal/pkg/trip"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/gin-contrib/cors"
)

func main() {
	godotenv.Load()

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://amuel.org", "https://dev.amuel.org"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "User-Agent"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Create services
	authenticationService := authentication.NewAuthenticationService()
	tripService := trip.NewTripService()
	profileService := profile.NewProfileService()

	server := app.NewServer(
		router,
		authenticationService,
		tripService,
		profileService,
	)

	server.Run()
}
