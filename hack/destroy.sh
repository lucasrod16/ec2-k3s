#!/bin/bash

echo "Goodbye infrastructure...until we meet again 🫡"
rm ./kubeconfig.yaml
go run main.go destroy