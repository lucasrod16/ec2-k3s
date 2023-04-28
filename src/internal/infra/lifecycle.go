package infra

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	projectName      string = "ec2-k3s"
	stackName        string = "dev"
	awsPluginVersion string = "v5.37.0"
)

// Up provisions AWS infrastructure
func Up(region, instanceType string) error {
	pulumiStack, ctx := configurePulumi(region, instanceType)

	// Wire up our update to stream progress to stdout
	stdoutStreamer := optup.ProgressStreams(os.Stdout)

	// Run the update to deploy our infrastructure
	if _, err := pulumiStack.Up(ctx, stdoutStreamer); err != nil {
		return err
	}

	// Wait for ec2 instance to be ready
	if err := WaitInstanceReady(region); err != nil {
		return err
	}

	// Install k3s on ec2 instance
	if err := InstallK3s(region); err != nil {
		return err
	}

	// Copy kubeconfig from remote host to local machine
	GetKubeconfig(region)

	return nil
}

// Down tears down AWS infrastructure
func Down(region, instanceType string) error {
	pulumiStack, ctx := configurePulumi(region, instanceType)

	// Wire up our destroy to stream progress to stdout
	stdoutStreamer := optdestroy.ProgressStreams(os.Stdout)

	// Destroy resources in the stack
	if _, err := pulumiStack.Destroy(ctx, stdoutStreamer); err != nil {
		return err
	}

	opts := auto.LocalWorkspace{}

	// Destroy the stack
	if err := opts.RemoveStack(ctx, stackName); err != nil {
		return err
	}

	fmt.Printf("Stack '%s' has been removed\n", stackName)

	return nil
}

func deployInfra(instanceType string) pulumi.RunFunc {
	deployFunc := func(ctx *pulumi.Context) error {
		// Create SSH keypair in AWS
		if _, err := CreateSSHKeyPair(ctx); err != nil {
			return err
		}

		// Create ec2 instance and security group in AWS
		infra, err := CreateInstance(ctx, instanceType)
		if err != nil {
			return err
		}

		// Print outputs to stdout
		ctx.Export("Instance ID", infra.Server.ID())
		ctx.Export("Public IP Address", infra.Server.PublicIp)
		ctx.Export("Hostname", infra.Server.PublicDns)
		ctx.Export("Instance Type", infra.Server.InstanceType)
		ctx.Export("AMI ID", infra.Server.Ami)
		ctx.Export("Instance Tags", infra.Server.Tags)

		return nil
	}

	return deployFunc
}

func configurePulumi(region, instanceType string) (auto.Stack, context.Context) {
	ctx := context.Background()

	stack, err := auto.UpsertStackInlineSource(ctx, stackName, projectName, deployInfra(instanceType))
	if err != nil {
		log.Fatal(err)
	}

	workspace := stack.Workspace()

	// For inline source programs, we must manage plugins ourselves
	if err := workspace.InstallPlugin(ctx, "aws", awsPluginVersion); err != nil {
		log.Fatal(err)
	}

	// Set stack configuration specifying the AWS region to deploy
	if err := stack.SetConfig(ctx, "aws:region", auto.ConfigValue{Value: region}); err != nil {
		log.Fatal(err)
	}

	// Refresh state
	if _, err := stack.Refresh(ctx); err != nil {
		log.Fatal(err)
	}

	return stack, ctx
}
