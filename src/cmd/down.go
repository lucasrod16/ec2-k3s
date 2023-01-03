/*
Copyright © 2022 Lucas Rodriguez
*/
package cmd

import (
	"github.com/lucasrod16/ec2-k3s/src/internal/infra"
	"github.com/spf13/cobra"
)

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Args:  cobra.MaximumNArgs(0),
	Short: "Teardown AWS infrastructure and k3s cluster",
	Run: func(cmd *cobra.Command, args []string) {
		infra.Down(o.InstanceType)
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
}
