#!/usr/bin/env bash
dir=$(dirname "$0")

# This script generates the config-clusters.go file in ocrequest dir.
# It gets the token of sa image-pruner from every cluster to generate this config

remember-current-cluster

{
cat <<EOT
package ocrequest

// The token is from sa image-pruner of namespace cluster-tasks
var Clusters = T_ClusterConfig{
	Config: map[T_clName]T_Cluster{
EOT

for cluster in dev-scp0 dev-scp1-c1 dev-scp1-c2 cid-scp0 ppr-scp0 vpt-scp0 pro-scp0 pro-scp1; do
    # shellcheck source=/dev/null
    . ocl $cluster cluster-tasks &>/dev/null
    ocw 1>&2
    secret="$(oc -n cluster-tasks get secret|rg image-pruner|rg token|head -n 1|pc 1)"
    token="$(oc -n cluster-tasks get secret "$secret" -o jsonpath='{.data.token}'|base64 -d)"
    cat <<EOT
		"$cluster": {
			Name:          "$cluster",
			Url:           "https://api.$cluster.sf-rz.de:6443",
			Token:         "$token",
			ConfigToolUrl: "https://scpconfig-service-master.apps.$cluster.sf-rz.de"},
EOT
done

cat <<EOT
	},
}
EOT
} > "$dir/../ocrequest/config-clusters.go"

switch-back-to-current-cluster
