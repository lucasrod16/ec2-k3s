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

// Up provisions AWS infrastructure
func Up(region, instanceType string) error {
	pulumiStack, ctx := configurePulumi(region, instanceType)

	// Wire up our update to stream progress to stdout
	stdoutStreamer := optup.ProgressStreams(os.Stdout)

	// Run the update to deploy our infrastructure
	s.Start()
	pterm.Info.Println("Updating stack...")
	if _, err := pulumiStack.Up(ctx, stdoutStreamer); err != nil {
		return err
	}
	s.Stop()

	pterm.Success.Println("Update succeeded!")

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

	s := spinner.New(spinner.CharSets[36], 1000*time.Millisecond)
	s.Start()

	pterm.Info.Println("Destroying the stack...")

	// Wire up our destroy to stream progress to stdout
	stdoutStreamer := optdestroy.ProgressStreams(os.Stdout)

	// Destroy our stack and exit early
	if _, err := pulumiStack.Destroy(ctx, stdoutStreamer); err != nil {
		return err
	}

	s.Stop()

	pterm.Success.Println("Stack successfully destroyed!")

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

		return nil
	}

	return deployFunc
}

func configurePulumi(region, instanceType string) (auto.Stack, context.Context) {
	ctx := context.Background()
	projectName := "ec2-k3s"
	stackName := "dev"

	stack, _ := auto.UpsertStackInlineSource(ctx, stackName, projectName, deployInfra(instanceType))

	pterm.Info.Println("Created/Selected stack " + stackName)

	workspace := stack.Workspace()

	// For inline source programs, we must manage plugins ourselves
	s.Start()
	pterm.Info.Println("Installing AWS plugin...")
	workspace.InstallPlugin(ctx, "aws", "v5.25.0")
	s.Stop()

	pterm.Success.Println("Successfully installed AWS plugin")

	// Set stack configuration specifying the AWS region to deploy
	stack.SetConfig(ctx, "aws:region", auto.ConfigValue{Value: region})
	pterm.Success.Println("Successfully set config")

	// Refresh state
	s.Start()
	pterm.Info.Println("Refreshing state...")

	if _, err := stack.Refresh(ctx); err != nil {
		pterm.Fatal.Printf("%v\n", err)
	}

	s.Stop()
	pterm.Success.Println("Refresh succeeded!")

	return stack, ctx
}
