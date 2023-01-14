# Provision a k3s Cluster in AWS

`ec2-k3s` can be used to:

- Provision AWS infrastructure
  - ec2 instance
  - security group
  - ssh keypair

- Create [k3s](https://docs.k3s.io/) cluster on the ec2 instance

## Prerequisites

- [go](https://go.dev/doc/install) `1.19`

- [pulumi cli](https://www.pulumi.com/docs/get-started/install/)
  - [Logged in to a state backend](https://www.pulumi.com/docs/intro/concepts/state/#logging-into-and-out-of-state-backends)

- [aws cli](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html)
  - [Configured](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-quickstart.html) to interact with an AWS account

- [k3sup](https://github.com/alexellis/k3sup)

- ssh keypair at `~/.ssh/id_rsa` and `~/.ssh/id_rsa.pub`
  - Can be generated by using `ssh-keygen` and following the prompts

## Default Configuration

- Uses the default region of the AWS profile that is used

- `Ubuntu 22.04 LTS` AMI for the ec2 instance

- Security group

  - Ingress rules
  
    - All ports and protocols allowed from workstation IP address only

  - Egress rules

    - All ports and protocols allowed to any IP address

- SSH keypair uses local SSH public key at `~/.ssh/id_rsa.pub`

## Usage

Clone the repository and change directories into it

```bash
git clone https://github.com/lucasrod16/ec2-k3s.git

cd ec2-k3s
```

List Makefile targets

```bash
make help
```

Build from source

The build step outputs the binary in the current working directory

```bash
make build
```

It can be executed by referencing it as `./ec2-k3s`

Provision a k3s cluster in AWS

```bash
./ec2-k3s up --region <your-aws-region>
```

By default, `ec2-k3s up` will use an ec2 instance type of `t2.micro`

To view a full list of options, see [here](https://aws.amazon.com/ec2/instance-types/)

To specify the ec2 instance type to use:

```bash
./ec2-k3s up --instance-type <instance-type>
```

Teardown AWS infrastructure and k3s cluster

```bash
./ec2-k3s down
```
