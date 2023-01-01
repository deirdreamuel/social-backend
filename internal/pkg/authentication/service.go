package authentication

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"speakeasy/internal/pkg/database"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type _Service struct {
	db database.Service[Authentication]
}

// NewAuthenticationService returns _AuthenticationService object
func NewAuthenticationService() Service {
	db := database.NewDatabaseService[Authentication]("Authentication")

	return &_Service{
		db,
	}
}

// Service interface which contains authentication operations
type Service interface {
	Login(request LoginRequest) (*LoginResponse, *Error)
	Signup(request SignupRequest) (*SignupReponse, *Error)
}

// Login function to get access token
func (service *_Service) Login(request LoginRequest) (*LoginResponse, *Error) {
	input := map[string]string{
		"PK": request.Email,
	}

	result, err := service.db.Read(input)

	if err != nil {
		log.Println("LoginError: ", err)
		return nil, &Error{Code: 503, Reason: "Internal Server Error"}
	}

	if result == nil {
		log.Println("LoginError: Item does not exist")
		return nil, &Error{Code: 401, Reason: "Invalid email or password, please try again"}
	}

	invalid := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(request.Password))
	if invalid != nil {
		log.Println("LoginError: ", invalid)
		return nil, &Error{Code: 401, Reason: "Invalid email or password, please try again"}
	}

	// create jwt token logic
	token, createTokenError := CreateToken(result.ID)
	if createTokenError != nil {
		log.Println("LoginError: error occurred when generating access token", createTokenError)
		return nil, &Error{Code: 500, Reason: "Internal Server Error"}
	}

	return &LoginResponse{
		AccessToken: token,
	}, nil
}

// Signup function to create an account
func (service *_Service) Signup(request SignupRequest) (*SignupReponse, *Error) {
	// get dynamodb item with user credentials
	input := map[string]string{
		"PK": request.Email,
	}

	result, err := service.db.Read(input)
	if err != nil {
		log.Println("SignupError:", err)
		return nil, &Error{Code: 503, Reason: "Internal Server Error"}
	}

	if result != nil {
		log.Println("SignupError: Account already exists")
		return nil, &Error{Code: 400, Reason: "Account already exists"}
	}

	// Generate hash password and store
	password := []byte(request.Password)
	hashed, hashErr := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if hashErr != nil {
		log.Println("SignupError: Unable to generate hash from password", hashErr)
		return nil, &Error{Code: 500, Reason: "Internal Server Error"}
	}

	// Create account if no issues exists
	account := Authentication{
		PK:       request.Email,
		ID:       uuid.New().String(),
		Email:    request.Email,
		Password: string(hashed),
		Name:     request.Name,
		Phone:    request.Phone,
	}

	// Save item to database
	err = service.db.Write(account)
	if err != nil {
		log.Printf("SignupError: %s", err)
		return nil, &Error{Code: 503, Reason: "Internal Server Error"}
	}

	return &SignupReponse{Status: true}, nil
}

// CreateToken function to create jwt
func CreateToken(userID string) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return token, nil
}

// ExtractToken function to extract jwt from http request
func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

// VerifyToken function to verify jwt from http request
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
