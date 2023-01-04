package database

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type _Service[T any] struct {
	db        *dynamodb.DynamoDB
	tableName string
}

// NewDatabaseService function to initialize Service object
func NewDatabaseService[T any](tableName string) Service[T] {
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
	db := dynamodb.New(sess, cfg)
	return &_Service[T]{db, tableName}
}

// Service interface which contains database operations
type Service[T any] interface {
	Read(keyObj interface{}) (*T, error)
	Write(obj T) error
	Delete(obj interface{}) error
}

// Read function to read data from database
func (service *_Service[T]) Read(keyObj interface{}) (*T, error) {
	// Create key object for DynamoDB Key
	key, marshallError := dynamodbattribute.MarshalMap(keyObj)
	if marshallError != nil {
		log.Println("ReadError: MarshalError: ", marshallError)
		return nil, marshallError
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(service.tableName),
		Key:       key,
	}

	result, err := service.db.GetItem(input)
	if err != nil {
		log.Println("ReadError: ", err)
		return nil, err
	}

	if result.Item == nil {
		return nil, nil
	}

	var out T
	dynamodbattribute.UnmarshalMap(result.Item, &out)

	return &out, nil
}

// Write function to write data from database
func (service *_Service[T]) Write(obj T) error {
	// Create item object for DynamoDB
	item, marshallError := dynamodbattribute.MarshalMap(obj)
	if marshallError != nil {
		log.Println("Error: ", marshallError)
		return marshallError
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(service.tableName),
		Item:      item,
	}

	_, err := service.db.PutItem(input)
	return err
}

// Delete function to delete data from database
func (service *_Service[T]) Delete(keyObj interface{}) error {
	// Create item object for DynamoDB
	key, marshalError := dynamodbattribute.MarshalMap(keyObj)
	if marshalError != nil {
		log.Println("Error: ", marshalError)
		return marshalError
	}

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(service.tableName),
		Key:       key,
	}

	_, err := service.db.DeleteItem(input)
	return err
}
