package ocrequest

import (
	"flag"
)

var CmdParams T_flags

// FamilyNamespaces T_famNs
func EvalFlags() {
	// Global Flags
	clusterPtr := flag.String("cluster", "", "shortname of cluster, eg. cid,int, ppr or pro")
	tokenPtr := flag.String("token", "", "token for cluster, its a alternative to login before exec")
	familyPtr := flag.String("family", "", "family name, eg. pkp, aps, ssp or fpc ")
	namespacePtr := flag.String("namespace", "", "namespace to look for istags")

	// Output flags
	isPtr := flag.Bool("is", false, "output of imageStreams as json")
	istagPtr := flag.Bool("istag", false, "output of imageStreamTags as json")
	shaPtr := flag.Bool("sha", false, "output of Sha's as json")
	usedPtr := flag.Bool("used", false, "output used imageStreams imageStreamTags and Sha's as json")
	allPtr := flag.Bool("all", false, "output all imageStreams imageStreamTags and Sha's as json")
	ocClientPtr := flag.Bool("occlient", false, "use oc client instead of api call for cluster communication")

	// Filter flags
	isnamePtr := flag.String("isname", "", "filter output of one imageStream as json, eg. -is=wvv-service")
	istagnamePtr := flag.String("istagname", "", "filter output of one imageStreamTag as json")
	shanamePtr := flag.String("shaname", "", "filter output of one Sha as json")
	flag.Parse()

	// define map with all flags
	flags := T_flags{
		Cluster:  string(*clusterPtr),
		Token:    string(*tokenPtr),
		Family:   string(*familyPtr),
		OcClient: bool(*ocClientPtr),
		// "namespace": *namespacePtr,
		Output: T_flagOut{
			Is:    *isPtr,
			Istag: *istagPtr,
			Sha:   *shaPtr,
			Used:  *usedPtr,
			All:   *allPtr,
		},
		Filter: T_flagFilt{
			Isname:    *isnamePtr,
			Istagname: *istagnamePtr,
			Shaname:   *shanamePtr,
			Namespace: *namespacePtr,
		},
	}

	InfoLogger.Println(GetJsonFromMap(flags))

	if flags.Cluster == "" {
		exitWithError("a shortname for cluster must given like: '-cluster=cid'. Is now: " + flags.Cluster)
	}
	if flags.Family == "" {
		exitWithError("a name for family must given like: '-family=pkp'")
	}
	if FamilyNamespaces[flags.Family] == nil {
		exitWithError("Family " + flags.Family + " is not defined")
	}

	foundNamespace := false
	for _, v := range FamilyNamespaces[flags.Family] {
		if flags.Filter.Namespace == v {
			foundNamespace = true
		}
	}
	if !foundNamespace && !(flags.Filter.Namespace == "") {
		exitWithError("Namespace " + flags.Filter.Namespace +
			" is no image namespace for family " + flags.Family)
	}

	if !(*isPtr || *istagPtr || *shaPtr || *allPtr || *usedPtr) {
		exitWithError("As least one of the output flags mus set")
	}
	CmdParams = flags
}

func FilterAllIstags(list T_result) T_result {
	outputflags := CmdParams.Output
	if !outputflags.All {
		if !outputflags.Is {
			list.Is = nil
		}
		if !outputflags.Istag {
			list.Istag = nil
		}
		if !outputflags.Sha {
			list.Sha = nil
		}
	}
	return list
}
