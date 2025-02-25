#!/usr/bin/env bash
set -eo pipefail

# Ensure git checkout master is executed on script exit
trap 'git checkout master' EXIT

scriptdir=$(dirname "$0")
dir=$(dirname "$scriptdir")
echo "Dir: $dir"
cd "$dir"
CLUSTER="${1:-$CLUSTER}"
ocl cid-scp0 -d > /dev/null
ocl > /dev/null
echo "CLUSTER: $CLUSTER"
tagversion=$(get-git-tag.sh)
git checkout "$tagversion"

echo "Generate swagger doc"
swag init
echo "### start go build with tag '$tagversion'"
go build -tags netgo -v
echo "### go build ready"

rm -f build.log
if echo && echo "### start image build with git tag $tagversion" && podman build . | tee build.log; then
    imagesha="$(tail -n 1 < build.log)"

    echo "tag $imagesha to  default-route-openshift-image-registry.apps.cid-scp0.sf-rz.de/scp-images/image-tool:$tagversion"
    podman tag "$imagesha" default-route-openshift-image-registry.apps.cid-scp0.sf-rz.de/scp-images/image-tool:$tagversion
    podman push default-route-openshift-image-registry.apps.cid-scp0.sf-rz.de/scp-images/image-tool:$tagversion

    for dst in pro-scp1;do  #  cid-scp0 pro-scp0
        echo '----------------------------------------------------------------------------------------------------------'
        copy-image.sh -v -sc=cid-scp0 -dc=$dst -sn=scp-images -dn=scp-images -i=image-tool:$tagversion;
        echo '----------------------------------------------------------------------------------------------------------'

        . ocl $dst

        if [[ $dst =~ pro-scp1 ]]; then
            echo '----------------------------------------------------------------------------------------------------------'
            echo "Deplopy image to registry-quay-quay.apps.pro-scp1.sf-rz.de/scp/image-tool:$tagversion"
            podman login -u "$USER" -p "$($LDAPPASSWORDPROVIDER)" registry-quay-quay.apps.pro-scp1.sf-rz.de
            podman tag "$imagesha"  registry-quay-quay.apps.pro-scp1.sf-rz.de/scp/image-tool:$tagversion
            podman push  registry-quay-quay.apps.pro-scp1.sf-rz.de/scp/image-tool:$tagversion
            echo '----------------------------------------------------------------------------------------------------------'
            oc project scp-ops-central
        else
            oc project scp-operations-"${dst/-scp[01]/}"
        fi
    done

    echo '----------------------------------------------------------------------------------------------------------'
    oc delete -f deploy-"$CLUSTER"-image-tool.yml
    echo '----------------------------------------------------------------------------------------------------------'
    cat deploy-"$CLUSTER"-image-tool.yml | \
        sed -e "s/image: registry-quay-quay.apps.pro-scp1.sf-rz.de\/scp\/image-tool.*/image: registry-quay-quay.apps.pro-scp1.sf-rz.de\/scp\/image-tool:$tagversion/" | \
            oc apply -f-
else
    echo "Build failed"
    exit 1
fi
