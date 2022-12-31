/*
Copyright © 2022 Lucas Rodriguez
*/
package cmd

import (
	"github.com/lucasrod16/ec2-k3s/src/internal/infra"
	"github.com/spf13/cobra"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Provision AWS infrastructure and k3s cluster",
	Run: func(cmd *cobra.Command, args []string) {
		infra.Up()
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}
