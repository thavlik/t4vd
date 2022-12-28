#!/bin/bash
set -euo pipefail
cd "$(dirname "$0")"
cd ../cmd
go build -o slideshow
echo "Running frame extraction benchmark..."
export AWS_REGION="us-west-2"
export AWS_PROFILE="personal"
export S3_ENDPOINT="https://nyc3.digitaloceanspaces.com"
export LOG_LEVEL="debug"

test() {
    out=$(./slideshow get-frame --read-ahead $1)
    echo "read-ahead of $1 bytes: $out"
}

#test 0
test 1024
test 2048
test 4096
test 8192
test 34000
test 50000

