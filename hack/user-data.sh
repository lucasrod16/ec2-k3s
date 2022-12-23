#!/bin/bash

# Install docker, k3d, and kubectl as user data script during launch of ec2 instance
apt-get update
apt-get install -y \
        ca-certificates \
        curl \
        gnupg \
        lsb-release

mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg

echo \
"deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
$(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null

apt-get update
apt-get install -y \
        docker-ce \
        docker-ce-cli \
        containerd.io

groupadd docker
usermod -aG docker ubuntu

apt-get update

curl -fsSLo /etc/apt/keyrings/kubernetes-archive-keyring.gpg https://packages.cloud.google.com/apt/doc/apt-key.gpg

echo \
    "deb [signed-by=/etc/apt/keyrings/kubernetes-archive-keyring.gpg] https://apt.kubernetes.io/ kubernetes-xenial main" \
    | sudo tee /etc/apt/sources.list.d/kubernetes.list

apt-get update
apt-get install -y kubectl

curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | bash