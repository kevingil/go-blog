package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"mime/multipart"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type File struct {
	Key          string
	LastModified time.Time
	Size         string
	SizeRaw      int64
	Url          string
	IsImage      bool
}

type Folder struct {
	Name         string
	Path         string
	IsHidden     bool
	LastModified time.Time
	FileCount    int
}

// List returns an array of files and common prefixes (folders)
func (s *Session) List(bucket, prefix string) ([]File, []Folder, error) {

	listObjectsOutput, err := s.Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket:    aws.String(bucket),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String("/"),
	})

	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}
	/*
		List S3 API response, for debugging
			for _, object := range listObjectsOutput.Contents {
				obj, _ := json.MarshalIndent(object, "", "\t")
				fmt.Println(string(obj))
			}
	*/
	listBucketsOutput, err := s.Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		log.Fatal(err)
	}

	for _, object := range listBucketsOutput.Buckets {
		obj, _ := json.MarshalIndent(object, "", "\t")
		fmt.Println(string(obj))
	}

	// {
	// 		"CreationDate": "2022-05-18T17:19:59.645Z",
	// 		"Name": "sdk-example"
	// }

	var files []File
	for _, item := range listObjectsOutput.Contents {
		file := File{
			Key:          *item.Key,
			LastModified: *item.LastModified,
			Size:         formatByteSize(*item.Size),
			SizeRaw:      *item.Size,
			Url:          fmt.Sprintf("%s/%s", s.UrlPrefix, *item.Key),
			IsImage:      isImageFile(*item.Key),
		}
		files = append(files, file)
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].LastModified.After(files[j].LastModified)
	})

	//  {
	//  	"ChecksumAlgorithm": null,
	//  	"ETag": "\"eb2b891dc67b81755d2b726d9110af16\"",
	//  	"Key": "ferriswasm.png",
	//  	"LastModified": "2022-05-18T17:20:21.67Z",
	//  	"Owner": null,
	//  	"Size": 87671,
	//  	"StorageClass": "STANDARD"
	//  }

	var folders []Folder
	for _, commonPrefix := range listObjectsOutput.CommonPrefixes {
		folderPath := *commonPrefix.Prefix
		folder := Folder{
			Name:     filepath.Base(folderPath),
			Path:     folderPath,
			IsHidden: folderIsHidden(filepath.Base(folderPath)),
		}

		// Get the folder contents to determine LastModified and FileCount
		folderContents, err := s.Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
			Bucket: aws.String(bucket),
			Prefix: aws.String(folderPath),
		})

		if err == nil {
			folder.FileCount = len(folderContents.Contents)
			if folder.FileCount > 0 {
				folder.LastModified = *folderContents.Contents[0].LastModified
				for _, item := range folderContents.Contents {
					if item.LastModified.After(folder.LastModified) {
						folder.LastModified = *item.LastModified
					}
				}
			}
		}

		folders = append(folders, folder)
	}

	sort.Slice(folders, func(i, j int) bool {
		return folders[i].LastModified.After(folders[j].LastModified)
	})

	return files, folders, nil
}

// Upload file to a specific directory
// func (s *Storage) Upload(bucket, key string, data []byte) error {
//}

func formatByteSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	exp := int(math.Log(float64(size)) / math.Log(unit))
	prefix := "KMGTPE"[exp-1]
	return fmt.Sprintf("%.2f %ciB", float64(size)/math.Pow(unit, float64(exp)), prefix)
}

// Helper function to check if a file has an image extension
func isImageFile(key string) bool {
	imageExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp"}
	ext := strings.ToLower(filepath.Ext(key))
	for _, imageExt := range imageExtensions {
		if ext == imageExt {
			return true
		}
	}
	return false
}

// Check if a folder is hidden based on naming conventions (e.g., starts with a dot)
func folderIsHidden(folderName string) bool {
	return strings.HasPrefix(folderName, ".")
}

func (s *Session) Upload(bucket, key string, file multipart.File) error {
	// Read the file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Create the PutObjectInput
	input := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(fileContent),
	}

	// Upload the file to S3
	_, err = s.Client.PutObject(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}

func (s *Session) Delete(bucket, key string) error {
	// Create the DeleteObjectInput request
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	// Delete the object
	_, err := s.Client.DeleteObject(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

func (s *Session) CreateFolder(bucket, folderPath string) error {
	// Ensure the folder path ends with a slash
	if !strings.HasSuffix(folderPath, "/") {
		folderPath += "/"
	}

	// Create an empty object with the folder path as the key
	input := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(folderPath),
		Body:   bytes.NewReader([]byte{}),
	}

	// Upload the empty object to S3
	_, err := s.Client.PutObject(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to create folder: %w", err)
	}

	return nil
}

func (s *Session) UpdateFolder(bucket, oldPath, newPath string) error {
	// Ensure both paths end with a slash
	if !strings.HasSuffix(oldPath, "/") {
		oldPath += "/"
	}
	if !strings.HasSuffix(newPath, "/") {
		newPath += "/"
	}

	// List objects in the old folder
	listInput := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(oldPath),
	}

	result, err := s.Client.ListObjectsV2(context.TODO(), listInput)
	if err != nil {
		return fmt.Errorf("failed to list objects: %w", err)
	}

	// Copy objects to the new folder and delete from the old folder
	for _, object := range result.Contents {
		// Create the new key
		newKey := strings.Replace(*object.Key, oldPath, newPath, 1)

		// Copy the object
		copyInput := &s3.CopyObjectInput{
			Bucket:     aws.String(bucket),
			CopySource: aws.String(fmt.Sprintf("%s/%s", bucket, *object.Key)),
			Key:        aws.String(newKey),
		}

		_, err := s.Client.CopyObject(context.TODO(), copyInput)
		if err != nil {
			return fmt.Errorf("failed to copy object %s: %w", *object.Key, err)
		}

		// Delete the old object
		deleteInput := &s3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    object.Key,
		}

		_, err = s.Client.DeleteObject(context.TODO(), deleteInput)
		if err != nil {
			return fmt.Errorf("failed to delete object %s: %w", *object.Key, err)
		}
	}

	return nil
}
