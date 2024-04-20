package storage

import (
	"bytes"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
	MaxImageSize = 100 * 1024 // 100KB
	HiddenDir    = "hidden"
	CacheDir     = "hidden/cache"
)

type Storage struct {
	S3Session *session.Session
}

type StorageInterface interface {
	Connect() (*session.Session, error)
	List(bucket, prefix string) (string, []File, error)
	Upload(bucket, key string, data []byte) (File, error)
}

type Bucket struct {
	Name string
}

type File struct {
	Key          string
	Name         string
	Path         string
	Size         int64
	Type         string
	Preview      string
	LastModified time.Time
}

// Connect to Amazon S3 or Cloudflare R2
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

// List returns a directory prefix, an array of files for current dir, and error if any
func (s *Storage) List(bucket, prefix string) (string, []File, error) {
	svc := s3.New(s.S3Session)
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	}
	resp, err := svc.ListObjectsV2(input)
	if err != nil {
		return "", nil, err
	}

	var files []File
	for _, item := range resp.Contents {
		file := File{
			Key:          *item.Key,
			LastModified: *item.LastModified,
			Size:         *item.Size,
		}
		files = append(files, file)
	}

	directory := fmt.Sprintf("%s/%s", bucket, prefix)
	return directory, files, nil
}

// Upload file
func (s *Storage) Upload(bucket, key string, data []byte) error {
	uploader := s3manager.NewUploader(s.S3Session)

	_, err := uploader.UploadWithContext(aws.BackgroundContext(), &s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	})

	return err
}

func (b *Bucket) String() string {
	return b.Name
}

func (f *File) String() string {
	return fmt.Sprintf("%s/%s", f.Path, f.Name)
}
