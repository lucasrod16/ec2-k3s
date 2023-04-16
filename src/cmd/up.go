package cmd

import (
	"log"
	"os"

	"github.com/lucasrod16/ec2-k3s/src/internal/infra"
	"github.com/lucasrod16/ec2-k3s/src/internal/types"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var configFile = types.ConfigFile{}

// upCmd represents the up command
var (
	upCmd = &cobra.Command{
		Use:   "up",
		Args:  cobra.MaximumNArgs(0),
		Short: "Provision AWS infrastructure and k3s cluster",
		Run: func(cmd *cobra.Command, args []string) {
			readConfigFile()
			validateConfigFile()
			infra.Up(configFile.Region, configFile.InstanceType)
		},
	}
)

func readConfigFile() {
	_, err := os.Stat(configFilePath)
	if os.IsNotExist(err) {
		log.Fatalf("File path %s does not exist\n", configFilePath)
	}

	configBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(configBytes, &configFile)
	if err != nil {
		log.Fatal(err)
	}
}

func validateConfigFile() {
	if configFile.Region == "" {
		log.Fatal("Region must be set")
	}

	if configFile.InstanceType == "" {
		log.Fatal("Instance type must be set")
	}
}

func init() {
	rootCmd.AddCommand(upCmd)
}
