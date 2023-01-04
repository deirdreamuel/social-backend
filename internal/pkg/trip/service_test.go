package trip

import (
	"errors"
	"speakeasy/pkg/database"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Mock DatabaseService where item exists
type _DatabaseServiceMockItemExists struct {
	database.Service[Trip]
}

func (db *_DatabaseServiceMockItemExists) Read(keyObj interface{}) (*Trip, error) {
	return &Trip{
		UserID:      "0000-0000-0000-0000",
		FromDate:    time.Now().Add(time.Hour * 24).UTC().Format(time.RFC3339),
		ToDate:      time.Now().Add(time.Hour * 24 * 2).UTC().Format(time.RFC3339),
		Name:        "trip.name",
		Description: "trip.description",
		Location: Location{
			City:    "New York",
			State:   "NY",
			Country: "US",
		},
		Participants: []Participant{
			{
				Email: "user.email@example.com",
			},
		},
	}, nil
}

func (db *_DatabaseServiceMockItemExists) Write(obj Trip) error {
	return nil
}

// Mock DatabaseService where item does not exist
type _DatabaseServiceMockItemNotFound struct {
	database.Service[Trip]
}

func (db *_DatabaseServiceMockItemNotFound) Read(keyObj interface{}) (*Trip, error) {
	return nil, nil
}

func (db *_DatabaseServiceMockItemNotFound) Write(obj Trip) error {
	return nil
}

// Mock DatabaseService where .Write returns an error
type _DatabaseServiceMockReadError struct {
	database.Service[Trip]
}

func (db *_DatabaseServiceMockReadError) Read(keyObj interface{}) (*Trip, error) {
	return nil, errors.New("ERROR")
}

func (db *_DatabaseServiceMockReadError) Write(obj Trip) error {
	return nil
}

// Mock DatabaseService where .Write returns an error
type _DatabaseServiceMockWriteError struct {
	database.Service[Trip]
}

func (db *_DatabaseServiceMockWriteError) Read(keyObj interface{}) (*Trip, error) {
	return nil, nil
}

func (db *_DatabaseServiceMockWriteError) Write(obj Trip) error {
	return errors.New("ERROR")
}

func TestNewTripService(t *testing.T) {
	t.Run("SUCCESS: RETURN NEW AUTHENTICATION SERVICE", func(t *testing.T) {
		svc := NewTripService()

		assert.NotEmpty(t, svc, "Service should not empty")
	})
}

func TestCreateTrip(t *testing.T) {
	t.Run("SUCCESS: RETURN CREATED MESSAGE WHEN FIELDS ARE VALID", func(t *testing.T) {
		svc := _Service{db: &_DatabaseServiceMockItemNotFound{}}

		err := svc.CreateTrip(&Trip{
			UserID:      "0000-0000-0000-0000",
			FromDate:    time.Now().Add(time.Hour * 24).UTC().Format(time.RFC3339),
			ToDate:      time.Now().Add(time.Hour * 24 * 2).UTC().Format(time.RFC3339),
			Name:        "trip.name",
			Description: "trip.description",
			Location: Location{
				City:    "New York",
				State:   "NY",
				Country: "US",
			},
			Participants: []Participant{
				{
					Email: "user.email@example.com",
				},
			},
		})

		assert.Empty(t, err)
	})

	t.Run("ERROR: RETURN 400 WHEN FROM_DATE IS INVALID", func(t *testing.T) {
		svc := _Service{db: &_DatabaseServiceMockItemNotFound{}}

		err := svc.CreateTrip(&Trip{
			UserID:      "0000-0000-0000-0000",
			FromDate:    "invalid.time",
			ToDate:      time.Now().Add(time.Hour * 24).UTC().Format(time.RFC3339),
			Name:        "trip.name",
			Description: "trip.description",
			Location: Location{
				City:    "New York",
				State:   "NY",
				Country: "US",
			},
			Participants: []Participant{
				{
					Email: "user.email@example.com",
				},
			},
		})

		assert.Equal(t, 400, err.Code, "Error should be 400")
	})

	t.Run("ERROR: RETURN 400 WHEN TO_DATE IS INVALID", func(t *testing.T) {
		svc := _Service{db: &_DatabaseServiceMockItemNotFound{}}

		err := svc.CreateTrip(&Trip{
			UserID:      "0000-0000-0000-0000",
			FromDate:    time.Now().Add(time.Hour * 24).UTC().Format(time.RFC3339),
			ToDate:      "invalid.time",
			Name:        "trip.name",
			Description: "trip.description",
			Location: Location{
				City:    "New York",
				State:   "NY",
				Country: "US",
			},
			Participants: []Participant{
				{
					Email: "user.email@example.com",
				},
			},
		})

		assert.Equal(t, 400, err.Code, "Error should be 400")
	})

	t.Run("ERROR: RETURN 400 WHEN DATES ARE IN THE PAST", func(t *testing.T) {
		svc := _Service{db: &_DatabaseServiceMockItemNotFound{}}

		err := svc.CreateTrip(&Trip{
			UserID:      "0000-0000-0000-0000",
			FromDate:    time.Now().Add(-time.Hour * 24 * 2).UTC().Format(time.RFC3339),
			ToDate:      time.Now().Add(-time.Hour * 24).UTC().Format(time.RFC3339),
			Name:        "trip.name",
			Description: "trip.description",
			Location: Location{
				City:    "New York",
				State:   "NY",
				Country: "US",
			},
			Participants: []Participant{
				{
					Email: "user.email@example.com",
				},
			},
		})

		assert.Equal(t, 400, err.Code, "Error should be 400")
	})

	t.Run("ERROR: RETURN 400 WHEN FROM_DATE IS AFTER TO_DATE", func(t *testing.T) {
		svc := _Service{db: &_DatabaseServiceMockItemNotFound{}}

		err := svc.CreateTrip(&Trip{
			UserID:      "0000-0000-0000-0000",
			FromDate:    time.Now().Add(time.Hour * 24 * 2).UTC().Format(time.RFC3339),
			ToDate:      time.Now().Add(time.Hour * 24).UTC().Format(time.RFC3339),
			Name:        "trip.name",
			Description: "trip.description",
			Location: Location{
				City:    "New York",
				State:   "NY",
				Country: "US",
			},
			Participants: []Participant{
				{
					Email: "user.email@example.com",
				},
			},
		})

		assert.Equal(t, 400, err.Code, "Error should be 400")
	})

	t.Run("ERROR: RETURN 400 WHEN USER ID IS EMPTY", func(t *testing.T) {
		svc := _Service{db: &_DatabaseServiceMockItemNotFound{}}

		err := svc.CreateTrip(&Trip{
			UserID:      "",
			FromDate:    time.Now().Add(-time.Hour * 24 * 2).UTC().Format(time.RFC3339),
			ToDate:      time.Now().Add(-time.Hour * 24).UTC().Format(time.RFC3339),
			Name:        "trip.name",
			Description: "trip.description",
			Location: Location{
				City:    "New York",
				State:   "NY",
				Country: "US",
			},
			Participants: []Participant{
				{
					Email: "user.email@example.com",
				},
			},
		})

		assert.Equal(t, 400, err.Code, "Error should be 400")
	})

	t.Run("ERROR: RETURN 400 WHEN NAME IS EMPTY", func(t *testing.T) {
		svc := _Service{db: &_DatabaseServiceMockItemNotFound{}}

		err := svc.CreateTrip(&Trip{
			UserID:      "0000-0000-0000-0000",
			FromDate:    time.Now().Add(-time.Hour * 24 * 2).UTC().Format(time.RFC3339),
			ToDate:      time.Now().Add(-time.Hour * 24).UTC().Format(time.RFC3339),
			Name:        "",
			Description: "trip.description",
			Location: Location{
				City:    "New York",
				State:   "NY",
				Country: "US",
			},
			Participants: []Participant{
				{
					Email: "user.email@example.com",
				},
			},
		})

		assert.Equal(t, 400, err.Code, "Error should be 400")
	})

	t.Run("ERROR: RETURN 503 WHEN DB WRITE RETURNS ERROR", func(t *testing.T) {
		svc := &_Service{db: &_DatabaseServiceMockWriteError{}}

		err := svc.CreateTrip(&Trip{
			UserID:      "0000-0000-0000-0000",
			FromDate:    time.Now().Add(time.Hour * 24).UTC().Format(time.RFC3339),
			ToDate:      time.Now().Add(time.Hour * 24 * 2).UTC().Format(time.RFC3339),
			Name:        "trip.name",
			Description: "trip.description",
			Location: Location{
				City:    "New York",
				State:   "NY",
				Country: "US",
			},
			Participants: []Participant{
				{
					Email: "user.email@example.com",
				},
			},
		})

		assert.Equal(t, 503, err.Code, "Error should be 503")
	})
}

func TestGetTrip(t *testing.T) {
	t.Run("SUCCESS: RETURN 200 WHEN ITEM IS FOUND", func(t *testing.T) {
		svc := _Service{db: &_DatabaseServiceMockItemExists{}}

		result, err := svc.GetTrip("0000-0000-0000-0000")

		assert.NotEmpty(t, result, "Result should be not be empty")
		assert.Empty(t, err, "Error should be empty")
	})

	t.Run("ERROR: RETURN 400 WHEN ITEM NOT FOUND", func(t *testing.T) {
		svc := _Service{db: &_DatabaseServiceMockItemNotFound{}}

		result, err := svc.GetTrip("0000-0000-0000-0000")

		assert.Empty(t, result, "Result should be empty")
		assert.Equal(t, 400, err.Code, "Error should be 400")
	})

	t.Run("ERROR: RETURN 503 WHEN DB WRITE RETURNS ERROR", func(t *testing.T) {
		svc := &_Service{db: &_DatabaseServiceMockReadError{}}

		result, err := svc.GetTrip("0000-0000-0000-0000")

		assert.Empty(t, result, "Result should be empty")
		assert.Equal(t, 503, err.Code, "Error should be 503")
	})
}
