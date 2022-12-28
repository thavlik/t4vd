#!/bin/bash
set -euo pipefail
cd "$(dirname "$0")"/kubectl
install() {
    sudo cp kubectl-$1 /usr/local/bin
    sudo chmod +x /usr/local/bin/kubectl-$1
}
install compiler
install gateway
install seer
install sources
