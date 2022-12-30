package infra

import (
	"context"
	"fmt"
	"os"

	"github.com/pterm/pterm"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func deployInfra() pulumi.RunFunc {
	deployFunc := func(ctx *pulumi.Context) error {
		// Create SSH keypair in AWS
		_, err := CreateSSHKeyPair(ctx)
		if err != nil {
			return err
		}

		// Create ec2 instance and security group in AWS
		infra, err := CreateInstance(ctx)
		if err != nil {
			return err
		}

		// Print outputs to stdout
		ctx.Export("instanceId", infra.Server.ID())
		ctx.Export("publicIp", infra.Server.PublicIp)
		ctx.Export("hostname", infra.Server.PublicDns)

		return nil
	}

	return deployFunc
}

func ConfigurePulumi() (auto.Stack, context.Context) {
	ctx := context.Background()

	projectName := "ec2-k3s"
	stackName := "dev"

	stack, _ := auto.UpsertStackInlineSource(ctx, stackName, projectName, deployInfra())

	pterm.Println(pterm.Green("Created/Selected stack " + stackName))

	workspace := stack.Workspace()

	// For inline source programs, we must manage plugins ourselves
	pterm.Println(pterm.Cyan("Installing AWS plugin..."))
	workspace.InstallPlugin(ctx, "aws", "v5.25.0")

	pterm.Println(pterm.Green("Successfully installed AWS plugin"))

	// Set stack configuration specifying the AWS region to deploy
	stack.SetConfig(ctx, "aws:region", auto.ConfigValue{Value: "us-east-1"})

	pterm.Println(pterm.Green("Successfully set config"))

	// Refresh state
	pterm.Println(pterm.Cyan("Refreshing state..."))
	_, err := stack.Refresh(ctx)

	if err != nil {
		fmt.Printf("Failed to refresh stack: %v\n", err)
		os.Exit(1)
	}
	pterm.Println(pterm.Green("Refresh succeeded!"))

	return stack, ctx
}

// Up provisions AWS infrastructure
func Up() error {
	pulumiStack, ctx := ConfigurePulumi()

	pterm.Println(pterm.Cyan("Updating stack..."))

	// Wire up our update to stream progress to stdout
	stdoutStreamer := optup.ProgressStreams(os.Stdout)

	// Run the update to deploy our infrastructure
	_, err := pulumiStack.Up(ctx, stdoutStreamer)

	if err != nil {
		return err
	}

	pterm.Println(pterm.Green("Update succeeded!"))

	// Wait for ec2 instance to be ready
	WaitInstanceReady()

	// TODO: create k3s cluster

	return nil
}

// Down tears down AWS infrastructure
func Down() error {
	pulumiStack, ctx := ConfigurePulumi()

	pterm.Println(pterm.Red("Starting stack destroy"))

	// Wire up our destroy to stream progress to stdout
	stdoutStreamer := optdestroy.ProgressStreams(os.Stdout)

	// Destroy our stack and exit early
	_, err := pulumiStack.Destroy(ctx, stdoutStreamer)
	if err != nil {
		return err
	}

	pterm.Println(pterm.Red("Stack successfully destroyed"))

	return nil
}
