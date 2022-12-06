# Kubernetes Development Environment

This is a [Pulumi](https://www.pulumi.com/) project that does the following:
- Provisions AWS infrastructure
    - ec2 instance
    - security group
    - ssh keypair

- Creates [k3d](https://k3d.io/) cluster on the ec2 instance
- Establishes an SSH connection to the ec2 instance, ready to run `kubectl` commands against

## Dependencies

The following dependencies are required to run this:
- git

- make

- [pulumi cli](https://www.pulumi.com/docs/get-started/install/) -- you will be prompted to create a Pulumi account

- [aws cli](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html) -- this must be [configured](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-quickstart.html) to interact with an AWS account

- [jq](https://stedolan.github.io/jq/download/)

- private ssh key at `~/.ssh/id_rsa` -- generate with `ssh-keygen` and follow prompts
    - You'll need to add your public ssh key to `internal/ec2/ec2.go`
    - You can view your public key by running:
    `cat ~/.ssh/id_rsa.pub`
    ```go
    // Create SSH keypair
		keypair, err := ec2.NewKeyPair(ctx, "lucas-dev-ssh", &ec2.KeyPairArgs{
			KeyName:   pulumi.String("ec2-k3d-keypair"),
			PublicKey: pulumi.String("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDPgLhdRFNkCK/CaRZI9B1EFPf5Ax1NOvhmBN6rKXUNextSIk3t+ZDyN4iv19aUZcr3IhG/8I9AIGV1+n48aZgDyuPh9MgvVeZRXTOpUp15m80RXcTrFUP8ubTESh8BiYee4DfmUcfccXjyB00OT5GK0OXNiWIGPkElpHPnwmRnRQ6bHyx8HJVMKC0MVwZe+RgydylDasUGJm+gE4+4xc+7F587mT+R17IjS4MZkIkIwIApez+euDp8lqtRuH3AGYpQdxkz09WuSRUBwgWOf4FkpB5+NZtDO3of22QJvL/6PZ2numx/llNhTO6ya1VrWpPH4q3ghaxjZy+v/Mh1+QrUx8r1RF6GSs18iKaBFHQin/it1KDSLW6wbHMQgsU9JrolWT93bZkOaahBDYkPvubgnGBEZ9kDDTVVzowUUJ6QNu932JJJk98dp0Q346RumUhAcgVgjenTPwgs6DgMc2q21pY/96vVGyJ87BlhJpkqZlVQMHfqnmCC/1hd8GWFmKU= lucas@Lucass-MacBook-Pro.local"),
    ```

## Usage

Clone the repository
```bash
git clone https://github.com/lucasrod16/ec2-k3d.git
```

List Makefile targets
```bash
make help
```

Create AWS infrastructure, k3d cluster, and ssh to ec2 instance
```bash
make all
```

Tear down AWS infrastructure
```bash
make destroy
```