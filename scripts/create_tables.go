package main

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/joho/godotenv"

	"fmt"
	"log"
)

func main() {
	godotenv.Load()

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

	// Create DynamoDB client
	svc := dynamodb.New(sess, cfg)

	if err := createAuthenticationTable(svc); err != nil {
		log.Fatalf("Got error calling CreateAuthenticationTable: %s", err)
	}

	if err := createApplicationTable(svc); err != nil {
		log.Fatalf("Got error calling createApplicationTable: %s", err)
	}
}

func createAuthenticationTable(ddb *dynamodb.DynamoDB) error {
	tableName := "AUTHENTICATION"

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("PK"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("PK"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
	}

	_, err := ddb.CreateTable(input)
	if err != nil {
		log.Printf("Got error calling CreateTable: %s", err)
		return err
	}

	fmt.Println("Created the table", tableName)
	return nil
}

func createApplicationTable(svc *dynamodb.DynamoDB) error {
	tableName := "APPLICATION"

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("PK"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("SK"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("PK"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("SK"),
				KeyType:       aws.String("RANGE"),
			},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("APPLICATION_GSI_1"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("SK"),
						KeyType:       aws.String("HASH"),
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(10),
					WriteCapacityUnits: aws.Int64(10),
				},
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
	}

	_, err := svc.CreateTable(input)
	if err != nil {
		log.Printf("Got error calling CreateTable: %s", err)
		return err
	}

	fmt.Println("Created the table", tableName)
	return nil
}
