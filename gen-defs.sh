#!/bin/bash
set -euo pipefail
cd "$(dirname "$0")"

build() {
    $1/pkg/definitions/build.sh
}

./base/pkg/iam/definitions/build.sh
build compiler
build filter
build hound
build seer
build slideshow
build sources
