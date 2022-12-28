#!/bin/bash
set -euo pipefail
cd "$(dirname "$0")"
cd ..
kindest build -v
kubectl rollout restart deployment -n bjjv bjjv-seer
sleep 5s
kubectl rollout restart deployment -n bjjv bjjv-compiler &
watch -n 4 kubectl get pod -n bjjv
