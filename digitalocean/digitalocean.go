package digitalocean

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/voxowl/objectstorage"
)

const (
	DIGITALOCEAN_SPACES_ENDPOINT = "https://%s.digitaloceanspaces.com"
)

type DigitalOceanConfig struct {
	Region     string
	Bucket     string
	AuthKey    string
	AuthSecret string
}

type DigitalOceanObjectStorageOpts struct {
	UsePathStyle bool
}

// DigitalOceanObjectStorage is an implementation of the ObjectStorage interface
// for DigitalOcean Spaces Object Storage. (https://www.digitalocean.com/products/spaces)
type DigitalOceanObjectStorage struct {
	region   string
	bucket   string
	s3Client *s3.Client
}

// NewDigitalOceanObjectStorage creates a new DigitalOceanObjectStorage instance
// authKey and authSecret are the credentials for the DigitalOcean Spaces API
// region is the region of the DigitalOcean Spaces API (e.g. "nyc3")
// bucket is the name of the DigitalOcean Spaces bucket (optional?)
func NewDigitalOceanObjectStorage(config DigitalOceanConfig, opts DigitalOceanObjectStorageOpts) (DigitalOceanObjectStorage, error) {

	// validate config
	if config.Region == "" {
		return DigitalOceanObjectStorage{}, fmt.Errorf("config.Region is not provided")
	}
	if config.Bucket == "" {
		return DigitalOceanObjectStorage{}, fmt.Errorf("config.Bucket is not provided")
	}

	// Create S3 client using the provided information
	S3Client := s3.New(s3.Options{
		Region:       config.Region,
		BaseEndpoint: aws.String(fmt.Sprintf(DIGITALOCEAN_SPACES_ENDPOINT, config.Region)),
		Credentials:  credentials.NewStaticCredentialsProvider(config.AuthKey, config.AuthSecret, ""),
		UsePathStyle: opts.UsePathStyle,
	})

	objStorage := DigitalOceanObjectStorage{
		region:   config.Region,
		bucket:   config.Bucket,
		s3Client: S3Client,
	}

	// test the connection
	err := objStorage.testConnection()
	if err != nil {
		return DigitalOceanObjectStorage{}, err
	}

	return objStorage, nil
}

// Implementation of the ObjectStorage interface

func (s DigitalOceanObjectStorage) Download(key string) (io.ReadCloser, error) {

	// Create the input parameters for getting the object
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	// Get the object from the bucket
	result, err := s.s3Client.GetObject(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	return result.Body, nil
}

func (s DigitalOceanObjectStorage) Upload(key string, dataReader io.Reader) error {

	// Define the parameters of the object you want to upload.
	input := &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   dataReader,
	}

	// Run the PutObject function with your parameters, catching for errors.
	_, err := s.s3Client.PutObject(context.TODO(), input)

	return err
}

func (s DigitalOceanObjectStorage) List(prefix string, opts objectstorage.ListOpts) ([]string, error) {

	if opts.Limit < 0 {
		return nil, fmt.Errorf("opts.Limit must be equal to or greater than 0")
	}

	// 0 means no limit
	var limit int = opts.Limit

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(prefix),
	}

	var files []string
	paginator := s3.NewListObjectsV2Paginator(s.s3Client, input)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, fmt.Errorf("failed to list objects: %w", err)
		}

		for _, obj := range output.Contents {
			files = append(files, *obj.Key)
			if limit != 0 && len(files) >= int(limit) {
				break
			}
		}

		if limit != 0 && len(files) >= int(limit) {
			break
		}
	}

	return files, nil
}

// TODO: Implement some kind of paginator API for DigitalOceanObjectStorage

// unexported methods

// Test the connection to the DigitalOcean Spaces API
func (s *DigitalOceanObjectStorage) testConnection() error {
	// test the connection to the DigitalOcean Spaces API
	_, err := s.List("", objectstorage.ListOpts{
		Limit: 1,
	})
	return err
}
