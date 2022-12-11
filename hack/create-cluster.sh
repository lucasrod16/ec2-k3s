#!/bin/bash

public_ip="$(aws ec2 describe-instances --filters Name=instance-state-name,Values=running | jq -r '.Reservations[].Instances[].PublicIpAddress')"

create_cluster="k3d cluster create \
                    --servers=1 \
                    --agents=3 \
                    --k3s-arg --disable=traefik@server:0 \
                    --k3s-arg --disable=metrics-server@server:0 \
                    --k3s-arg --tls-san=${public_ip}@server:0 \
                    --port 80:80@loadbalancer \
                    --port 443:443@loadbalancer \
                    --api-port 6443"


echo "Creating the k3d cluster on the ec2 instance..."
ssh -i ~/.ssh/id_rsa -o StrictHostKeyChecking=no -o IdentitiesOnly=yes ubuntu@"${public_ip}" "${create_cluster}"