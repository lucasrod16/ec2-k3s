#!/bin/bash

PUBLIC_IP="$(aws ec2 describe-instances --filters Name=instance-state-name,Values=running | jq -r '.Reservations[].Instances[].PublicIpAddress')"

k3sup install \
        --ip="${PUBLIC_IP}" \
        --user ubuntu \
        --k3s-extra-args "--disable traefik"