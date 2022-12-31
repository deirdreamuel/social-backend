package authentication

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type _AuthenticationService struct {
	ddb       *dynamodb.DynamoDB
	tableName string
}

func NewAuthenticationService(ddb *dynamodb.DynamoDB, tableName string) AuthenticationService {
	return &_AuthenticationService{
		ddb,
		tableName,
	}
}

type AuthenticationService interface {
	Login(request LoginRequest) (LoginResponse, error)
	Signup(request SignupRequest) (SignupReponse, error)
}

func (service *_AuthenticationService) Login(request LoginRequest) (LoginResponse, error) {
	// get dynamodb item with user credentials
	input := &dynamodb.GetItemInput{
		TableName: aws.String(service.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String(request.Email),
			},
		},
	}

	result, err := service.ddb.GetItem(input)

	if err != nil {
		fmt.Println("error occurred when getting item", err)
		return LoginResponse{}, err
	}

	if result.Item == nil {
		fmt.Println("error item does not exist")
		return LoginResponse{}, errors.New("item does not exists")
	}

	item := Authentication{}
	dynamodbattribute.UnmarshalMap(result.Item, &item)

	invalidPasswordError := bcrypt.CompareHashAndPassword([]byte(item.Password), []byte(request.Password))
	if invalidPasswordError != nil {
		fmt.Println("error validating password")
		return LoginResponse{}, invalidPasswordError
	}

	// create jwt token logic
	token, createTokenError := CreateToken(item.Id)
	if createTokenError != nil {
		fmt.Println("Login: error occurred when generating access token", createTokenError)
		return LoginResponse{}, err
	}

	response := LoginResponse{
		AccessToken: token,
	}

	return response, nil
}

func (service *_AuthenticationService) Signup(request SignupRequest) (SignupReponse, error) {
	// get dynamodb item with user credentials
	input := &dynamodb.GetItemInput{
		TableName: aws.String(service.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String(request.Email),
			},
		},
	}

	result, err := service.ddb.GetItem(input)
	if err != nil {
		fmt.Println("error occurred when getting item", err)
		return SignupReponse{}, err
	}

	if result.Item != nil {
		fmt.Println("error account already exists")
		return SignupReponse{}, errors.New("account already exists")
	}

	// hash password for storage
	password := []byte(request.Password)
	hashed, hashErr := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if hashErr != nil {
		fmt.Println("error occurred when getting item", hashErr)
		return SignupReponse{}, err
	}

	// create account if no issues exists
	account := Authentication{
		PK:       request.Email,
		Id:       uuid.New().String(),
		Email:    request.Email,
		Password: string(hashed),
		Name:     request.Name,
		Phone:    request.Phone,
	}

	// save item to database
	accountItem, marshallErr := dynamodbattribute.MarshalMap(account)
	putItemInput := &dynamodb.PutItemInput{
		Item:      accountItem,
		TableName: aws.String(service.tableName),
	}

	if marshallErr != nil {
		fmt.Println("error occurred when getting item", hashErr)
		return SignupReponse{}, err
	}

	_, err = service.ddb.PutItem(putItemInput)
	if err != nil {
		log.Printf("error putting dynamodb item: %s", err)
		return SignupReponse{}, err
	}

	return SignupReponse{Status: true}, nil
}

func CreateToken(userId string) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["id"] = userId
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return token, nil
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
