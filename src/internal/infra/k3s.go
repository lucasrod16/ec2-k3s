package infra

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	ssh "github.com/lucasrod16/ec2-k3s/src/internal/ssh-client"
	"github.com/lucasrod16/ec2-k3s/src/internal/utils"
)

// InstallK3s installs k3s on an ec2 instance via SSH
func InstallK3s(region string) error {
	sshClient, err := ssh.ConfigureSSHClient(region)
	if err != nil {
		return err
	}

	// Close the underlying network connection
	defer sshClient.Close()

	ip, err := utils.GetInstanceIp(region)
	if err != nil {
		return err
	}

	installK3sCommand := "curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC='--tls-san=" + ip + "' sh -s - --disable traefik"

	if _, err = sshClient.Execute(installK3sCommand); err != nil {
		return err
	}

	return nil
}

// GetKubeconfig fetches the kubeconfig from the remote host
// and writes it to working directory on local disk
func GetKubeconfig(region string) error {
	sshClient, err := ssh.ConfigureSSHClient(region)
	if err != nil {
		return err
	}

	getConfigCommand := "sudo cat /etc/rancher/k3s/k3s.yaml"

	output, err := sshClient.ExecuteOutput(getConfigCommand, false)
	if err != nil {
		return err
	}

	ip, err := utils.GetInstanceIp(region)
	if err != nil {
		return err

	}

	kubeconfig := editKubeconfig(string(output.StdOut), ip)

	if err := writeKubeconfig([]byte(kubeconfig)); err != nil {
		return err
	}

	return nil
}

// Edit kubeconfig file with public IP of ec2 instance to connect to
func editKubeconfig(kubeconfig string, ip string) []byte {
	kubeconfigChanges := strings.NewReplacer(
		"127.0.0.1", ip,
		"localhost", ip,
	)

	return []byte(kubeconfigChanges.Replace(kubeconfig))
}

// Write kubeconfig file to disk
func writeKubeconfig(data []byte) error {
	absPath, err := getAbsolutePath()
	if err != nil {
		return err
	}

	filePath := path.Join(absPath, "kubeconfig")

	if err := os.WriteFile(filePath, []byte(data), 0600); err != nil {
		return err
	}

	return nil
}

// Get the absolute path of the current working directory
func getAbsolutePath() (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	absPath, err := filepath.Abs(workingDir)
	if err != nil {
		return "", err
	}

	return absPath, nil
}
