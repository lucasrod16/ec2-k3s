#!/bin/bash

go clean

if [ -f "./kubeconfig" ]; then
    rm ./kubeconfig
fi