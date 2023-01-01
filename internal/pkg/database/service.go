package database

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type _DatabaseService[T any] struct {
	db        *dynamodb.DynamoDB
	tableName string
}

func NewDatabaseService[T any](tableName string) DatabaseService[T] {
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
	return &_DatabaseService[T]{db, tableName}
}

type DatabaseService[T any] interface {
	Read(keyObj interface{}) (*T, error)
	Write(obj T) error
	Delete(obj interface{}) error
}

func (service *_DatabaseService[T]) Read(keyObj interface{}) (*T, error) {
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

func (service *_DatabaseService[T]) Write(obj T) error {
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

func (service *_DatabaseService[T]) Delete(keyObj interface{}) error {
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
