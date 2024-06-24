package storage

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Session struct {
	BucketName      string
	AccountId       string
	AccessKeyId     string
	AccessKeySecret string
	Endpoint        string
	Region          string

	Client    *s3.Client
	UrlPrefix string
}

type SessionInterface interface {
	Connect() (*s3.Client, error)
	List(bucket, prefix string) ([]File, []string, error)
	Upload(bucket, key string, data []byte) error
}

type Config struct {
}

const (
	MaxImageSize = 100 * 1024 // 100KB
	HIDDEN       = "./hidden"
	CACHE        = "./cache"
)

// Connect to S3 endpoint with ENV variables
func (s Session) Connect() (Session, error) {

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: s.Endpoint,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(s.AccessKeyId, s.AccessKeySecret, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		log.Fatal(err)
	}

	s.Client = s3.NewFromConfig(cfg)

	return s, err
}
