package main

import (
	"bufio"
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	log "github.com/sirupsen/logrus"
)

var (
	awsProfile    string
	awsRegion     string
	instanceId    string
	encryptedPass string
	keyFile       = "aws-dev-us-east-1.pem"
)

type instanceDetails struct {
	id        string
	name      string
	privateIp string
	platform  string
	password  string
}

func init() {
	//flag setup
	flag.StringVar(&awsProfile, "profile", "aws-dev", "AWS Profile")
	flag.StringVar(&awsRegion, "region", "us-east-1", "AWS Region")
	flag.StringVar(&instanceId, "id", "", "Ec2 Instance Id")
	flag.Parse()
}

func checkParams() bool {
	ok := true

	if len(instanceId) <= 0 {
		log.Error("-id is required")
		ok = false
	}
	return ok
}

func printRequiredParamsMessage() {
	fmt.Println("id is required to be set:")
	fmt.Println("Examples:")
	fmt.Println("-id is AWS ec2 Instance Id")
	fmt.Println("get-instance-details -profile aws-dev -id i-xxxxxx")
}

func getInstances(client *ec2.Client, ec2Instance *instanceDetails, instanceId string) {
	result, err := client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
		Filters: []ec2Types.Filter{
			{
				Name:   aws.String("instance-id"),
				Values: []string{instanceId},
			},
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	for _, r := range result.Reservations {
		for _, i := range r.Instances {
			ec2Instance.id = *i.InstanceId
			ec2Instance.privateIp = *i.PrivateIpAddress
			ec2Instance.platform = *i.PlatformDetails
		}
	}
}

func getEncryptedEc2Password(client *ec2.Client, ec2Instance *instanceDetails) []byte {
	encrPass, err := client.GetPasswordData(context.TODO(), &ec2.GetPasswordDataInput{
		InstanceId: &ec2Instance.id,
	})

	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Error(err)
	}

	if len(*encrPass.PasswordData) == 0 {
		log.Error("Unable to lookup password")
	}

	encryptedPasswd, err := base64.StdEncoding.DecodeString(*encrPass.PasswordData)

	if err != nil {
		log.Fatal(err)
	}

	return encryptedPasswd
}

func getKey() (*rsa.PrivateKey, error) {
	file, err := os.Open(keyFile)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	//Create a byte slice (pemBytes) the size of the file size
	pemFileInfo, _ := file.Stat()
	var size = pemFileInfo.Size()
	pemBytes := make([]byte, size)

	//Create new reader for the file and read into pemBytes
	buffer := bufio.NewReader(file)
	_, err = buffer.Read(pemBytes)
	if err != nil {
		return nil, err
	}

	//Now decode the byte slice
	data, _ := pem.Decode(pemBytes)
	if data == nil {
		return nil, errors.New("could not read pem file")
	}

	privKey, err := x509.ParsePKCS1PrivateKey(data.Bytes)
	if err != nil {
		return nil, err
	}

	return privKey, nil
}

func decryptEc2Password(cipherText []byte, privKey *rsa.PrivateKey) string {
	plaintext, err := rsa.DecryptPKCS1v15(nil, privKey, cipherText)
	if err != nil {
		log.Error(err)
	}

	return string(plaintext)
}

func main() {
	if !checkParams() {
		printRequiredParamsMessage()
		log.Fatal("Missing required params, exiting...")
	}

	config, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(awsProfile), config.WithRegion(awsRegion))

	if err != nil {
		log.Fatal("Failed to load configuration.")
	}

	// Create new EC2 client
	client := ec2.NewFromConfig(config)

	ec2Instance := instanceDetails{}

	getInstances(client, &ec2Instance, instanceId)

	if err != nil {
		log.Fatal(err)
	}

	if ec2Instance.platform == "Windows" {
		encryptedPass := getEncryptedEc2Password(client, &ec2Instance)

		//get key from keyfile
		privKey, err := getKey()

		if err != nil {
			log.Error("Error getting key ", err)
		}

		log.Info("Decrypting password")
		ec2Instance.password = decryptEc2Password(encryptedPass, privKey)
	}

	fmt.Println("Instance Id: " + ec2Instance.id)
	fmt.Println("Platform: " + ec2Instance.platform)
	if ec2Instance.platform == "Windows" {
		fmt.Println("Password: " + ec2Instance.password)
	}
}
