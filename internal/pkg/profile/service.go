package profile

import (
	"errors"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

type Profile struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	BirthDate int       `json:"birth_date"`
	Gender    string    `json:"gender"`
	Email     string    `json:"email"`
}

type _ProfileService struct {
	ddb       *dynamodb.DynamoDB
	tableName string
}

type ProfileService interface {
	PutProfile(profile Profile) error
	GetProfile(id string) (Profile, error)
	DeleteProfile(id string) error
}

func NewProfileService(ddb *dynamodb.DynamoDB, tableName string) ProfileService {
	return &_ProfileService{
		ddb:       ddb,
		tableName: tableName,
	}
}

// PutProfile puts Profile object to database
// Takes in Profile object as argument
// It returns any error encountered
func (service *_ProfileService) PutProfile(profile Profile) error {
	item, err := dynamodbattribute.MarshalMap(profile)
	if err != nil {
		log.Printf("error mapping profile object: %s", err)
	}

	pk := "PROFILE"
	id := uuid.New().String()
	item["PK"] = &dynamodb.AttributeValue{S: &pk}
	item["SK"] = &dynamodb.AttributeValue{S: &id}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(service.tableName),
	}

	_, err = service.ddb.PutItem(input)
	if err != nil {
		log.Printf("error putting dynamodb item: %s", err)
		return err
	}

	return nil
}

func (service *_ProfileService) GetProfile(id string) (Profile, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(service.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String("PROFILE"),
			},
			"SK": {
				S: aws.String(id),
			},
		},
	}

	result, _ := service.ddb.GetItem(input)

	if result.Item == nil {
		msg := "could not find item"
		return Profile{}, errors.New(msg)
	}

	item := Profile{}
	err := dynamodbattribute.UnmarshalMap(result.Item, &item)

	return Profile{}, err
}

func (service *_ProfileService) DeleteProfile(id string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(service.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String("PROFILE"),
			},
			"SK": {
				S: aws.String(id),
			},
		},
	}

	_, err := service.ddb.DeleteItem(input)
	if err != nil {
		log.Fatalf("error deleting dynamodb item: %s", err)
		return err
	}

	return nil
}
