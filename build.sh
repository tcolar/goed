#!/bin/bash

# Bundle file resources into the goed binary and build
(cd core && date +%s > res/resources_version.txt)
(cd core && go-bindata -pkg core -o resources_gen.go res/...)
go build 
go build api/goed_api/
go build actions/goed_builtin/
go install ./...
