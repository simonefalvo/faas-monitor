#!/bin/bash

# to let environment vars setting to take effect invoke this script with the `source` command
PROMETHEUS_URL="http://$(minikube ip):30007"
export PROMETHEUS_URL
echo "set PROMETHEUS_URL=$PROMETHEUS_URL"

# deploy test functions
faas-cli store deploy nodeinfo
echo "function nodeinfo deployed"
faas-cli store deploy figlet
echo "function figlet deployed"
faas-cli store deploy sleep
echo "function sleep deployed"

echo "test" | faas-cli invoke nodeinfo
echo "test" | faas-cli invoke figlet
echo "test" | faas-cli invoke sleep

# set proper replicas number
kubectl scale --replicas=1 -n openfaas-fn deployment nodeinfo
echo "function nodeinfo scaled to 1"
kubectl scale --replicas=0 -n openfaas-fn deployment figlet
echo "function figlet scaled to 0"
