#!/bin/bash
set -e

[ -f image-tool ] && cp image-tool "old-version/image-tool_$(date +%F_%T)"
[ -f image-tool.exe ] && cp image-tool.exe "old-version/image-tool.exe_$(date +%F_%T).exe"

# echo Generate the config-clusters.go
# build/scripts/generate_config.sh

echo "Generate swagger doc"
swag i
echo "Build linux binary of image-tool"
go build -v

echo "Build windows binary of image-tool"
GOOS=windows GOARCH=amd64 go build -v

if ./image-tool -family=ebs -used > /dev/null; then
    echo "Push to artifactory"

    artifactory-upload.sh -lf=image-tool       -tr=scptools-bin-dev-local   -tf=/tools/image-tools/image-tools-linux/
    artifactory-upload.sh -lf=image-tool       -tr=scpas-bin-dev-local      -tf=/istag_and_image_management/image-tools-linux/

    artifactory-upload.sh -lf=image-tool.exe   -tr=scptools-bin-dev-local   -tf=/tools/image-tools/image-tools-windows/
    artifactory-upload.sh -lf=image-tool.exe   -tr=scpas-bin-dev-local      -tf=/istag_and_image_management/image-tools-windows/

    # jf rt u --server-id default --flat image-tool  /scptools-bin-develop/tools/image-tools/image-tools-linux/image-tool
    # jf rt u --server-id default --flat image-tool  /scpas-bin-develop/istag_and_image_management/image-tools-linux/image-tool
    # jf rt u --server-id default --flat image-tool.exe  /scptools-bin-develop/tools/image-tools/image-tools-windows/image-tool.exe
    # jf rt u --server-id default --flat image-tool.exe  /scpas-bin-develop/istag_and_image_management/image-tools-windows/image-tool.exe

    echo "Copy it to share folder PEWI4124://Daten"
    cp image-tool image-tool.exe  /gast-drive-d/Daten/
fi