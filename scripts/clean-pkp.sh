#!/usr/bin/env bash
scriptdir="$(dirname $0)"
logpostfix="log-$CLUSTER-$(date +%F_%T).log"
{
    image-tools -family=pkp -cluster=$CLUSTER -delete -snapshot
    image-tools -family=pkp -cluster=$CLUSTER -delete -isname=hybris -minage=365
    image-tools -family=pkp -cluster=$CLUSTER -delete -isname=db2 -minage=200
    image-tools -family=pkp -cluster=$CLUSTER -delete  -minage=365
    image-tools -family=pkp -cluster=$CLUSTER -delete -tagname=pkp-3.1.[0-9]+-build-*
    image-tools -family=pkp -cluster=$CLUSTER -delete -tagname=pkp-3.6.*-build-*
    image-tools -family=pkp -cluster=$CLUSTER -delete -tagname=pkp-3.8.*-build-*
    image-tools -family=pkp -cluster=$CLUSTER -delete -tagname=pkp-3.9.*-build-*
    image-tools -family=pkp -cluster=$CLUSTER -delete  -isname=priorityservice -minage=0

    image-tools -family=aps -cluster=$CLUSTER -delete -snapshot
    image-tools -family=aps -cluster=$CLUSTER -delete -tagname=CBDE1.[01].* 2>/dev/null|rg -v postgr

    image-tools -family=fpc -cluster=$CLUSTER -delete -tagname=FPC* -minage=350

    image-tools -family=vps -cluster=$CLUSTER -delete -minage=50

    cat "$scriptdir/hub-list.txt"

} 2>/dev/null | tee "delete-$logpostfix" | \
    rg -av "db2|mail|maildev" | \
    xargs -n 1 -I{} bash -c "{}"

# prune registry
prune-registry-of-current-cluster.sh