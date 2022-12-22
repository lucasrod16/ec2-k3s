#!/bin/bash

public_ip="$(aws ec2 describe-instances --filters Name=instance-state-name,Values=running | jq -r '.Reservations[].Instances[].PublicIpAddress')"

echo "Copying the kubeconfig from the ec2 instance..."
ssh -o StrictHostKeyChecking=no -o IdentitiesOnly=yes ubuntu@"${public_ip}" 'kubectl config view --minify=true --raw=true > ./kubeconfig.yaml'
scp -o StrictHostKeyChecking=no -o IdentitiesOnly=yes ubuntu@"${public_ip}":./kubeconfig.yaml .

echo "Editing the cluster api-server IP address in the kubeconfig..."
gsed -i "s|0\.0\.0\.0|${public_ip}|g" ./kubeconfig.yaml

echo "To access your remote k3d cluster you can set your kubeconfig environment variable: export KUBECONFIG=./kubeconfig.yaml"