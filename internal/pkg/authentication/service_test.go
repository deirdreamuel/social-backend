package authentication

import (
	"errors"
	"speakeasy/internal/pkg/database"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// Mock DatabaseService where item exists
type _DatabaseServiceMock_ItemExists struct {
	database.DatabaseService[Authentication]
}

func (db *_DatabaseServiceMock_ItemExists) Read(keyObj interface{}) (*Authentication, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("correct.password"), bcrypt.DefaultCost)
	return &Authentication{Password: string(hash)}, nil
}

func (db *_DatabaseServiceMock_ItemExists) Write(obj Authentication) error {
	return nil
}

// Mock DatabaseService where .Read returns an error
type _DatabaseServiceMock_ReadError struct {
	database.DatabaseService[Authentication]
}

func (db *_DatabaseServiceMock_ReadError) Read(keyObj interface{}) (*Authentication, error) {
	return nil, errors.New("ERROR")
}

func (db *_DatabaseServiceMock_ReadError) Write(obj Authentication) error {
	return nil
}

// Mock DatabaseService where .Write returns an error
type _DatabaseServiceMock_WriteError struct {
	database.DatabaseService[Authentication]
}

func (db *_DatabaseServiceMock_WriteError) Read(keyObj interface{}) (*Authentication, error) {
	return nil, nil
}

func (db *_DatabaseServiceMock_WriteError) Write(obj Authentication) error {
	return errors.New("ERROR")
}

// Mock DatabaseService where item does not exist
type _DatabaseServiceMock_ItemNotFound struct {
	database.DatabaseService[Authentication]
}

func (db *_DatabaseServiceMock_ItemNotFound) Read(keyObj interface{}) (*Authentication, error) {
	return nil, nil
}

func (db *_DatabaseServiceMock_ItemNotFound) Write(obj Authentication) error {
	return nil
}

func TestNewAuthenticationService(t *testing.T) {
	t.Run("SUCCESS: RETURN NEW AUTHENTICATION SERVICE", func(t *testing.T) {
		svc := NewAuthenticationService()

		assert.NotEmpty(t, svc, "Service should not empty")
	})
}

func TestLogin(t *testing.T) {
	t.Run("SUCCESS: RETURN JWT WHEN USER PASSWORD IS CORRECT", func(t *testing.T) {
		svc := &_AuthenticationService{db: &_DatabaseServiceMock_ItemExists{}}

		result, err := svc.Login(LoginRequest{
			Email:    "user@email.com",
			Password: "correct.password",
		})

		assert.Empty(t, err, "Error should be empty")
		assert.NotEmpty(t, result.AccessToken, "Access token should be returned")
	})

	t.Run("ERROR: RETURN 401 WHEN USER PASSWORD IS INCORRECT", func(t *testing.T) {
		svc := &_AuthenticationService{db: &_DatabaseServiceMock_ItemExists{}}

		result, err := svc.Login(LoginRequest{
			Email:    "user@email.com",
			Password: "wrong.password",
		})

		assert.Equal(t, 401, err.Code, "Error should be 401")
		assert.Empty(t, result.AccessToken, "Access token should be empty")
	})

	t.Run("ERROR: RETURN 401 WHEN USER DOES NOT EXIST", func(t *testing.T) {
		svc := &_AuthenticationService{db: &_DatabaseServiceMock_ItemNotFound{}}

		result, err := svc.Login(LoginRequest{
			Email:    "user@email.com",
			Password: "correct.password",
		})

		assert.Equal(t, 401, err.Code, "Error should be 401")
		assert.Empty(t, result.AccessToken, "Access token should be empty")
	})

	t.Run("ERROR: RETURN 503 WHEN DB RETURNS ERROR", func(t *testing.T) {
		svc := &_AuthenticationService{db: &_DatabaseServiceMock_ReadError{}}

		result, err := svc.Login(LoginRequest{
			Email:    "user@email.com",
			Password: "correct.password",
		})

		assert.Equal(t, 503, err.Code, "Error should be 503")
		assert.Empty(t, result, "Result should be empty")
	})
}

func TestSignup(t *testing.T) {
	t.Run("SUCCESS: RETURN CREATED MESSAGE WHEN ITEM NOT FOUND", func(t *testing.T) {
		svc := &_AuthenticationService{db: &_DatabaseServiceMock_ItemNotFound{}}
		result, err := svc.Signup(SignupRequest{
			Email:    "user@email.com",
			Password: "correct.password",
			Name:     "user.name",
			Phone:    "user.phone",
		})

		assert.Equal(t, true, result.Status, "Status should be true")
		assert.Empty(t, err, "Error should be empty")
	})

	t.Run("ERROR: RETURN 400 ERROR WHEN ITEM EXISTS", func(t *testing.T) {
		svc := &_AuthenticationService{db: &_DatabaseServiceMock_ItemExists{}}
		result, err := svc.Signup(SignupRequest{
			Email:    "user@email.com",
			Password: "correct.password",
			Name:     "user.name",
			Phone:    "user.phone",
		})

		assert.Equal(t, 400, err.Code, "Error should be 400")
		assert.Empty(t, result, "Result should be empty")
	})

	t.Run("ERROR: RETURN 503 WHEN DB READ RETURNS ERROR", func(t *testing.T) {
		svc := &_AuthenticationService{db: &_DatabaseServiceMock_ReadError{}}

		result, err := svc.Signup(SignupRequest{
			Email:    "user@email.com",
			Password: "correct.password",
			Name:     "user.name",
			Phone:    "user.phone",
		})

		assert.Equal(t, 503, err.Code, "Error should be 503")
		assert.Empty(t, result, "Result should be empty")
	})

	t.Run("ERROR: RETURN 503 WHEN DB WRITE RETURNS ERROR", func(t *testing.T) {
		svc := &_AuthenticationService{db: &_DatabaseServiceMock_WriteError{}}

		result, err := svc.Signup(SignupRequest{
			Email:    "user@email.com",
			Password: "correct.password",
			Name:     "user.name",
			Phone:    "user.phone",
		})

		assert.Equal(t, 503, err.Code, "Error should be 503")
		assert.Empty(t, result, "Result should be empty")
	})
}
