# Provision a k3s Cluster in AWS

`ec2-k3s` can be used to:

- Provision AWS infrastructure
  - ec2 instance
  - security group
  - ssh keypair

- Create [k3s](https://docs.k3s.io/) cluster on the ec2 instance

## Dependencies

The following dependencies are required to run this:

- [go](https://go.dev/doc/install) `1.19`

- [pulumi cli](https://www.pulumi.com/docs/get-started/install/)
  - [Logged in to a state backend](https://www.pulumi.com/docs/intro/concepts/state/#logging-into-and-out-of-state-backends)

- [aws cli](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html)
  - [Configured](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-quickstart.html) to interact with an AWS account

- [k3sup](https://github.com/alexellis/k3sup)

- ssh keypair at `~/.ssh/id_rsa` and `~/.ssh/id_rsa.pub`
  - Can be generated by using `ssh-keygen` and following the prompts

## Technical Specs

- `Ubuntu 22.04 LTS` AMI for the ec2 instance

- `t3.2xlarge` ec2 instance type

- Security group

  - Ingress rules
  
    - All ports and protocols allowed from workstation IP address only

  - Egress rules

    - All ports and protocols allowed to any IP address

- SSH keypair uses local SSH public key at `~/.ssh/id_rsa.pub`

## Usage

Clone the repository

```bash
git clone https://github.com/lucasrod16/ec2-k3s.git
```

List Makefile targets

```bash
make help
```

Build from source

```bash
make build
```

Create AWS infrastructure and k3s cluster

```bash
make up
```

Tear down AWS infrastructure

```bash
make down
```
