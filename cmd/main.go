package main

import (
	"speakeasy/internal/app"
	"speakeasy/internal/pkg/authentication"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/gin-contrib/cors"
)

func main() {
	godotenv.Load()

	router := gin.Default()
	router.Use(cors.Default())

	// Create services
	authenticationSvc := authentication.NewAuthenticationService()

	server := app.NewServer(
		router,
		authenticationSvc,
	)

	server.Run()
}
