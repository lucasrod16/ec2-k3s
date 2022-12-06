package main

import "main/internal/ec2"

func main() {
	// Create security group, ec2 instance, and ssh key pair
	ec2.Create()
}
