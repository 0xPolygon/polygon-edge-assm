package aws

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)


func StoreGenesis(path string) error{
	// read genesis.json file
	gen, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("could not read genesis.json file: %w", err)
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("could not load aws configurtation: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.Region = Region
	})

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(BucketName),
		Key: aws.String("genesis.json"),
		Body: strings.NewReader(string(gen)),
	})
	if err != nil {
		return fmt.Errorf("could not put genesis.json to S3 bucket %w",err)
	}

	return nil
}