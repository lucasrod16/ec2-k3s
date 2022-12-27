#!/bin/bash

sleep 5

KUBECONFIG="$(pwd)/kubeconfig"
export KUBECONFIG

for deployment in "local-path-provisioner" "coredns" "metrics-server"
do
    kubectl rollout status deployment/"${deployment}" --namespace kube-system
done