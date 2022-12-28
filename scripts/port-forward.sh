#!/bin/bash
set -euo pipefail
while true
do
    pod=$(kubectl get pods --selector=app=$1 -o jsonpath='{.items[*].metadata.name}')
    kubectl port-forward $pod $2:80
done
