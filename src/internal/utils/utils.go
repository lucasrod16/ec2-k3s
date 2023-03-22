package utils

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
	"path"

	"crypto/sha256"
	"encoding/hex"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const (
	publicKeyFile  string = ".ssh/id_rsa.pub"
	privateKeyFile string = ".ssh/id_rsa"
)

var InstanceOwner string = createInstanceOwnerTag()

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

	publicSSHKey := path.Join(userHomeDir, publicKeyFile)
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

	privateSSHKey := path.Join(userHomeDir, privateKeyFile)
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

// GetInstanceStatus returns the reachability status of the ec2 instance
func GetInstanceStatus(region string) (string, error) {
	client := SetupEC2Client(region)
	instanceId, err := getInstanceId(region)
	if err != nil {
		return "", err
	}

	input := &ec2.DescribeInstanceStatusInput{
		InstanceIds: []*string{
			aws.String(instanceId),
		},
		Filters: []*ec2.Filter{
			{
				Name: aws.String("instance-status.reachability"),
				Values: []*string{
					aws.String("passed"),
					aws.String("failed"),
					aws.String("initializing"),
					aws.String("insufficient-data"),
				},
			},
		},
	}

	// Describe the status of running instances
	result, err := client.DescribeInstanceStatus(input)
	if err != nil {
		return "", err
	}

	// Convert string pointer to string
	instanceStatusPointer := result.InstanceStatuses[0].InstanceStatus.Details[0].Status
	instanceStatus := DerefString(instanceStatusPointer)

	return instanceStatus, nil
}

// GetInstanceIp returns the public IP address of the ec2 instance
func GetInstanceIp(region string) (string, error) {
	client := SetupEC2Client(region)
	instanceId, err := getInstanceId(region)
	if err != nil {
		return "", err
	}

	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceId),
		},
	}

	// Get the IP address of the EC2 instance
	result, err := client.DescribeInstances(input)
	if err != nil {
		return "", err
	}

	// Convert string pointer to string
	publicIpAddressPointer := result.Reservations[0].Instances[0].PublicIpAddress
	publicIpAddress := DerefString(publicIpAddressPointer)

	return publicIpAddress, nil
}

// getInstanceId gets the instance ID of the EC2 instance
func getInstanceId(region string) (string, error) {
	client := SetupEC2Client(region)
	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:Owner"),
				Values: []*string{
					aws.String(InstanceOwner),
				},
			},
		},
	}

	// Get the instance ID of the EC2 instance
	result, err := client.DescribeInstances(input)
	if err != nil {
		return "", err
	}

	instanceIdPointer := result.Reservations[0].Instances[0].InstanceId
	instanceId := DerefString(instanceIdPointer)

	return instanceId, nil
}

// createInstanceOwnerName creates a unique name for the ec2 instance owner tag value
func createInstanceOwnerTag() string {
	instanceOwner := GetCurrentUser() + "-" + createHash()

	return instanceOwner
}

func GetCurrentUser() string {
	userData, err := user.Current()
	if err != nil {
		fmt.Println("failed to get the current user's name")
		os.Exit(1)
	}
	userName := userData.Username

	return userName
}

// createHash creates a unique hash
func createHash() string {
	hasher := sha256.New()
	stringHash := hex.EncodeToString(hasher.Sum(nil))

	return stringHash
}
