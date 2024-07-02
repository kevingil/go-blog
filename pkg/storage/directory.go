package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"path/filepath"
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
	Name     string
	Path     string
	IsHidden bool
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
		folders = append(folders, folder)
	}

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
