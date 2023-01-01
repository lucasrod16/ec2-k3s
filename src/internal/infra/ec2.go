package infra

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/briandowns/spinner"
	"github.com/lucasrod16/ec2-k3s/src/internal/types"
	"github.com/lucasrod16/ec2-k3s/src/internal/utils"
	"github.com/pterm/pterm"
	"k8s.io/apimachinery/pkg/util/wait"

	pec2 "github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// CreateSecurityGroup creates a security group in AWS
func CreateSecurityGroup(ctx *pulumi.Context) (*types.Infrastructure, error) {
	securityGroup, err := pec2.NewSecurityGroup(ctx, "security-group", &pec2.SecurityGroupArgs{
		Description: pulumi.String("Allow all inbound traffic from the workstation IP address only"),
		Ingress: pec2.SecurityGroupIngressArray{
			&pec2.SecurityGroupIngressArgs{
				Description: pulumi.String("All ports and protocols from workstation IP"),
				FromPort:    pulumi.Int(0),
				ToPort:      pulumi.Int(0),
				Protocol:    pulumi.String("-1"),
				CidrBlocks: pulumi.StringArray{
					pulumi.String(utils.LocalIP()),
				},
			},
		},
		Egress: pec2.SecurityGroupEgressArray{
			&pec2.SecurityGroupEgressArgs{
				FromPort: pulumi.Int(0),
				ToPort:   pulumi.Int(0),
				Protocol: pulumi.String("-1"),
				CidrBlocks: pulumi.StringArray{
					pulumi.String("0.0.0.0/0"),
				},
			},
		},
		Tags: pulumi.StringMap{
			"Name": pulumi.String("allow all ports and protocols from workstation IP"),
		},
	})
	if err != nil {
		return nil, err
	}

	return &types.Infrastructure{
		SecurityGroup: securityGroup,
	}, nil
}

// CreateSSHKeyPair creates an SSH keypair in AWS
func CreateSSHKeyPair(ctx *pulumi.Context) (*types.Infrastructure, error) {
	keypair, err := pec2.NewKeyPair(ctx, "ssh-keypair", &pec2.KeyPairArgs{
		KeyName:   pulumi.String("ec2-k3s-keypair"),
		PublicKey: pulumi.String(utils.GetPublicSSHKey()),
	})
	if err != nil {
		return nil, err
	}

	return &types.Infrastructure{
		Keypair: keypair,
	}, nil
}

// CreateInstance creates an ec2 instance in AWS
func CreateInstance(ctx *pulumi.Context) (*types.Infrastructure, error) {
	computeInfra, err := getUbuntuAMI(ctx)
	if err != nil {
		return nil, err
	}

	securityInfra, err := CreateSecurityGroup(ctx)
	if err != nil {
		return nil, err
	}

	server, err := pec2.NewInstance(ctx, "ec2-instance", &pec2.InstanceArgs{
		Ami:                 pulumi.String(computeInfra.Ami.ImageId),
		InstanceType:        pulumi.String("t3.2xlarge"),
		KeyName:             pulumi.String("ec2-k3s-keypair"),
		VpcSecurityGroupIds: pulumi.StringArray{securityInfra.SecurityGroup.ID()},
		Tags: pulumi.StringMap{
			"Name":         pulumi.String("ec2-k3s"),
			"Cluster-type": pulumi.String("k3s"),
		},
	})
	if err != nil {
		return nil, err
	}

	return &types.Infrastructure{
		Server: server,
	}, nil
}

// getUbuntuAMI returns the latest Ubuntu 22.04 AMI ID
func getUbuntuAMI(ctx *pulumi.Context) (*types.Infrastructure, error) {
	ami, err := pec2.LookupAmi(ctx, &pec2.LookupAmiArgs{
		MostRecent: pulumi.BoolRef(true),
		Filters: []pec2.GetAmiFilter{
			pec2.GetAmiFilter{
				Name: "name",
				Values: []string{
					"ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*",
				},
			},
			pec2.GetAmiFilter{
				Name: "virtualization-type",
				Values: []string{
					"hvm",
				},
			},
		},
		Owners: []string{
			"099720109477",
		},
	}, nil)
	if err != nil {
		return nil, err
	}

	return &types.Infrastructure{
		Ami: ami,
	}, nil
}

// GetInstanceIp returns the public IP address of the ec2 instance
func GetInstanceIp() (string, error) {
	client := utils.SetupEC2Client()

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
	publicIpAddress := utils.DerefString(publicIpAddressPointer)

	return publicIpAddress, nil
}

// GetInstanceStatus returns the reachability status of the ec2 instance
func GetInstanceStatus() (string, error) {
	client := utils.SetupEC2Client()

	input := &ec2.DescribeInstanceStatusInput{
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
	instanceStatus := utils.DerefString(instanceStatusPointer)

	return instanceStatus, nil
}

// WaitInstanceReady waits for instance health checks to return "passed"
func WaitInstanceReady() error {
	s := spinner.New(spinner.CharSets[36], 1000*time.Millisecond)
	s.Start()

	pterm.Info.Println("Waiting for ec2 instance to be ready...")

	err := wait.Poll(1*time.Second, 3*time.Minute, func() (bool, error) {
		status, err := GetInstanceStatus()
		if err != nil {
			return false, err
		}

		if status == "passed" {
			s.Stop()
			pterm.Success.Println("Instance is ready!")
			return true, nil
		}

		return false, nil
	})
	if err != nil {
		return err
	}

	return nil
}
