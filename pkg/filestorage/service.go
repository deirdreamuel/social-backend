package filestorage

import (
	"mime/multipart"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type _Service struct {
	client     *s3.S3
	bucketName string
}

// Service interface which contains file storage operations
type Service interface {
	GetUploadUrl(filename string) (string, error)
	UploadFile(filename string, file *multipart.File) error
}

// NewDatabaseService function to initialize filestorage.Service object
func NewFileStorageService(bucketName string) Service {
	// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-credentials
	// Initialize session and config for initializing client
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	cfg := &aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	}

	client := s3.New(sess, cfg)
	return &_Service{client: client, bucketName: bucketName}
}

// GetUploadUrl to get url to upload directly to storage instead of going through the server
func (svc *_Service) GetUploadUrl(filename string) (string, error) {
	response, _ := svc.client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: &svc.bucketName,
		Key:    &filename,
	})

	// Get presign url which allows upload with necessary permission
	return response.Presign(5 * time.Minute)
}

// UploadFile uploads file to storage
func (svc *_Service) UploadFile(filename string, file *multipart.File) error {
	_, err := svc.client.PutObject(&s3.PutObjectInput{
		Bucket: &svc.bucketName,
		Key:    &filename,
		Body:   *file,
	})

	if err != nil {
		return err
	}

	return nil
}
