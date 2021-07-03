package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func main() {
	profile := "aws-dev"
	// Load the Shared AWS Configuration (~/.aws/config)
	config, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile))

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Create an Amazon ec2 service client
	ec2Client := ec2.NewFromConfig(config)

	// build the request
	resp, err := ec2Client.GetPasswordData()

}
