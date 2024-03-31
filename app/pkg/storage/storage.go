package storage

import (
	"bytes"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Storage struct {
	S3Session *session.Session
}

type SessionInterface interface {
	Connect() (*session.Session, error)
	Upload(bucket, key string, data []byte) error
	Download(bucket, key string, filePath string) error
}

type Bucket struct {
	Name string
}

type Directory struct {
	Bucket Bucket
	Path   string
}

type File struct {
	Directory Directory
	Name      string
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

// Download file
func (s *Storage) Download(bucket, key string, filePath string) error {
	downloader := s3manager.NewDownloader(s.S3Session)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = downloader.DownloadWithContext(aws.BackgroundContext(), file, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	return err
}

func (b *Bucket) String() string {
	return b.Name
}

func (d *Directory) String() string {
	return fmt.Sprintf("%s/%s", d.Bucket.Name, d.Path)
}

func (f *File) String() string {
	return fmt.Sprintf("%s/%s", f.Directory.String(), f.Name)
}
