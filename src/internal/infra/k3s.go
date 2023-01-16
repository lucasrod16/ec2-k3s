package infra

import ssh "github.com/lucasrod16/ec2-k3s/src/internal/ssh-client"

// InstallK3s installs k3s on an ec2 instance via SSH
func InstallK3s(region string) error {
	sshClient, err := ssh.ConfigureSSHClient(region)
	if err != nil {
		return err
	}

	// Close the underlying network connection
	defer sshClient.Close()

	installK3sCommand := "curl -sfL https://get.k3s.io | sh -s -"

	if _, err = sshClient.Execute(installK3sCommand); err != nil {
		return err
	}

	return nil
}
