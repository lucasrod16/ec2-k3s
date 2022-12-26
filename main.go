package main

import (
	"main/internal/ec2"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create SSH keypair
		_, err := ec2.CreateSSHKeyPair(ctx)
		if err != nil {
			return err
		}

		// Create ec2 instance and security group
		infra, err := ec2.CreateInstance(ctx)
		if err != nil {
			return err
		}

		// Print infrastructure details to stdout
		ctx.Export("EC2 Instance", infra.Server.ID())
		ctx.Export("publicIp", infra.Server.PublicIp)
		ctx.Export("publicHostName", infra.Server.PublicDns)
		return nil
	})
}
