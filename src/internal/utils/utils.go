package utils

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// DerefString Dereferences string pointers to strings
func DerefString(s *string) string {
	if s != nil {
		return *s
	}

	return ""
}

// GetPublicSSHKey returns the public ssh key at ~/.ssh/id_rsa.pub
func GetPublicSSHKey() []byte {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	publicSSHKey := userHomeDir + "/.ssh/id_rsa.pub"
	keyData, err := os.ReadFile(publicSSHKey)
	if err != nil {
		log.Panicf("Failed reading data from public ssh key: %s", err)
	}

	return keyData
}

// GetPrivateSSHKey returns the private ssh key at ~/.ssh/id_rsa
func GetPrivateSSHKey() []byte {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	privateSSHKey := userHomeDir + "/.ssh/id_rsa"
	keyData, err := os.ReadFile(privateSSHKey)
	if err != nil {
		log.Panicf("Failed reading data from private ssh key: %s", err)
	}

	return keyData
}

// LocalIP returns the IP address of the machine that executed the program
func LocalIP() []byte {
	resp, err := http.Get("https://checkip.amazonaws.com")
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	trimmedBody := bytes.Trim(body, "\n")
	suffix := "/32"
	cidr := append([]byte(trimmedBody), suffix...)

	fmt.Printf("\nWorkstation IP address: %s", cidr)

	return cidr
}

// SetupEC2Client configures a client to make EC2 API calls
func SetupEC2Client(region string) *ec2.EC2 {
	sess := session.Must(session.NewSession())
	svc := ec2.New(sess, aws.NewConfig().WithRegion(region))

	return svc
}

// GetInstanceIp returns the public IP address of the ec2 instance
func GetInstanceIp(region string) (string, error) {
	client := SetupEC2Client(region)
	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("ip-address"),
				Values: []*string{
					aws.String("*"),
				},
			},
		},
	}

	// Describe the status of running instances
	result, err := client.DescribeInstances(input)
	if err != nil {
		return "", err
	}

	// Convert string pointer to string
	publicIpAddressPointer := result.Reservations[0].Instances[0].PublicIpAddress
	publicIpAddress := DerefString(publicIpAddressPointer)

	return publicIpAddress, nil
}
