#!/usr/bin/env bash
script="$0"
dir=$(dirname "$script")
set -Eeo pipefail
# set -x
# This script generates the config-clusters.go file in ocrequest dir.
# It gets the token of sa image-pruner from every cluster to generate this config
domain="$1"
if [ -z "$domain" ]; then
    echo "the domain for the cluster must be givan as parameter. eg. mydomain.de"
    exit 1
fi

remember-current-cluster || exit 1

{
cat <<EOT
package ocrequest

// The token is from sa image-pruner of namespace cluster-tasks
var Clusters = T_ClusterConfig{
	Config: map[T_clName]T_Cluster{
EOT
for cluster in $(cluster_list all); do
	echo >&2
	echo -n "Login into '$cluster'" >&2
    # shellcheck source=/dev/null
	timeout 5 ocl "$cluster" cluster-tasks &>/dev/null || { echo >&2;echo '------------------'  >&2; continue; }
	echo -n ", get namespace: " >&2
	ocw 1>&2 || continue
	sa_name='image-pruner'
	# echo -n "get secret $secretname" >&2
	# secret="$(oc -n cluster-tasks get secret|rg "$secretname"|rg token|head -n 1|pc 1)" || { echo ", no secret of $secretname found !" >&2;echo '------------------'  >&2; continue; }
	# echo -n ", get token" >&2
	token="$(oc -n cluster-tasks create token "$sa_name" --duration=1576800m)" || continue
	echo -n ", write config" >&2
	cat <<-EOT
	        "$cluster": {
	            Name:          "$cluster",
	            Url:           "https://api.$cluster.$domain:6443",
	            Token:         "$token",
	            ConfigToolUrl: "https://scpconfig-service-master.apps.$cluster.$domain"},
	EOT
	echo ", config written." >&2
	echo '-----------------' >&2
done

cat <<EOT
	},
}
EOT
} > "$dir/../../ocrequest/config-clusters.go"

ls -l "$dir/../../ocrequest/config-clusters.go"
switch-back-to-current-cluster

echo
echo "#########################################################################################"
echo "To use this generated copnfig, you must execute build/build-and-deploy-image.sh pro-scp1"
echo "#########################################################################################"

