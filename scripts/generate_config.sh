#!/usr/bin/env bash
dir=$(dirname "$0")

{
cat <<EOT
package ocrequest

// The token is from sa image-pruner of namespace cluster-tasks
var Clusters = T_ClusterConfig{
	Config: map[T_clName]T_Cluster{
EOT

for cluster in dev-scp0 cid-scp0 ppr-scp0 vpt-scp0 pro-scp0; do
    # shellcheck source=/dev/null
    . ocl $cluster cluster-tasks &>/dev/null
    ocw 1>&2
    secret="$(oc get secret|rg image-pruner|rg token|head -n 1|pc 1)"
    token="$(oc get secret "$secret" -o jsonpath='{.data.token}'|base64 -d)"
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