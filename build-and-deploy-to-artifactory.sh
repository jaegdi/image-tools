#!/bin/bash
set -e

cp image-tool "old-version/image-tools_$(date +%F_%T)"
cp image-tool.exe "old-version/image-tools_$(date +%F_%T).exe"

echo Generate the config-clusters.go
scripts/generate_config.sh

echo "Build linux binary of image-tool"
go build

echo "Build windows binary of image-tool"
GOOS=windows GOARCH=amd64 go build

# aritfactory-upload.sh is a tool in the admintools repo
if ./image-tool -family=ebs -used > /dev/null; then
    echo "Push to artifactory"
    artifactory-upload.sh  -lf=image-tool       -tr=scptools-bin-develop   -tf=tools/image-tools/image-tools-linux
    artifactory-upload.sh  -lf=image-tool       -tr=scpas-bin-develop      -tf=istag_and_image_management/image-tools-linux
    artifactory-upload.sh  -lf=image-tool.exe   -tr=scptools-bin-develop   -tf=tools/image-tools/image-tools-windows
    artifactory-upload.sh  -lf=image-tool.exe   -tr=scpas-bin-develop      -tf=istag_and_image_management/image-tools-windows

    echo "Copy it to share folder PEWI4124://Daten"
    cp image-tools image-tools.exe  /gast-drive-d/Daten/
fi