package infra

import (
	"context"
	"os"
	"time"

	"github.com/briandowns/spinner"
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

	pterm.Info.Println("Created/Selected stack " + stackName)

	workspace := stack.Workspace()

	// For inline source programs, we must manage plugins ourselves
	s.Start()
	pterm.Info.Println("Installing AWS plugin...")
	workspace.InstallPlugin(ctx, "aws", "v5.25.0")
	s.Stop()

	pterm.Success.Println("Successfully installed AWS plugin")

	// Set stack configuration specifying the AWS region to deploy
	stack.SetConfig(ctx, "aws:region", auto.ConfigValue{Value: "us-east-1"})
	pterm.Success.Println("Successfully set config")

	// Refresh state
	s.Start()
	pterm.Info.Println("Refreshing state...")
	_, err := stack.Refresh(ctx)
	s.Stop()

	if err != nil {
		pterm.Fatal.Printf("%v\n", err)
	}
	pterm.Success.Println("Refresh succeeded!")

	return stack, ctx
}

// Up provisions AWS infrastructure
func Up() error {
	pulumiStack, ctx := ConfigurePulumi()

	// Wire up our update to stream progress to stdout
	stdoutStreamer := optup.ProgressStreams(os.Stdout)

	// Run the update to deploy our infrastructure
	s.Start()
	pterm.Info.Println("Updating stack...")
	_, err := pulumiStack.Up(ctx, stdoutStreamer)
	if err != nil {
		return err
	}
	s.Stop()

	pterm.Success.Println("Update succeeded!")

	// Wait for ec2 instance to be ready
	WaitInstanceReady()

	// Create k3s cluster on ec2 instance
	K3sUp()

	// Wait for cluster node to be ready
	K3sReady()

	return nil
}

// Down tears down AWS infrastructure
func Down() error {
	pulumiStack, ctx := ConfigurePulumi()

	s := spinner.New(spinner.CharSets[36], 1000*time.Millisecond)
	s.Start()

	pterm.Info.Println("Destroying the stack...")

	// Wire up our destroy to stream progress to stdout
	stdoutStreamer := optdestroy.ProgressStreams(os.Stdout)

	// Destroy our stack and exit early
	_, err := pulumiStack.Destroy(ctx, stdoutStreamer)
	if err != nil {
		return err
	}

	s.Stop()

	pterm.Success.Println("Stack successfully destroyed!")

	return nil
}
