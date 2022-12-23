#!/bin/bash

echo "Goodbye infrastructure...until we meet again 🫡"
rm ./kubeconfig.yaml
pulumi destroy --yes