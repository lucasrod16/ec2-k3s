/*
Copyright © 2022 Lucas Rodriguez
*/
package cmd

import (
	"github.com/lucasrod16/ec2-k3s/src/internal/infra"
	"github.com/lucasrod16/ec2-k3s/src/internal/types"
	"github.com/spf13/cobra"
)

var io = &types.InstanceOptions{}
var co = &types.ConfigOptions{}

// upCmd represents the up command
var (
	upCmd = &cobra.Command{
		Use:   "up",
		Args:  cobra.MaximumNArgs(0),
		Short: "Provision AWS infrastructure and k3s cluster",
		Run: func(cmd *cobra.Command, args []string) {
			infra.Up(co.Region, io.InstanceType)
		},
	}
)

func init() {
	upCmd.Flags().StringVar(&io.InstanceType, "instance-type", "t2.micro", "ec2 instance type to use")
	upCmd.Flags().StringVar(&co.Region, "region", "", "AWS region to deploy to")

	upCmd.MarkFlagRequired("region")

	rootCmd.AddCommand(upCmd)
}
