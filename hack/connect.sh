#!/bin/bash

instance_id="$(aws ec2 describe-instances --filters Name=instance-state-name,Values=running | jq -r '.Reservations[].Instances[].InstanceId')"
public_ip="$(aws ec2 describe-instances --filters Name=instance-state-name,Values=running | jq -r '.Reservations[].Instances[].PublicIpAddress')"

aws ec2 wait instance-running --output json --no-cli-pager --instance-ids "${instance_id}" &> /dev/null
echo "Giving the ec2 instance some time to boot up...🤖"
echo "This is a great time to grab a snack...🍿😋"
sleep 120
ssh -i ~/.ssh/id_rsa -o StrictHostKeyChecking=no -o IdentitiesOnly=yes ubuntu@"${public_ip}"