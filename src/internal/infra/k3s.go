package infra

import (
	"os"

	execute "github.com/alexellis/go-execute/pkg/v1"
	"github.com/pterm/pterm"
)

/*
K3sUp creates a k3s cluster on ec2 instance via SSH
and copies the kubeconfig to the directory where
`ec2-k3s up` was executed.
*/
func K3sUp(region string) error {
	publicIp, err := GetInstanceIp(region)
	if err != nil {
		return err
	}

	userName := "ubuntu"
	disableTraefik := "--disable traefik"

	task := execute.ExecTask{
		Command: "k3sup",
		Args: []string{
			"install",
			"--ip=" + publicIp,
			"--user=" + userName,
			"--k3s-extra-args=" + disableTraefik,
		},
		StreamStdio: false,
	}

	s.Start()
	pterm.Info.Println("Creating k3s cluster on ec2 instance...")
	_, err = task.Execute()
	if err != nil {
		return err
	}
	s.Stop()

	pterm.Success.Println("Successfully created k3s cluster on ec2 instance!")

	return nil
}

// K3sReady waits for k3s cluster node to be ready
func K3sReady() error {
	kubeconfig := "kubeconfig"
	task := execute.ExecTask{
		Command: "k3sup",
		Args: []string{
			"ready",
			"--kubeconfig=" + kubeconfig,
		},
		StreamStdio: false,
	}

	s.Start()
	pterm.Info.Println("Waiting for k3s cluster node to be ready...")
	_, err := task.Execute()
	if err != nil {
		return err
	}
	s.Stop()

	pterm.Success.Println("Node is ready!")

	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}

	kubeconfigPath := workingDir + "/" + kubeconfig

	pterm.Info.Printf(`Access your cluster:
	export KUBECONFIG=%v
	kubectl get node -o wide`, kubeconfigPath)

	return nil
}
