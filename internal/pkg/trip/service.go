package trip

import (
	"fmt"
	"log"
	"speakeasy/pkg"
	"speakeasy/pkg/database"
	"strings"
	"time"
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
}

// CreateTrip function to create trip
func (service *_Service) CreateTrip(trip *Trip) *pkg.Error {
	if len(strings.TrimSpace(trip.UserID)) == 0 {
		return &pkg.Error{Code: 400, Reason: "User ID cannot be empty"}
	}

	if len(strings.TrimSpace(trip.Name)) == 0 {
		return &pkg.Error{Code: 400, Reason: "Trip name cannot be empty"}
	}

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

	writeError := service.db.Write(*trip)
	if writeError != nil {
		log.Println("CreateTripError:", writeError)
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

	result, err := service.db.Read(input)
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
