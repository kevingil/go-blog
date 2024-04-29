package storage

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

type Storage struct {
	S3Session *session.Session
}

type StorageInterface interface {
	Connect() (*session.Session, error)
	List(bucket, prefix string) ([]File, []string, error)
	Upload(bucket, key string, data []byte) error
}

type Config struct {
	AccessKey    string
	SecretKey    string
	SessionToken string
	Endpoint     string
	Region       string
}

const (
	MaxImageSize = 100 * 1024 // 100KB
	HIDDEN       = "./hidden"
	CACHE        = "./cache"
)

// Connect to S3 endpoint with ENV variables
func NewSession(c Config) (*Storage, error) {
	config := &aws.Config{
		Region:      aws.String(c.Region), // Adjust according to your setup
		Credentials: credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, c.SessionToken),
		Endpoint:    aws.String(c.Endpoint),
	}

	s3Session := session.Must(session.NewSessionWithOptions(session.Options{
		Config: *config,
	}))

	return &Storage{S3Session: s3Session}, nil
}
