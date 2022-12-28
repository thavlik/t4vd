#!/bin/bash
kubectl get secret -n t4vd $1 -o json \
    | jq .data.$2 \
    | xargs echo \
    | base64 -d
