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

type InstanceOptions struct {
	InstanceType string
}

type ConfigOptions struct {
	Region string
}
