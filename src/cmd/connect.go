package cmd

import (
	"fmt"
	"os"

	"github.com/lucasrod16/ec2-k3s/src/internal/ssh-client"
	"github.com/spf13/cobra"
)

// connectCmd represents the connect command
var (
	command    string
	connectCmd = &cobra.Command{
		Use:   "connect",
		Args:  cobra.MaximumNArgs(0),
		Short: "Connect to the ec2 instance via SSH",
		Run: func(cmd *cobra.Command, args []string) {
			sshClient, err := ssh.ConfigureSSHClient(region)
			if err != nil {
				fmt.Println("Unable to connect to server via SSH")
				os.Exit(1)
			}

			// Close the underlying network connection
			defer sshClient.Close()

			sshClient.Execute(command)
		},
	}
)

func init() {
	connectCmd.Flags().StringVar(&command, "command", "", "commands to execute on the ec2 instance via SSH")
	connectCmd.Flags().StringVar(&region, "region", "", "AWS region the ec2 instance is in")

	connectCmd.MarkFlagRequired("command")
	connectCmd.MarkFlagRequired("region")

	rootCmd.AddCommand(connectCmd)
}
