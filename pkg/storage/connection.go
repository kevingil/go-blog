package storage

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
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

type File struct {
	Key          string
	Name         string
	Path         string
	Size         int64
	LastModified time.Time
}

type Folder struct {
	Name  string
	Path  string
	Files []File
}

const (
	MaxImageSize = 100 * 1024 // 100KB
	HIDDEN       = "./hidden"
	CACHE        = "./cache"
)

// Connect to S3 endpoint
func (s *Storage) Connect() (*session.Session, error) {
	if s.S3Session == nil {
		config := &aws.Config{
			Region: aws.String("us-west-2"),
		}
		s.S3Session = session.Must(session.NewSessionWithOptions(session.Options{
			Config: *config,
		}))
	}
	return s.S3Session, nil
}
