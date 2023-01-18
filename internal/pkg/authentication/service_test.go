package authentication

import (
	"errors"
	"speakeasy/pkg/database"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// Mock DatabaseService where item exists
type _DatabaseServiceMockItemExists struct {
	database.Service[Authentication]
}

func (db *_DatabaseServiceMockItemExists) Get(keyObj interface{}) (*Authentication, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("correct.password"), bcrypt.DefaultCost)
	return &Authentication{Password: string(hash)}, nil
}

func (db *_DatabaseServiceMockItemExists) Write(obj ...*Authentication) error {
	return nil
}

// Mock DatabaseService where .Get returns an error
type _DatabaseServiceMockGetError struct {
	database.Service[Authentication]
}

func (db *_DatabaseServiceMockGetError) Get(keyObj interface{}) (*Authentication, error) {
	return nil, errors.New("ERROR")
}

func (db *_DatabaseServiceMockGetError) Write(obj ...*Authentication) error {
	return nil
}

// Mock DatabaseService where .Write returns an error
type _DatabaseServiceMockWriteError struct {
	database.Service[Authentication]
}

func (db *_DatabaseServiceMockWriteError) Get(keyObj interface{}) (*Authentication, error) {
	return nil, nil
}

func (db *_DatabaseServiceMockWriteError) Write(obj ...*Authentication) error {
	return errors.New("ERROR")
}

// Mock DatabaseService where item does not exist
type _DatabaseServiceMockItemNotFound struct {
	database.Service[Authentication]
}

func (db *_DatabaseServiceMockItemNotFound) Get(keyObj interface{}) (*Authentication, error) {
	return nil, nil
}

func (db *_DatabaseServiceMockItemNotFound) Write(obj ...*Authentication) error {
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
		svc := &_Service{db: &_DatabaseServiceMockItemExists{}}

		result, err := svc.Login(LoginRequest{
			Email:    "user@email.com",
			Password: "correct.password",
		})

		assert.Empty(t, err, "Error should be empty")
		assert.NotEmpty(t, result.AccessToken, "Access token should be returned")
	})

	t.Run("ERROR: RETURN 401 WHEN USER PASSWORD IS INCORRECT", func(t *testing.T) {
		svc := &_Service{db: &_DatabaseServiceMockItemExists{}}

		result, err := svc.Login(LoginRequest{
			Email:    "user@email.com",
			Password: "wrong.password",
		})

		assert.Equal(t, 401, err.Code, "Error should be 401")
		assert.Empty(t, result, "Access token should be empty")
	})

	t.Run("ERROR: RETURN 401 WHEN USER DOES NOT EXIST", func(t *testing.T) {
		svc := &_Service{db: &_DatabaseServiceMockItemNotFound{}}

		result, err := svc.Login(LoginRequest{
			Email:    "user@email.com",
			Password: "correct.password",
		})

		assert.Equal(t, 401, err.Code, "Error should be 401")
		assert.Empty(t, result, "Access token should be empty")
	})

	t.Run("ERROR: RETURN 503 WHEN DB RETURNS ERROR", func(t *testing.T) {
		svc := &_Service{db: &_DatabaseServiceMockGetError{}}

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
		svc := &_Service{db: &_DatabaseServiceMockItemNotFound{}}
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
		svc := &_Service{db: &_DatabaseServiceMockItemExists{}}
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
		svc := &_Service{db: &_DatabaseServiceMockGetError{}}

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
		svc := &_Service{db: &_DatabaseServiceMockWriteError{}}

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
