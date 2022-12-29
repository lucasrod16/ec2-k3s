/*
Copyright © 2022 Lucas Rodriguez
*/
package cmd

import (
	"github.com/lucasrod16/ec2-k3s/internal/infra"
	"github.com/spf13/cobra"
)

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Teardown AWS infrastructure",
	Long:  "Teardown AWS infrastructure managed by Pulumi.",
	Run: func(cmd *cobra.Command, args []string) {
		infra.Down()
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
}
