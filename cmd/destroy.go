/*
Copyright © 2022 Lucas Rodriguez
*/
package cmd

import (
	"github.com/lucasrod16/ec2-k3s/internal/infra"
	"github.com/spf13/cobra"
)

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Tear down cloud infrastructure",
	Long:  "Tear down cloud infrastructure managed by Pulumi.",
	Run: func(cmd *cobra.Command, args []string) {
		infra.Down()
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)
}
