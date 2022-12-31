package infra

import (
	execute "github.com/alexellis/go-execute/pkg/v1"
)

/*
K3sUp creates a k3s cluster on ec2 instance via SSH
and copies the kubeconfig to the directory where
`ec2-k3s up` was executed.
*/
func K3sUp() error {
	publicIp, err := GetInstanceIp()
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
		StreamStdio: true,
	}

	_, err = task.Execute()
	if err != nil {
		return err
	}

	return nil
}

// K3sReady waits for k3s cluster node to be ready
func K3sReady() error {
	kubeconfig := "./kubeconfig"
	task := execute.ExecTask{
		Command: "k3sup",
		Args: []string{
			"ready",
			"--kubeconfig=" + kubeconfig,
		},
		StreamStdio: true,
	}

	_, err := task.Execute()
	if err != nil {
		return err
	}

	return nil
}
