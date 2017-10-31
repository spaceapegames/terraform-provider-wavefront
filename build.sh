#!/bin/bash

VERSION=$(cat version)
echo "building terraform-provider-wavefront_${VERSION}_linux_amd64..."
env GOOS=linux GOARCH=amd64 go build -o terraform-provider-wavefront_${VERSION}_linux_amd64
echo "building terraform-provider-wavefront_${VERSION}_darwin_amd64..."
env GOOS=darwin GOARCH=amd64 go build -o terraform-provider-wavefront_${VERSION}_darwin_amd64
