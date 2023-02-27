package main

import (
	"context"
	"log"
	"os"
	"speakeasy/internal/app"
	"speakeasy/internal/pkg/authentication"
	"speakeasy/internal/pkg/profile"
	"speakeasy/internal/pkg/trip"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/gin-contrib/cors"

	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
)

var ginLambda *ginadapter.GinLambda

func main() {
	godotenv.Load()

	authenticationService := authentication.NewAuthenticationService()
	tripService := trip.NewTripService()
	profileService := profile.NewProfileService()

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://amuel.org", "https://dev.amuel.org"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "User-Agent"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	server := app.NewServer(
		router,
		authenticationService,
		tripService,
		profileService,
	)

	if inLambda() {
		ginLambda = ginadapter.New(server.Routes())
		lambda.Start(Handler)
		return
	} else {
		server.Run()
	}
}

func inLambda() bool {
	if runtime, _ := os.LookupEnv("AWS_LAMBDA_RUNTIME_API"); runtime != "" {
		log.Println("Found Lambda environment.")
		return true
	} else {
		log.Println("Lambda environment not found.")
		return false
	}
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return ginLambda.ProxyWithContext(ctx, req)
}
