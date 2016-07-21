#!/bin/bash
set -ex

# Update resouyrce bundles

date +%s > res/resources_version.txt
go-bindata -pkg core -o core/resources_gen.go res/...

# Validate build 

go test ./...

go build ./...

# Create git release branch

git status
echo "Continue ?"
read should_continue

echo "Version (ie: 0.0.3)?"
read version

git checkout -b release_$version

echo -e "package core\n\nconst Version = \"$version\"\n" > core/version.go

git add core/version.go core/resources_gen.go res/resources_version.txt
git commit -m "Release $version"

git push origin release_$version

echo "Pushed branch release_$version to github, merge ready !"

# cross compile binaries

gox -osarch="darwin/amd64 darwin/386 linux/amd64 linux/386 linux/arm"\
 -output="/tmp/goed/${version}/{{.OS}}/{{.Arch}}/goed" 

# Publish to bintray
echo "publish to bintray ?"
read should_publish

curl -T /tmp/goed/$version/linux/amd64/goed -utcolar:$BINTRAY_KEY\
 https://api.bintray.com/content/tcolar/Goed/Goed/$version/linux_amd64/goed
curl -T /tmp/goed/$version/linux/386/goed -utcolar:$BINTRAY_KEY\
 https://api.bintray.com/content/tcolar/Goed/Goed/$version/linux_386/goed
curl -T /tmp/goed/$version/linux/arm/goed -utcolar:$BINTRAY_KEY\
 https://api.bintray.com/content/tcolar/Goed/Goed/$version/linux_arm/goed
curl -T /tmp/goed/$version/darwin/amd64/goed -utcolar:$BINTRAY_KEY\
 https://api.bintray.com/content/tcolar/Goed/Goed/$version/darwin_amd64/goed
curl -T /tmp/goed/$version/darwin/386/goed -utcolar:$BINTRAY_KEY\
 https://api.bintray.com/content/tcolar/Goed/Goed/$version/darwin_386/goed
 
curl -X POST -utcolar:$BINTRAY_KEY\
 https://api.bintray.com/content/tcolar/Goed/Goed/$version/publish 