package attachment

import (
	"time"

	"speakeasy/internal/pkg/profile"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Interaction struct {
	Id        string          `json:"id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	From      profile.Profile `json:"from"`
	To        profile.Profile `json:"to"`
}

type interactionService struct {
	ddb       *dynamodb.DynamoDB
	tableName string
}

type InteractionService interface {
	Follow(follow Interaction) error
	Unfollow(id string) error
	GetFollowings(id string) ([]Interaction, error)
	GetFollowers(id string) ([]Interaction, error)
}

func NewInteractionService(ddb *dynamodb.DynamoDB, tableName string) InteractionService {
	return &interactionService{
		ddb:       ddb,
		tableName: tableName,
	}
}

func (service *interactionService) Follow(follow Interaction) error {
	return nil
}

func (service *interactionService) Unfollow(userId string) error {
	return nil
}

func (service *interactionService) GetFollowings(userId string) ([]Interaction, error) {
	return nil, nil
}

func (service *interactionService) GetFollowers(userId string) ([]Interaction, error) {
	return nil, nil
}
