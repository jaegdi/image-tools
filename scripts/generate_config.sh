#!/usr/bin/env bash
script="$0"
dir=$(dirname "$script")
set -Eeo pipefail
# set -x
# This script generates the config-clusters.go file in ocrequest dir.
# It gets the token of sa image-pruner from every cluster to generate this config

remember-current-cluster || exit 1

{
cat <<EOT
package ocrequest

// The token is from sa image-pruner of namespace cluster-tasks
var Clusters = T_ClusterConfig{
	Config: map[T_clName]T_Cluster{
EOT
#                                   dev-scp1-c2
for cluster in dev-scp0 dev-scp1-c1             cid-scp0 ppr-scp0 vpt-scp0 pro-scp0 pro-scp1; do
	echo >&2
	CLUSTER="$cluster"
	echo -n "Login into $cluster" >&2
    # shellcheck source=/dev/null
	timeout 3 ocl $CLUSTER cluster-tasks &>/dev/null || { echo >&2;echo '------------------'  >&2; continue; }
	echo -n ", get namespace: " >&2
	ocw 1>&2 || continue
	echo -n "get secret" >&2
	secret="$(oc -n cluster-tasks get secret|rg image-pruner|rg token|head -n 1|pc 1)" || continue
	echo -n ", get token" >&2
	token="$(oc -n cluster-tasks get secret "$secret" -o jsonpath='{.data.token}'|base64 -d)" || continue
	echo -n ", write config" >&2
	cat <<EOT
		"$cluster": {
			Name:          "$cluster",
			Url:           "https://api.$cluster.sf-rz.de:6443",
			Token:         "$token",
			ConfigToolUrl: "https://scpconfig-service-master.apps.$cluster.sf-rz.de"},
EOT
	echo ", config written." >&2
	echo '-----------------' >&2
done

cat <<EOT
	},
}
EOT
} > "$dir/../ocrequest/config-clusters.go"

switch-back-to-current-cluster