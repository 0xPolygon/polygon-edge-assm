package s3storage

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"log"
	"strings"
)

type Adapter struct {
	s3         *s3.Client
	bucketName *string
}

func NewAdapter(region, bucketName string) (*Adapter, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("could not load new aws config err=%w", err)
	}

	return &Adapter{
		s3:         s3.NewFromConfig(cfg),
		bucketName: aws.String(bucketName),
	}, nil
}

func (a Adapter) WriteData(key, data string) error {
	log.Println("writing data to s3: ", data)

	input := &s3.PutObjectInput{
		Bucket: a.bucketName,
		Key:    aws.String(key),
		Body:   strings.NewReader(data),
	}

	_, err := a.s3.PutObject(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("could not put object to S3 err=%w", err)
	}

	return nil
}

func (a Adapter) FetchData(key string) (string, error) {
	// if object does not exist return empty json string
	// because we will create a new object
	listInput := &s3.ListObjectsV2Input{Bucket: a.bucketName}

	s3Objects, err := a.s3.ListObjectsV2(context.TODO(), listInput)
	if err != nil {
		return "", fmt.Errorf("could not list s3 bucket err=%w", err)
	}
	// fetch and return the data only if our S3 object is found
	for _, object := range s3Objects.Contents {
		if *object.Key == key {
			input := &s3.GetObjectInput{
				Bucket: a.bucketName,
				Key:    aws.String(key),
			}
			// fetch s3 object
			s3Object, err := a.s3.GetObject(context.TODO(), input)
			if err != nil {
				return "", fmt.Errorf("could not fetch S3 object err=%w", err)
			}

			// read from stream and return as string
			buf := new(strings.Builder)

			_, err = io.Copy(buf, s3Object.Body)
			if err != nil {
				return "", fmt.Errorf("could not read s3Object body err=%w", err)
			}

			fmt.Println("read data from s3: ", buf.String())

			return buf.String(), nil
		}
	}

	return "{}", nil
}
