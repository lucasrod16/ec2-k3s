#!/bin/bash

rm -rf ./ec2-k3s

if [ -f "./kubeconfig" ]; then
    rm ./kubeconfig
fi