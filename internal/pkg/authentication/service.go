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

type _AuthenticationService struct {
	db database.DatabaseService[Authentication]
}

func NewAuthenticationService() AuthenticationService {
	db := database.NewDatabaseService[Authentication]("Authentication")

	return &_AuthenticationService{
		db,
	}
}

type AuthenticationService interface {
	Login(request LoginRequest) (LoginResponse, *AuthenticationError)
	Signup(request SignupRequest) (SignupReponse, *AuthenticationError)
}

// Get dynamodb item with user credentials
func (service *_AuthenticationService) Login(request LoginRequest) (LoginResponse, *AuthenticationError) {
	input := map[string]string{
		"PK": request.Email,
	}

	result, err := service.db.Read(input)

	if err != nil {
		log.Println("LoginError: ", err)
		return LoginResponse{}, &AuthenticationError{Code: 503, Reason: "Internal Server Error"}
	}

	if result == nil {
		log.Println("LoginError: Item does not exist")
		return LoginResponse{}, &AuthenticationError{Code: 401, Reason: "Invalid email or password, please try again"}
	}

	invalid := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(request.Password))
	if invalid != nil {
		log.Println("LoginError: ", invalid)
		return LoginResponse{}, &AuthenticationError{Code: 401, Reason: "Invalid email or password, please try again"}
	}

	// create jwt token logic
	token, createTokenError := CreateToken(result.Id)
	if createTokenError != nil {
		log.Println("LoginError: error occurred when generating access token", createTokenError)
		return LoginResponse{}, &AuthenticationError{Code: 500, Reason: "Internal Server Error"}
	}

	return LoginResponse{
		AccessToken: token,
	}, nil
}

func (service *_AuthenticationService) Signup(request SignupRequest) (SignupReponse, *AuthenticationError) {
	// get dynamodb item with user credentials
	input := map[string]string{
		"PK": request.Email,
	}

	result, err := service.db.Read(input)
	if err != nil {
		fmt.Println("SignupError:", err)
		return SignupReponse{}, &AuthenticationError{Code: 503, Reason: "Internal Server Error"}
	}

	if result != nil {
		fmt.Println("SignupError: Account already exists")
		return SignupReponse{}, &AuthenticationError{Code: 400, Reason: "Account already exists"}
	}

	// Generate hash password and store
	password := []byte(request.Password)
	hashed, hashErr := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if hashErr != nil {
		fmt.Println("SignupError: Unable to generate hash from password", hashErr)
		return SignupReponse{}, &AuthenticationError{Code: 500, Reason: "Internal Server Error"}
	}

	// Create account if no issues exists
	account := Authentication{
		PK:       request.Email,
		Id:       uuid.New().String(),
		Email:    request.Email,
		Password: string(hashed),
		Name:     request.Name,
		Phone:    request.Phone,
	}

	// Save item to database
	err = service.db.Write(account)
	if err != nil {
		log.Printf("SignupError: %s", err)
		return SignupReponse{}, &AuthenticationError{Code: 503, Reason: "Internal Server Error"}
	}

	return SignupReponse{Status: true}, nil
}

func CreateToken(userId string) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userId
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
