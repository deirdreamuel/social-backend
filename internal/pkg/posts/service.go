package posts

import (
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Post struct {
	Id          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Description string    `json:"description"`
	Media       string    `json:"media"`
}

type postService struct {
	ddb       *dynamodb.DynamoDB
	tableName string 
}

type PostService interface {
	PutPost(post Post) error
	GetPost(id string) (Post, error)
	DeletePost(id string) error
}

func NewPostService(ddb *dynamodb.DynamoDB, tableName string) PostService {
	return &postService{
		ddb:       ddb,
		tableName: tableName,
	}
}

func (service *postService) PutPost(follow Post) error {
	return nil
}

func (service *postService) GetPost(id string) (Post, error) {
	return Post{}, nil
}

func (service *postService) DeletePost(id string) error {
	return nil
}
