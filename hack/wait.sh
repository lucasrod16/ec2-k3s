#!/bin/bash

instance_id="$(aws ec2 describe-instances --filters Name=instance-state-name,Values=running | jq -r '.Reservations[].Instances[].InstanceId')"

aws ec2 wait instance-running --output json --no-cli-pager --instance-ids "${instance_id}" &> /dev/null
until [ "$(aws ec2 describe-instance-status --instance-ids "${instance_id}" | jq -r '.InstanceStatuses[].InstanceStatus.Details[].Status')" == "passed" ]
do
    echo "Waiting for the ec2 instance to be ready...🤖"
    sleep 3
done

echo "The instance is ready!"