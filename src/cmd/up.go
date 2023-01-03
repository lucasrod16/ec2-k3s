/*
Copyright © 2022 Lucas Rodriguez
*/
package cmd

import (
	"github.com/lucasrod16/ec2-k3s/src/internal/infra"
	"github.com/lucasrod16/ec2-k3s/src/internal/types"
	"github.com/spf13/cobra"
)

var o = &types.InstanceOptions{}

// upCmd represents the up command
var (
	upCmd = &cobra.Command{
		Use:   "up",
		Args:  cobra.MaximumNArgs(0),
		Short: "Provision AWS infrastructure and k3s cluster",
		Run: func(cmd *cobra.Command, args []string) {
			infra.Up(o.InstanceType)
		},
	}
)

func init() {
	upCmd.Flags().StringVar(&o.InstanceType, "instance-type", "t2.micro", "ec2 instance type to use")

	rootCmd.AddCommand(upCmd)
}
