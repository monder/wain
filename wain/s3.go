package wain

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	awsS3 "github.com/aws/aws-sdk-go/service/s3"

	"bytes"
	"io/ioutil"
)

type S3Connection struct {
	aws    *awsS3.S3
	bucket string
}

func CreateS3(config *Config) (map[string]*S3Connection, error) {
	s3 := make(map[string]*S3Connection)

	for _, bucket := range config.Buckets {

		svc := awsS3.New(session.New(), &aws.Config{
			Region:      aws.String(bucket.Region),
			Credentials: credentials.NewStaticCredentials(bucket.AccessKey, bucket.AccessSecret, ""),
		})

		s3[bucket.Name] = &S3Connection{svc, bucket.Name}
	}

	return s3, nil
}

func (s3 *S3Connection) GetObject(key string) ([]byte, error) {

	resp, err := s3.aws.GetObject(&awsS3.GetObjectInput{
		Bucket: aws.String(s3.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	return bytes, err
}

func (s3 *S3Connection) PutObject(key string, data []byte, contentType string) error {

	_, err := s3.aws.PutObject(&awsS3.PutObjectInput{
		Bucket:       aws.String(s3.bucket),
		Key:          aws.String(key),
		Body:         bytes.NewReader(data),
		StorageClass: aws.String("REDUCED_REDUNDANCY"),
		ACL:          aws.String("public-read"),
		ContentType:  aws.String(contentType),
	})

	return err
}
