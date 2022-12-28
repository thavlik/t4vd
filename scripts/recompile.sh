#!/bin/bash
set -euo pipefail
cd "$(dirname "$0")"
cd ..
kindest build -v
kubectl rollout restart deployment -n t4vd t4vd-seer
sleep 5s
kubectl rollout restart deployment -n t4vd t4vd-compiler &
watch -n 4 kubectl get pod -n t4vd
