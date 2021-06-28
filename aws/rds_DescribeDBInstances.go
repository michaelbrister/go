package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

var (
	profile string
	region  string
)

func main() {
	// Load the Shared AWS Configuration (~/.aws/config)
	profile = "aws-dev"
	region = "us-east-1"
	config, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile), config.WithRegion(region))

	// Create an Amazon service client
	client := rds.NewFromConfig(config)

	output, err := client.DescribeDBInstances(context.TODO(), &rds.DescribeDBInstancesInput{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("List RDS Instances in: " + profile)
	for _, object := range output.DBInstances {
		// Lookup RDS endpoint for database
		result, err := client.DescribeDBClusterEndpoints(context.TODO(), &rds.DescribeDBClusterEndpointsInput{
			DBClusterIdentifier: *&object.DBClusterIdentifier,
		})
		if err != nil {
			log.Fatal(err)
		}

		// log.Println("RDS Instance Identifier: " + *object.DBInstanceIdentifier)
		log.Println("RDS Cluster Identifier: " + *object.DBClusterIdentifier)

		// log.Println("RDS cluster endpoints")
		for _, endpoint := range result.DBClusterEndpoints {
			log.Println("RDS Cluster Endpoint: " + *endpoint.EndpointType)
		}
	}
}
