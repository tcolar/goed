#!/bin/bash
# Bundle file resources into the goed binary.
date +%s > res/resources_version.txt
go-bindata -pkg core -o resources_gen.go res/...
