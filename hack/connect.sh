#!/bin/bash

PUBLIC_IP="$(aws ec2 describe-instances --filters Name=instance-state-name,Values=running | jq -r '.Reservations[].Instances[].PublicIpAddress')"

ssh -o StrictHostKeyChecking=no -o IdentitiesOnly=yes ubuntu@"${PUBLIC_IP}"