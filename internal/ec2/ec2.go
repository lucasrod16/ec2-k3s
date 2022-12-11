package ec2

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Install docker, k3d, and kubectl
const userData string = `#!/bin/bash
	apt-get update
	apt-get install \
			ca-certificates \
			curl \
			gnupg \
			lsb-release

	mkdir -p /etc/apt/keyrings
	curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg

	echo \
	"deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
	$(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null

	apt-get update
	apt-get install -y \
			docker-ce \
			docker-ce-cli \
			containerd.io

	groupadd docker
	usermod -aG docker ubuntu

	apt-get update

	curl -fsSLo /etc/apt/keyrings/kubernetes-archive-keyring.gpg https://packages.cloud.google.com/apt/doc/apt-key.gpg

	echo \
		"deb [signed-by=/etc/apt/keyrings/kubernetes-archive-keyring.gpg] https://apt.kubernetes.io/ kubernetes-xenial main" \
		| sudo tee /etc/apt/sources.list.d/kubernetes.list

	apt-get update
	apt-get install -y kubectl

	curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | bash
`

func Create() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create security group
		sg, err := ec2.NewSecurityGroup(ctx, "lucas-dev-sg", &ec2.SecurityGroupArgs{
			Description: pulumi.String("Allow all inbound traffic from the workstation IP address only"),
			Ingress: ec2.SecurityGroupIngressArray{
				&ec2.SecurityGroupIngressArgs{
					Description: pulumi.String("All ports and protocols from VPC"),
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
				"Name": pulumi.String("allow all"),
			},
		})
		if err != nil {
			return err
		}

		// Create SSH keypair
		keypair, err := ec2.NewKeyPair(ctx, "lucas-dev-ssh", &ec2.KeyPairArgs{
			KeyName:   pulumi.String("ec2-k3d-keypair"),
			PublicKey: pulumi.String("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDPgLhdRFNkCK/CaRZI9B1EFPf5Ax1NOvhmBN6rKXUNextSIk3t+ZDyN4iv19aUZcr3IhG/8I9AIGV1+n48aZgDyuPh9MgvVeZRXTOpUp15m80RXcTrFUP8ubTESh8BiYee4DfmUcfccXjyB00OT5GK0OXNiWIGPkElpHPnwmRnRQ6bHyx8HJVMKC0MVwZe+RgydylDasUGJm+gE4+4xc+7F587mT+R17IjS4MZkIkIwIApez+euDp8lqtRuH3AGYpQdxkz09WuSRUBwgWOf4FkpB5+NZtDO3of22QJvL/6PZ2numx/llNhTO6ya1VrWpPH4q3ghaxjZy+v/Mh1+QrUx8r1RF6GSs18iKaBFHQin/it1KDSLW6wbHMQgsU9JrolWT93bZkOaahBDYkPvubgnGBEZ9kDDTVVzowUUJ6QNu932JJJk98dp0Q346RumUhAcgVgjenTPwgs6DgMc2q21pY/96vVGyJ87BlhJpkqZlVQMHfqnmCC/1hd8GWFmKU= lucas@Lucass-MacBook-Pro.local"),
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
			UserData: pulumi.String(userData),
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
		print(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		print(err)
	}
	trimmedBody := bytes.Trim(body, "\n")
	suffix := "/32"
	cidr := append([]byte(trimmedBody), suffix...)
	fmt.Println(string(cidr))
	return cidr
}
