#!/usr/bin/env bash
set -eo pipefail

scriptdir=$(dirname "$0")
dir=$(dirname "$scriptdir")
echo "Dir: $dir"
cd $dir

if go build -v && podman build . | tee build.log; then
    imagesha="$(tail -n 1 < build.log)"
    rm build.log


    podman tag "$imagesha" default-route-openshift-image-registry.apps.cid-scp0.sf-rz.de/scp-images/image-tool:latest

    for dst in dev pro;do
        echo "############   copy from cid to $dst   ############"
        copy-image.sh scl=cid-scp0 dcl=$dst-scp0 sns=scp-images dns=scp-images image=image-tool:latest;
    done

    oc project scp-operations-${CLUSTER/-scp0/}

    oc delete -f deploy-$CLUSTER-image-tool.yml
    oc apply -f deploy-$CLUSTER-image-tool.yml
else
    echo "Build failed"
    exit 1
fi