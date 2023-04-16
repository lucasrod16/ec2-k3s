package utils

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/google/uuid"
	"github.com/lucasrod16/ec2-k3s/src/internal/types"
)

var (
	InstanceOwner = createInstanceOwnerTag()
	cfg           = types.ConfigFile{}
)

// GetPublicSSHKey returns the public ssh key at the provided path.
func GetPublicSSHKey() []byte {
	keyData, err := os.ReadFile(cfg.SSHKeyPath + ".pub")
	if err != nil {
		log.Fatal(err)
	}

	return keyData
}

// GetPrivateSSHKey returns the private ssh key at the provided path.
func GetPrivateSSHKey() []byte {
	keyData, err := os.ReadFile(cfg.SSHKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	return keyData
}

// LocalIP returns the IP address of the machine that executed the program
func LocalIP() []byte {
	resp, err := http.Get("https://checkip.amazonaws.com")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	trimmedBody := bytes.Trim(body, "\n")
	suffix := "/32"
	cidr := append([]byte(trimmedBody), suffix...)

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
	}

	// Describe the status of running instances
	result, err := client.DescribeInstanceStatus(input)
	if err != nil {
		return "", err
	}

	// Convert string pointer to string
	instanceStatusPointer := result.InstanceStatuses[0].InstanceStatus.Details[0].Status
	instanceStatus := aws.StringValue(instanceStatusPointer)

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
	publicIpAddress := aws.StringValue(publicIpAddressPointer)

	return publicIpAddress, nil
}

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
	instanceId := aws.StringValue(instanceIdPointer)

	return instanceId, nil
}

// createInstanceOwnerTag creates a unique name for the ec2 instance owner tag value
func createInstanceOwnerTag() string {
	instanceOwner := GetCurrentUser() + "-" + createUUID()

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

// createUUID creates a uuid
func createUUID() string {
	uuid := uuid.New()
	uuidString := fmt.Sprintf("%v", uuid)

	return uuidString
}
