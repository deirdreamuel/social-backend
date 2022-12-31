package main

import (
	"os"
	"speakeasy/internal/app"
	"speakeasy/internal/pkg/authentication"
	"speakeasy/internal/pkg/profile"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/gin-contrib/cors"
)

func main() {
	godotenv.Load()

	router := gin.Default()
	router.Use(cors.Default())

	// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-credentials
	// Initialize session and config for initializing client
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	region := os.Getenv("AWS_REGION")
	endpoint := os.Getenv("DYNAMODB_ENDPOINT")

	cfg := &aws.Config{
		Region:   &region,
		Endpoint: &endpoint,
	}

	// Create new DynamoDB client
	ddb := dynamodb.New(sess, cfg)

	// Create services
	profileSvc := profile.NewProfileService(ddb, "Profile")
	authenticationSvc := authentication.NewAuthenticationService(ddb, "Authentication")

	server := app.NewServer(router, profileSvc, authenticationSvc)
	server.Run()
}
