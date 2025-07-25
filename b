#!/bin/bash
set -e

export GOOS=linux
export GOARCH=amd64

# Ensure the build directory exists
mkdir -p build

# Build the Go executable
go build -o build/snap-rq

echo "✅ Build succeeded."

# Run the executable
./build/snap-rq
