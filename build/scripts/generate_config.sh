#!/usr/bin/env bash
script="$0"
dir=$(dirname "$script")
set -Eeo pipefail
# This script gets image-pruner tokens and writes the runtime cluster configuration.
domain="$1"
if [ -z "$domain" ]; then
    echo "the domain for the cluster must be givan as parameter. eg. mydomain.de"
    exit 1
fi

command -v jq >/dev/null || { echo "jq is required" >&2; exit 1; }
remember-current-cluster || exit 1

umask 077
entries=$(mktemp "$dir/../../.cluster-entries.XXXXXX")
output=$(mktemp "$dir/../../.cluster-config.XXXXXX")
cleanup() {
    rm -f "$entries" "$output"
    switch-back-to-current-cluster
}
trap cleanup EXIT

for cluster in $(cluster_list all); do
	echo >&2
	echo -n "Login into '$cluster'" >&2
    # shellcheck source=/dev/null
	timeout 5 ocl "$cluster" cluster-tasks &>/dev/null || { echo >&2;echo '------------------'  >&2; continue; }
	echo -n ", get namespace: " >&2
	ocw 1>&2 || continue
	secretname='image-pruner'
	echo -n "get secret $secretname" >&2
	secret="$(oc -n cluster-tasks get secret|rg "$secretname"|rg token|head -n 1|pc 1)" || { echo ", no secret of $secretname found !" >&2;echo '------------------'  >&2; continue; }
	echo -n ", get token" >&2
	token="$(oc -n cluster-tasks get secret "$secret" -o jsonpath='{.data.token}'|base64 -d)" || continue
	echo -n ", write config" >&2
	jq -nc --arg name "$cluster" \
	    --arg url "https://api.$cluster.$domain:6443" \
	    --arg token "$token" \
	    --arg configtoolurl "https://crpconfig-service-master.apps.$cluster.$domain" \
	    '{name: $name, url: $url, token: $token, configtoolurl: $configtoolurl}' >> "$entries"
	echo ", config written." >&2
	echo '-----------------' >&2
done

if [ ! -s "$entries" ]; then
    echo "No cluster configuration was generated" >&2
    exit 1
fi
jq -s '{config: (map({key: .name, value: .}) | from_entries)}' "$entries" > "$output"
mv "$output" "$dir/../../clusterconfig.json"
echo "Wrote clusterconfig.json; keep this file private." >&2
