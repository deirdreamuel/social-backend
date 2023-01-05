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
	Get(keyObj interface{}) (*T, error)
	Write(obj ...T) error
	Delete(obj interface{}) error
	Query(filterObj interface{}, condition string) (*[]T, error)
	QueryWithIndex(filterObj interface{}, condition string, filterExpr string, index string) (*[]T, error)
}

// Get function to read data from database
func (service *_Service[T]) Get(keyObj interface{}) (*T, error) {
	// Create key object for DynamoDB Key
	key, marshallError := dynamodbattribute.MarshalMap(keyObj)
	if marshallError != nil {
		log.Println("GetError: MarshalError: ", marshallError)
		return nil, marshallError
	}

	input := &dynamodb.GetItemInput{
		TableName: &service.tableName,
		Key:       key,
	}

	result, err := service.db.GetItem(input)
	if err != nil {
		log.Println("GetError: ", err)
		return nil, err
	}

	if result.Item == nil {
		return nil, nil
	}

	var out T
	err = dynamodbattribute.UnmarshalMap(result.Item, &out)

	return &out, err
}

// Write function to write data from database
func (service *_Service[T]) Write(objs ...T) error {
	items := []*dynamodb.WriteRequest{}

	for _, obj := range objs {
		// Create item object for DynamoDB
		item, err := dynamodbattribute.MarshalMap(obj)
		if err != nil {
			log.Println("Error: ", err)
			return err
		}

		req := dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: item,
			},
		}

		items = append(items, &req)
	}

	input := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			service.tableName: items,
		},
	}

	_, err := service.db.BatchWriteItem(input)
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
		TableName: &service.tableName,
		Key:       key,
	}

	_, err := service.db.DeleteItem(input)
	return err
}

// Query function to query data from database
func (service *_Service[T]) Query(filterObj interface{}, condition string) (*[]T, error) {
	// Create item object for DynamoDB
	filter, err := dynamodbattribute.MarshalMap(filterObj)
	if err != nil {
		log.Println("QueryError: ", err)
		return nil, err
	}

	input := &dynamodb.QueryInput{
		TableName:                 &service.tableName,
		KeyConditionExpression:    &condition,
		ExpressionAttributeValues: filter,
	}

	result, err := service.db.Query(input)
	if err != nil {
		log.Println("QueryError: ", err)
		return nil, err
	}

	var out []T

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &out)

	return &out, err
}

// Query function to query data from database
func (service *_Service[T]) QueryWithIndex(filterObj interface{}, condition string, filterExpr string, index string) (*[]T, error) {
	// Create item object for DynamoDB
	filter, err := dynamodbattribute.MarshalMap(filterObj)
	if err != nil {
		log.Println("QueryWithIndexError: ", err)
		return nil, err
	}

	input := &dynamodb.QueryInput{
		TableName:                 &service.tableName,
		IndexName:                 &index,
		KeyConditionExpression:    &condition,
		FilterExpression:          &filterExpr,
		ExpressionAttributeValues: filter,
	}

	result, err := service.db.Query(input)
	if err != nil {
		log.Println("QueryWithIndexError: ", err)
		return nil, err
	}

	var out []T

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &out)

	return &out, err
}
