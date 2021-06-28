package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var profile string

func main() {
	// Load the Shared AWS Configuration (~/.aws/config)
	profile = "aws-dev"
	config, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile), config.WithRegion("us-east-1"))

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(config)

	// Get the first page of results for ListObjectsV2 for a bucket
	result, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})

	if err != nil {
		log.Fatal(err)
	}

	for _, bucket := range result.Buckets {
		fmt.Println(*bucket.Name)
	}
}
