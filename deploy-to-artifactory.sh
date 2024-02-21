#!/bin/bash
set -e

echo "Build linux binary of image-tools"
go build

echo "Build windows binary of image-tools"
GOOS=windows GOARCH=amd64 go build

echo "Push to artifactory"
artifactory-upload.sh  -lf=image-tools       -tr=scptools-bin-develop   -tf=tools/image-tools/image-tools-linux
artifactory-upload.sh  -lf=image-tools       -tr=scpas-bin-develop      -tf=istag_and_image_management/image-tools-linux
artifactory-upload.sh  -lf=image-tools.exe   -tr=scptools-bin-develop   -tf=tools/image-tools/image-tools-windows
artifactory-upload.sh  -lf=image-tools.exe   -tr=scpas-bin-develop      -tf=istag_and_image_management/image-tools-windows

echo "Copy it to share folder PEWI4124://Daten"
cp image-tools image-tools.exe  /gast-drive-d/Daten/
