package ocrequest

import (
	"flag"
	"fmt"
	"os"
)

var CmdParams T_flags

func cmdUsage() {
	usageText := `
istag-mgmt has the following functions.

  - It operates cluster and family specific. That means it works for one cluster like
	'cid, int, ppr, vpt or pro' and for families like 'pkp, sps, fpc, aps, ...'
	The cluster must be defined by the mandatory parametter '-cluster=[cid|int|ppr|vpt|pro]'
	The family must be defined by the mandatory parameter '-family=[aps|fpc|pkp|ssp]
	

  - Generate JSON reports about imagestreamtags, imagestreams and images. The content of the JSON 
    report can be defined by the mandatory parameter '-output=[is|istag|image|used|all]'.

  - filter data for reports

  - delete istags based on filters like 'older than n days' and/or 'istag name like pattern'

For this reports the data is collected from the oc cluster defined by parameter '-cluster=...' and
the parameter 'family=...'.
`

	fmt.Printf("Usage: %s [OPTIONS] argument ...\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Println(usageText)
}

// FamilyNamespaces T_famNs
func EvalFlags() {
	flag.Usage = cmdUsage
	// Global Flags
	clusterPtr := flag.String("cluster", "", "shortname of cluster, eg. cid,int, ppr or pro")
	tokenPtr := flag.String("token", "", "token for cluster, its a alternative to login before exec")
	familyPtr := flag.String("family", "", "family name, eg. pkp, aps, ssp or fpc ")
	namespacePtr := flag.String("namespace", "", "namespace to look for istags")
	ocClientPtr := flag.Bool("occlient", false, "use oc client instead of api call for cluster communication")

	// Output format of result data
	jsonPtr := flag.Bool("json", false, "defines JSON as the output format for the reported data. This is the DEFAULT")
	yamlPtr := flag.Bool("yaml", false, "defines YAML as the output format for the reported data")
	csvPtr := flag.Bool("csv", false, "defines CSV as the output format for the reported data")
	tablePtr := flag.Bool("table", false, "defines formated ASCI table as the output format for the reported data")
	tabgroupPtr := flag.Bool("tabgroup", false, "defines formated ASCII table with grouped rows as the output format for the reported data")

	// Output flags
	isPtr := flag.Bool("is", false, "output of imageStreams as json")
	istagPtr := flag.Bool("istag", false, "output of imageStreamTags as json")
	shaPtr := flag.Bool("image", false, "output of Sha's as json")
	usedPtr := flag.Bool("used", false, "output used imageStreams imageStreamTags and Sha's as json")
	allPtr := flag.Bool("all", false, "output all imageStreams imageStreamTags and Sha's as json")

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
		Json:     bool(*jsonPtr) || !(bool(*yamlPtr) || bool(*csvPtr) || bool(*tablePtr) || bool(*tabgroupPtr)),
		Yaml:     bool(*yamlPtr) && !bool(*jsonPtr),
		Csv:      bool(*csvPtr) && !bool(*jsonPtr),
		Table:    bool(*tablePtr) && !bool(*jsonPtr),
		TabGroup: bool(*tabgroupPtr) && !bool(*jsonPtr),

		Output: T_flagOut{
			Is:    *isPtr,
			Istag: *istagPtr,
			Image: *shaPtr,
			Used:  *usedPtr,
			All:   *allPtr,
		},
		Filter: T_flagFilt{
			Isname:    *isnamePtr,
			Istagname: *istagnamePtr,
			Imagename:   *shanamePtr,
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
		if !outputflags.Image {
			list.Sha = nil
		}
	}
	return list
}
