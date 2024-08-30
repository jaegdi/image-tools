#!/usr/bin/env bash
set -o pipefail

scriptdir=$(dirname "$0")
dir=$(dirname "$scriptdir")
echo "Dir: $dir"
cd "$dir"
CLUSTER="${1:-$CLUSTER}"
ocl cid-scp0 -d
ocl
echo "CLUSTER: $CLUSTER"

if echo && echo "### start go build" && go build -v && echo "### go build ready" && \
   echo && echo "### start image build" && podman build . | tee build.log; then
    imagesha="$(tail -n 1 < build.log)"
    rm build.log

    echo "tag $imagesha to  default-route-openshift-image-registry.apps.cid-scp0.sf-rz.de/scp-images/image-tool:latest"
    podman tag "$imagesha" default-route-openshift-image-registry.apps.cid-scp0.sf-rz.de/scp-images/image-tool:latest
    podman push default-route-openshift-image-registry.apps.cid-scp0.sf-rz.de/scp-images/image-tool:latest

    for dst in pro-scp1;do
        echo '----------------------------------------------------------------------------------------------------------'
        copy-image.sh -v scl=cid-scp0 dcl=$dst sns=scp-images dns=scp-images image=image-tool:latest;
    done

    # oc project scp-operations-"${CLUSTER/-scp0/}"

    # oc delete -f deploy-"$CLUSTER"-image-tool.yml
    # oc apply -f deploy-"$CLUSTER"-image-tool.yml
else
    echo "Build failed"
    exit 1
fi