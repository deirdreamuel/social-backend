package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"fmt"
	"log"
)

func main() {
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	creds := credentials.NewStaticCredentials("123", "123", "")
	region := "us-east-1"
	endpoint := "http://127.0.0.1:8000"
	awsConfig := &aws.Config{
		Credentials: creds,
		Region:      &region,
		Endpoint:    &endpoint,
	}

	// Create DynamoDB client
	svc := dynamodb.New(sess, awsConfig)

	if err := CreateAuthenticationTable(svc); err != nil {
		log.Fatalf("Got error calling CreateAuthenticationTable: %s", err)
	}

	if err := CreateProfileTable(svc); err != nil {
		log.Fatalf("Got error calling CreateProfileTable: %s", err)
	}
}

func CreateAuthenticationTable(ddb *dynamodb.DynamoDB) error {
	tableName := "Authentication"

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

func CreateProfileTable(svc *dynamodb.DynamoDB) error {
	tableName := "Profile"

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
