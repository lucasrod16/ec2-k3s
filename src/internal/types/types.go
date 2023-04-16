package types

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
)

type Infrastructure struct {
	Ami           *ec2.LookupAmiResult
	Keypair       *ec2.KeyPair
	SecurityGroup *ec2.SecurityGroup
	Server        *ec2.Instance
}

type ConfigFile struct {
	Region       string `json:"region" yaml:"region"`
	InstanceType string `json:"instanceType" yaml:"instanceType"`
	SSHKeyPath   string `json:"sshKeyPath" yaml:"sshKeyPath"`
}
