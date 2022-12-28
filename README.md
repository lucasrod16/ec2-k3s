# Kubernetes Development Environment

This is a [Pulumi](https://www.pulumi.com/) project that does the following:

- Provisions AWS infrastructure
  - ec2 instance
  - security group
  - ssh keypair

- Creates [k3s](https://docs.k3s.io/) cluster on the ec2 instance

## Dependencies

The following dependencies are required to run this:

- go `1.19`

- [pulumi cli](https://www.pulumi.com/docs/get-started/install/) -- `brew install pulumi`

- [aws cli](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html) -- this must be [configured](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-quickstart.html) to interact with an AWS account

- [k3sup](https://github.com/alexellis/k3sup)

- ssh keypair at `~/.ssh/id_rsa` and `~/.ssh/id_rsa.pub` -- can be generated by using `ssh-keygen` and following the prompts

## Technical Specs

- Ubuntu 22.04 LTS AMI for the ec2 instance

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

Create AWS infrastructure and k3s cluster

```bash
make all
```

Tear down AWS infrastructure

```bash
make destroy
```
