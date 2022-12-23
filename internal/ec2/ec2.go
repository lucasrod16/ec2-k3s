package ec2

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Create() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create security group
		sg, err := ec2.NewSecurityGroup(ctx, "lucas-dev-sg", &ec2.SecurityGroupArgs{
			Description: pulumi.String("Allow all inbound traffic from the workstation IP address only"),
			Ingress: ec2.SecurityGroupIngressArray{
				&ec2.SecurityGroupIngressArgs{
					Description: pulumi.String("All ports and protocols from workstation IP"),
					FromPort:    pulumi.Int(0),
					ToPort:      pulumi.Int(0),
					Protocol:    pulumi.String("-1"),
					CidrBlocks: pulumi.StringArray{
						pulumi.String(localIP()),
					},
				},
			},
			Egress: ec2.SecurityGroupEgressArray{
				&ec2.SecurityGroupEgressArgs{
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
			return err
		}

		// Create SSH keypair
		keypair, err := ec2.NewKeyPair(ctx, "lucas-dev-ssh", &ec2.KeyPairArgs{
			KeyName:   pulumi.String("ec2-k3d-keypair"),
			PublicKey: pulumi.String(getPublicSSHKey()),
		})
		if err != nil {
			return err
		}
		fmt.Println(keypair)

		// Get Ubuntu 22.04 AMI ID
		ubuntu, err := ec2.LookupAmi(ctx, &ec2.LookupAmiArgs{
			MostRecent: pulumi.BoolRef(true),
			Filters: []ec2.GetAmiFilter{
				ec2.GetAmiFilter{
					Name: "name",
					Values: []string{
						"ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*",
					},
				},
				ec2.GetAmiFilter{
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
			return err
		}

		// Create ec2 instance
		_, err = ec2.NewInstance(ctx, "lucas-dev-ec2", &ec2.InstanceArgs{
			Ami:                 pulumi.String(ubuntu.Id),
			InstanceType:        pulumi.String("t3.2xlarge"),
			KeyName:             pulumi.String("ec2-k3d-keypair"),
			VpcSecurityGroupIds: pulumi.StringArray{sg.ID()},
			Tags: pulumi.StringMap{
				"Name":         pulumi.String("lucas-dev"),
				"Cluster-type": pulumi.String("k3d"),
				"Workload":     pulumi.String("bigbang"),
			},
			UserData: pulumi.String(getUserDataScript()),
		})
		if err != nil {
			return err
		}
		return nil
	})
}

// localIP returns the IP address of the machine that executed the program
func localIP() []byte {
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

// getPublicSSHKey returns the public ssh key at ~/.ssh/id_rsa.pub
func getPublicSSHKey() []byte {
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

// getUserDataScript returns the user data script in this repo at hack/user-data.sh
func getUserDataScript() []byte {
	userDataScript := "hack/user-data.sh"
	scriptData, err := os.ReadFile(userDataScript)
	if err != nil {
		log.Panicf("Failed reading data from user data script: %s", err)
	}

	return scriptData
}
