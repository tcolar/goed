#!/bin/bash

# Bundle fonts, only if missing as not expected to change
#go-bindata -nomemcopy -pkg fonts -o ui/fonts/fonts_gen.go fonts/...

# Bundle file resources into the goed binary and build
date +%s > res/resources_version.txt
go-bindata -pkg core -o core/resources_gen.go res/...

go install ./...
