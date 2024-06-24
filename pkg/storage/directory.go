package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type File struct {
	Key          string
	LastModified time.Time
	Size         int64
	Url          string
}

type Folder struct {
	Name  string
	Path  string
	Files []File
}

// List returns an array of files and common prefixes (folders)
func (s *Session) List(bucket, prefix string) ([]File, []string, error) {

	listObjectsOutput, err := s.Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket:    aws.String(bucket),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String("/"),
	})

	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}

	for _, object := range listObjectsOutput.Contents {
		obj, _ := json.MarshalIndent(object, "", "\t")
		fmt.Println(string(obj))
	}

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
			Size:         *item.Size,
			Url:          fmt.Sprintf("%s/%s", s.UrlPrefix, *item.Key),
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

	var folders []string
	for _, commonPrefix := range listObjectsOutput.CommonPrefixes {
		folders = append(folders, *commonPrefix.Prefix)
	}

	return files, folders, err
}

// Upload file to a specific directory
// func (s *Storage) Upload(bucket, key string, data []byte) error {
//}

// Check if a file or folder is hidden based on naming conventions (e.g., starts with a dot)
func IsHidden(key string) bool {
	return key[0] == '.'
}
