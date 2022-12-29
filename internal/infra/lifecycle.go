package infra

import (
	"context"
	"fmt"
	"os"

	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var pulumiStack, ctx = configurePulumi()

func deployInfra() pulumi.RunFunc {
	deployFunc := func(ctx *pulumi.Context) error {
		// Create SSH keypair in AWS
		_, err := CreateSSHKeyPair(ctx)
		if err != nil {
			return err
		}

		// Create ec2 instance and security group in AWS
		i, err := CreateInstance(ctx)
		if err != nil {
			return err
		}

		// Print outputs to stdout
		ctx.Export("instanceId", i.Server.ID())
		ctx.Export("publicIp", i.Server.PublicIp)
		ctx.Export("hostname", i.Server.PublicDns)

		return nil
	}

	return deployFunc
}

func configurePulumi() (auto.Stack, context.Context) {
	ctx := context.Background()

	projectName := "ec2-k3s"
	stackName := "dev"

	stack, _ := auto.UpsertStackInlineSource(ctx, stackName, projectName, deployInfra())

	fmt.Printf("Created/Selected stack %q\n", stackName)

	workspace := stack.Workspace()

	fmt.Println("Installing the AWS plugin")

	// For inline source programs, we must manage plugins ourselves
	workspace.InstallPlugin(ctx, "aws", "v5.25.0")

	fmt.Println("Successfully installed AWS plugin")

	// Set stack configuration specifying the AWS region to deploy
	stack.SetConfig(ctx, "aws:region", auto.ConfigValue{Value: "us-east-1"})

	fmt.Println("Successfully set config")

	fmt.Println("Starting refresh")

	_, err := stack.Refresh(ctx)
	if err != nil {
		fmt.Printf("Failed to refresh stack: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Refresh succeeded!")

	return stack, ctx
}

// Up provisions AWS infrastructure
func Up() error {
	fmt.Println("Starting update")

	// Wire up our update to stream progress to stdout
	stdoutStreamer := optup.ProgressStreams(os.Stdout)

	// Run the update to deploy our infrastructure
	_, err := pulumiStack.Up(ctx, stdoutStreamer)
	if err != nil {
		return err
	}

	fmt.Println("Update succeeded!")

	// Return the ec2 instance reachability status
	GetInstanceStatus(ctx)

	// TODO: Use GetInstanceStatus result to check for "passed" status before creating cluster

	return nil
}

// Down tears down AWS infrastructure
func Down() error {
	fmt.Println("Starting stack destroy")

	// Wire up our destroy to stream progress to stdout
	stdoutStreamer := optdestroy.ProgressStreams(os.Stdout)

	// Destroy our stack and exit early
	_, err := pulumiStack.Destroy(ctx, stdoutStreamer)
	if err != nil {
		return err
	}

	fmt.Println("Stack successfully destroyed")

	return nil
}
