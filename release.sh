#!/bin/bash

./build.sh

echo "Version?"
read version
gox -osarch="darwin/amd64 darwin/386 linux/amd64 linux/386 linux/arm" -output="bin/${version}/{{.OS}}/{{.Arch}}/goed" 
