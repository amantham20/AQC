#!/bin/bash

platforms=("linux/amd64" "linux/arm64" "darwin/amd64" "darwin/arm64" "windows/amd64")

mkdir -p build


for platform in "${platforms[@]}"
do
    IFS="/" read -r GOOS GOARCH <<< "$platform"
    output_name="aqc-${GOOS}-${GOARCH}"
    if [ "$GOOS" = "windows" ]; then
        output_name+='.exe'
    fi
    env GOOS=$GOOS GOARCH=$GOARCH go build -o build/$output_name .
    echo "Built: $output_name"
done
