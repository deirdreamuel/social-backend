package trip

import (
	"fmt"
	"log"
	"speakeasy/pkg"
	"speakeasy/pkg/database"
	"strings"
	"time"

	"github.com/google/uuid"
)

type _Service struct {
	db database.Service[Trip]
}

// NewTripService returns _Service object
func NewTripService() Service {
	db := database.NewDatabaseService[Trip]("APPLICATION")

	return &_Service{
		db,
	}
}

// Service interface which contains trip operations
type Service interface {
	CreateTrip(trip *Trip) *pkg.Error
	GetTrip(tripID string) (*Trip, *pkg.Error)
	GetTripsByUser(userID string) (*[]Trip, *pkg.Error)
	GetTripParticipants(tripID string) (*[]Trip, *pkg.Error)
}

// CreateTrip function to create trip
func (service *_Service) CreateTrip(trip *Trip) *pkg.Error {
	// Validate create trip request
	if len(strings.TrimSpace(trip.CreatedBy)) == 0 {
		return &pkg.Error{Code: 400, Reason: "User ID cannot be empty"}
	}

	if len(strings.TrimSpace(trip.Name)) == 0 {
		return &pkg.Error{Code: 400, Reason: "Trip name cannot be empty"}
	}

	// Validate trip dates
	from, err := time.Parse(time.RFC3339, trip.FromDate)
	if err != nil {
		log.Println("CreateTripError:", err)
		return &pkg.Error{Code: 400, Reason: "Trip dates are invalid"}
	}

	to, err := time.Parse(time.RFC3339, trip.ToDate)
	if err != nil {
		log.Println("CreateTripError:", err)
		return &pkg.Error{Code: 400, Reason: "Trip dates are invalid"}
	}

	if from.Before(time.Now()) || to.Before(time.Now()) {
		log.Println("CreateTripError: dates are in the past")
		return &pkg.Error{Code: 400, Reason: "Trip dates are invalid"}
	}

	if from.After(to) {
		log.Println("CreateTripError: from_date if after to_date ")
		return &pkg.Error{Code: 400, Reason: "Trip dates are invalid"}
	}

	// Add primary and sort key to item
	uid := uuid.New().String()
	trip.ID = uid
	trip.PK = fmt.Sprintf("TRIP#%s", uid)
	trip.SK = fmt.Sprintf("TRIP#%s", uid)

	// Create another item for user reference
	userTrip := *trip
	userTrip.PK = fmt.Sprintf("USER#%s", userTrip.CreatedBy)

	err = service.db.Write(*trip, userTrip)
	if err != nil {
		log.Println("CreateTripError:", err)
		return &pkg.Error{Code: 503, Reason: "Internal Server Error"}
	}

	return nil
}

// GetTrip function to get trip by id
func (service *_Service) GetTrip(tripID string) (*Trip, *pkg.Error) {
	input := map[string]string{
		"PK": fmt.Sprintf("TRIP#%s", tripID),
		"SK": fmt.Sprintf("TRIP#%s", tripID),
	}

	result, err := service.db.Get(input)
	if err != nil {
		log.Println("GetTrip:", err)
		return nil, &pkg.Error{Code: 503, Reason: "Internal Server Error"}
	}

	if result == nil {
		log.Println("GetTrip: item not found")
		return nil, &pkg.Error{Code: 400, Reason: "Trip not found"}
	}

	return result, nil
}

func (service *_Service) GetTripsByUser(userID string) (*[]Trip, *pkg.Error) {
	filter := map[string]string{
		":PK": fmt.Sprintf("USER#%s", userID),
		":SK": "TRIP",
	}

	condition := "PK = :PK And begins_with(SK, :SK)"

	results, err := service.db.Query(filter, condition)
	if err != nil {
		log.Println("GetTripsByUser:", err)
		return nil, &pkg.Error{Code: 503, Reason: "Internal Server Error"}
	}

	return results, nil
}

func (service *_Service) GetTripParticipants(tripID string) (*[]Trip, *pkg.Error) {
	filter := map[string]string{
		":SK": fmt.Sprintf("TRIP#%s", tripID),
		":PK": "USER",
	}

	log.Println(filter)

	condition := "SK = :SK"
	filterExpr := "begins_with(PK, :PK)"

	results, err := service.db.QueryWithIndex(filter, condition, filterExpr, "APPLICATION_GSI_1")
	if err != nil {
		log.Println("GetTripParticipants:", err)
		return nil, &pkg.Error{Code: 503, Reason: "Internal Server Error"}
	}

	return results, nil
}
