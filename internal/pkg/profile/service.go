package profile

import (
	"fmt"
	"log"
	"mime/multipart"
	"speakeasy/pkg"
	"speakeasy/pkg/database"
	"speakeasy/pkg/filestorage"
	"time"
)

var PROFILE_PK string = "USER#%s"
var PROFILE_SK string = "__PROFILE__"

// Profile object to store in database.
// PK (Primary Key) should be in the format of PROFILE_PK value,
// SK (Sort Key) should be PROFILE_SK value.
type Profile struct {
	UpdatedAt     time.Time `json:"updated_at"`
	PK            string    `json:"PK,omitempty"`
	SK            string    `json:"SK,omitempty"`
	UserID        string    `json:"user_id"`
	Name          string    `json:"name"`
	Bio           string    `json:"bio"`
	ProfilePicUrl string    `json:"profile_pic_url,omitempty"`
}

type _Service struct {
	db      database.Service[Profile]
	storage filestorage.Service
}

type Service interface {
	PutProfile(profile *Profile) *pkg.Error
	GetProfile(id string) (*Profile, *pkg.Error)
	UploadProfilePicture(userID string, file multipart.File) error
}

// NewProfileService initializes database and returns Service object
func NewProfileService() Service {
	db := database.NewDatabaseService[Profile]("APPLICATION")
	storage := filestorage.NewFileStorageService("profile.image.amuel.org")

	return &_Service{db, storage}
}

// PutProfile function to configure db keys and update information.
func (service *_Service) PutProfile(profile *Profile) *pkg.Error {
	profile.PK = fmt.Sprintf(PROFILE_PK, profile.UserID)
	profile.SK = PROFILE_SK

	profile.UpdatedAt = time.Now().UTC()

	err := service.db.Write(profile)
	if err != nil {
		log.Printf("(PutProfile) error: %s", err)
		return &pkg.Error{Code: 503, Reason: "Internal Server Error"}
	}

	return nil
}

// GetProfile function to configure db keys and update information.
func (service *_Service) GetProfile(userID string) (*Profile, *pkg.Error) {
	input := map[string]string{
		"PK": fmt.Sprintf(PROFILE_PK, userID),
		"SK": PROFILE_SK,
	}

	profile, err := service.db.Get(input)
	if err != nil {
		log.Println("(GetProfile) error:", err)
		return nil, &pkg.Error{Code: 503, Reason: "Internal Server Error"}
	}

	if profile == nil {
		log.Println("(GetProfile) error: item not found")
		return &Profile{UserID: userID}, nil
	}

	removePrivateFieldsFromJSON(profile)

	// TODO: add profile pic url based on environment

	return profile, nil
}

// UploadProfilePicture function to upload file with userID as name.
// Profile picture can be seen by anyone, it's ok to use userID as name
func (service *_Service) UploadProfilePicture(userID string, file multipart.File) error {
	return service.storage.UploadFile(userID, &file)
}

// removePrivateFieldsFromJSON sets keys to empty and removes from being serialized
func removePrivateFieldsFromJSON(profile *Profile) {
	profile.PK = ""
	profile.SK = ""
}
