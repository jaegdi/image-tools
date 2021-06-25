#!/bin/bash
echo "Build linux and windows binary of image-tools"
go build
GOOS=windows GOARCH=amd64 go build

echo "Push to artifactory"
curl -ujaegdi:AP5guU7a8TN9S6zAtrRadGyMBUt -T image-tools     "https://artifactory-pro.sf-rz.de:8443/artifactory/scpas-bin-dev-local/istag_and_image_management/image-tools-linux/image-tools"
curl -ujaegdi:AP5guU7a8TN9S6zAtrRadGyMBUt -T image-tools.exe "https://artifactory-pro.sf-rz.de:8443/artifactory/scpas-bin-dev-local/istag_and_image_management/image-tools-windows/image-tools.exe"

cp image-tools image-tools.exe  /gast-drive-d/Daten/