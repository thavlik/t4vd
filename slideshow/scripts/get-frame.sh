#!/bin/bash
set -euo pipefail
cd "$(dirname "$0")"
cd ../cmd
go build -o slideshow
echo "generating frame..."
export AWS_REGION="us-west-2"
export AWS_PROFILE="personal"
export S3_ENDPOINT="https://nyc3.digitaloceanspaces.com"
export LOG_LEVEL="debug"
./slideshow get-frame
