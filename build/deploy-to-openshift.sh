#!/usr/bin/env bash
set -eo pipefail

script="$(basename "$0")"
scriptdir="$(dirname "$0")"
dir=$(dirname "$scriptdir")
echo "script: $script, scriptdir: $scriptdir, Dir: $dir"
cd "$dir"


hilfe() {
    if [[ -n $1 ]]; then
        echo
        echo '***'  "$*"  '***'
        echo
    fi
    cat <<-EOH

    SYNOPSIS

        $script [-c|--cluster <clustername>] [-b|--build] [-i|imagebuild] [-d|--deploy] [-t|--tag <tag>] [-v|--version] [-h|--help]  [--loglevel[=]<value>]

    OPTIONS

        -c cluster | --cluster=cluster[,cluster...]
            The cluster to deploy the image-tool to. Default is cid-scp0

        -b | --build
            Build the image-tool and deploy to artifactory

        -i | --imagebuild
            create docker image of the image-tool

        -d | --deploy
            Deploy the image-tool pon openshift clusters

        -t git-tag | --tag git-tag
            The git-tag to build and deploy. Default is the latest tag that is defined in the repo.
            The image and the deployment are tagged with this tag.
            The repo is checked out to this tag for build and deploy unless the tag is 'latest',
            then the build and deploy is executed on the current branch and tagged as 'latest'.

        -h | --help
            Print this help message

        -v | --version
            Print tag version of repo and exit

    DESCRIPTION

        $script builds and deploys the image-tool to the specified cluster.

EOH
}

CLUSTER=cid-scp0
optspec=":bidvhct-:"
tagversion="$(get-git-tag.sh)"
build='false'
imagebuild='false'
deploy='false'
while getopts "$optspec" optchar; do
    case "${optchar}" in
        -)
            case "${OPTARG}" in
                cluster)
                    val="${!OPTIND}"; OPTIND=$(( $OPTIND + 1 ))
                    echo "Parsing option: '--${OPTARG}', value: '${val}'" >&2;
                    CLUSTER="$val"
                    ;;
                build)
                    val="${!OPTIND}";
                    echo "Parsing option: '--${OPTARG}', value: '${val}'" >&2;
                    build='true'
                    ;;
                imagebuild)
                    val="${!OPTIND}";
                    echo "Parsing option: '--${OPTARG}', value: '${val}'" >&2;
                    imagebuild='true'
                    ;;
                deploy)
                    val="${!OPTIND}";
                    echo "Parsing option: '--${OPTARG}', value: '${val}'" >&2;
                    deploy='true'
                    ;;
                tag)
                    val="${!OPTIND}"; OPTIND=$(( $OPTIND + 1 ))
                    echo "Parsing option: '--${OPTARG}', value: '${val}'" >&2;
                    tagversion="$val"
                    ;;
                version)
                    val="${!OPTIND}";
                    echo "Parsing option: '--${OPTARG}'" >&2
                    echo "latest tag in repo: S(get-git-tag.sh)"
                    exit
                    ;;
                help)
                    hilfe ''
                    exit
                    ;;
                *)
                    if [ "$OPTERR" = 1 ] && [ "${optspec:0:1}" != ":" ]; then
                        hilfe "Unknown option --${OPTARG}"
                        exit
                    fi
                    ;;
            esac;;
        c)
            val="${!OPTIND}"; OPTIND=$(( $OPTIND + 1 ))
            echo "Parsing option: '-${optchar}', value: '${val}'" >&2
            CLUSTER=$val
            ;;
        b)
            val="${!OPTIND}";
            echo "Parsing option: '-${optchar}'" >&2
            build=true
            ;;
        i) val="${!OPTIND}";
            echo "Parsing option: '-${optchar}'" >&2
            imagebuild=true
            ;;
        d)
            val="${!OPTIND}";
            echo "Parsing option: '-${optchar}'" >&2
            deploy=true
            ;;
        t)
            val="${!OPTIND}"; OPTIND=$(( $OPTIND + 1 ))
            echo "Parsing option: '-${optchar}', value: '${val}'" >&2
            tagversion=$val
            ;;
        v)
            val="${!OPTIND}";
            echo "Parsing option: '-${optchar}'" >&2
            echo "latest tag in repo: $(get-git-tag.sh)"
            exit
            ;;
        h)
            val="${!OPTIND}";
            hilfe ''
            exit 0
            ;;
        *)
            if [ "$OPTERR" != 1 ] || [ "${optspec:0:1}" = ":" ]; then
                echo "Non-option argument: '-${OPTARG}'" >&2
            fi
            ;;
    esac
done

remember-current-cluster
ocl cid-scp0 -d > /dev/null
ocl > /dev/null
echo "working on CLUSTER: $CLUSTER"

# B U I L D   L O C A L
if [ "$build" == 'true' ]; then
    if [ "$tagversion" != 'latest' ]; then
        # Ensure git checkout master is executed on script exit
        trap 'git checkout master' EXIT
        git checkout "$tagversion"
    fi
    echo "Build image-tool local and deploy to artifactory"
    "$scriptdir"/build-and-deploy-to-artifactory.sh
fi

# B U I L D   I M A G E
if [ "$imagebuild" == 'true' ]; then
    if [ "$tagversion" != 'latest' ]; then
        # Ensure git checkout master is executed on script exit
        trap 'git checkout master' EXIT
        git checkout "$tagversion"
    fi

    echo "Generate swagger doc"
    swag init
    echo "### start go build with tag '$tagversion'"
    go build -tags netgo -v
    echo "### go build ready"

    rm -f build.log
    if echo && echo "### start image build with git tag $tagversion" && podman build . | tee build.log; then
        imagesha="$(tail -n 1 < build.log)"

        echo "tag $imagesha to  default-route-openshift-image-registry.apps.cid-scp0.sf-rz.de/scp-images/image-tool:$tagversion"
        podman tag "$imagesha" "default-route-openshift-image-registry.apps.cid-scp0.sf-rz.de/scp-images/image-tool:$tagversion"
        podman push "default-route-openshift-image-registry.apps.cid-scp0.sf-rz.de/scp-images/image-tool:$tagversion"

    else
        echo "Build failed"
        switch-back-to-current-cluster
        exit 1
    fi
fi

# D E P L O Y
if [ "$deploy" == 'true' ]; then
    for dst in pro-scp1 cid-scp0 pro-scp0;do
        . ocl "$dst"
        echo "Deploying on $CLUSTER with tag $tagversion"
        echo '----------------------------------------------------------------------------------------------------------'
        oc delete --ignore-not-found -f "$scriptdir/deploy-$dst-image-tool.yml"
        echo '----------------------------------------------------------------------------------------------------------'
        # shellcheck disable=SC2002
        cat "$scriptdir/deploy-$dst-image-tool.yml" | \
            sed -e "s/image: registry-quay-quay.apps.pro-scp1.sf-rz.de\/scp\/image-tool.*/image: registry-quay-quay.apps.pro-scp1.sf-rz.de\/scp\/image-tool:$tagversion/" | \
                oc apply -f-
        echo
    done
fi
switch-back-to-current-cluster
