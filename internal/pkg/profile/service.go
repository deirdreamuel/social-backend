package profile

import (
	"errors"
	"log"
	"speakeasy/internal/pkg/database"
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	PK        string
	SK        string
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	BirthDate int       `json:"birth_date"`
	Gender    string    `json:"gender"`
	Email     string    `json:"email"`
}

type _ProfileService struct {
	db        database.DatabaseService[Profile]
	tableName string
}

type ProfileService interface {
	PutProfile(profile Profile) error
	GetProfile(id string) (Profile, error)
	DeleteProfile(id string) error
}

func NewProfileService(tableName string) ProfileService {
	db := database.NewDatabaseService[Profile]("Profile")
	return &_ProfileService{
		db:        db,
		tableName: tableName,
	}
}

// PutProfile puts Profile object to database
// Takes in Profile object as argument
// It returns any error encountered
func (service *_ProfileService) PutProfile(profile Profile) error {
	profile.PK = "PROFILE"
	profile.SK = uuid.New().String()

	err := service.db.Write(profile)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}

	return nil
}

func (service *_ProfileService) GetProfile(id string) (Profile, error) {
	// Query profile from database
	obj := Profile{
		PK: "PROFILE",
		SK: id,
	}

	result, err := service.db.Read(obj)

	if err != nil {
		log.Printf("error deleting dynamodb item: %s", err)
		return Profile{}, err
	}

	if result == nil {
		return Profile{}, errors.New("item not found")
	}

	return obj, nil
}

func (service *_ProfileService) DeleteProfile(id string) error {
	obj := Profile{
		PK: "PROFILE",
		SK: id,
	}

	err := service.db.Delete(obj)
	if err != nil {
		log.Fatalf("error deleting dynamodb item: %s", err)
		return err
	}

	return nil
}
