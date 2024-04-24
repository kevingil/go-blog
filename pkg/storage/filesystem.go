package storage

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// List returns an array of files and common prefixes (folders)
func (s *Storage) List(bucket, prefix string) ([]File, []string, error) {
	svc := s3.New(s.S3Session)
	input := &s3.ListObjectsV2Input{
		Bucket:    aws.String(bucket),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String("/"),
	}
	resp, err := svc.ListObjectsV2(input)
	if err != nil {
		return nil, nil, err
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

	var folders []string
	for _, commonPrefix := range resp.CommonPrefixes {
		folders = append(folders, *commonPrefix.Prefix)
	}

	return files, folders, nil
}

// Upload file to a specific directory
func (s *Storage) Upload(bucket, key string, data []byte) error {
	uploader := s3manager.NewUploader(s.S3Session)

	_, err := uploader.UploadWithContext(aws.BackgroundContext(), &s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	})

	return err
}

// Check if a file or folder is hidden based on naming conventions (e.g., starts with a dot)
func IsHidden(key string) bool {
	return key[0] == '.'
}
