package ocrequest

import (
	"flag"
	"fmt"
	"os"
)

var CmdParams T_flags

func cmdUsage() {
	usageText := `

DESCRIPTION

  istag-mgmt reports image date for a application family 
  (eg. pkp, fpc, aps, ssp)

  - For existing Images it operates cluster and family specific. 

  	That means it works for one cluster like 
        'cid, int, ppr, vpt or pro' 

	and for families like 
		'pkp, sps, fpc, aps, ...'

	The cluster must be defined by the mandatory parameter 
		'-cluster=[cid|int|ppr|vpt|pro]'

	The family must be defined by the mandatory parameter 
		'-family=[aps|fpc|pkp|ssp]
		
  - For used images it looks in all clusters and reports the istags used 
	  by any deploymentconfig, job, cronjob or pod of all namespaces that 
	  belong to the application family.

  - Generate reports about imagestreamtags, imagestreams and images. 
	  The content of the report can be defined by the mandatory parameter 
	  '-output=[is|istag|image|used|all]'.
	
  - Variable output format: json(default), yaml, csv, csvfile, table and tabgroup 
	  (table with grouped rows for identical content). 
	  Output as table or tabgroup is best used when piped into less

  - filter data for reports. 
  	By specifying one of the parameters 
  		-isname=..., -istagname=..., tagname=... or -shaname=...
  	the report is filtered.

  - delete istags based on filters (not yet implemented)
	  The idea is to delete istags by filterpatterns like 
	  'older than n days' and/or 'istag name like pattern' 
	  (not yet implemented)

	For this reports the data is collected from the openshift cluster defined by 
	parameter '-cluster=...' and the 
	parameter 'family=...'. #
	For type '-used' (also included in type '-all') the data is collected
	from all clusters.

EXAMPLES

	Report all information for family pkp in cluster cid as json
	(which is the default output format)

		./report-istags -cluster=cid -family=pkp -all
		
	Report only used istags for family pkp as pretty printed table 
	(the output is paginated to fit your screen size and piped to 
		the pager define in the environment variable $PAGER/%PAGER%. 
		If $PAGER is not set, it try to use 'more')

		./report-istags -cluster=cid -family=pkp -used -table
		
	Report istags for family aps in cluster int as yaml report

		./report-istags -cluster=int -family=aps -istag -yaml

-----------------------------------------------------------------------------------------------------
`

	fmt.Printf("Usage: %s [OPTIONS] argument ...\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Println(usageText)
}

// FamilyNamespaces T_famNs

// EvalFlags evaluate all command line flags and set a struct with their values
func EvalFlags() {
	flag.Usage = cmdUsage
	// Global Flags
	clusterPtr := flag.String("cluster", "", "shortname of cluster, eg. cid,int, ppr or pro")
	tokenPtr := flag.String("token", "", "token for cluster, its a alternative to login before exec")
	familyPtr := flag.String("family", "", "family name, eg. pkp, aps, ssp or fpc ")
	namespacePtr := flag.String("namespace", "", "namespace to look for istags")

	// Collect flags
	isPtr := flag.Bool("is", false, "collect and report data of of existing imageStreams in the cluster")
	istagPtr := flag.Bool("istag", false, "collect and report data of of existing imageStreamTags in the cluster")
	shaPtr := flag.Bool("image", false, "collect and report data of existing Image's in the cluster")
	usedPtr := flag.Bool("used", false, "collect and report data of used imageStreamTags from all clusters")
	allPtr := flag.Bool("all", false, "collect and report data all imageStreams, imageStreamTags, Image's and used imageStreamTags")

	// Output format of result data
	jsonPtr := flag.Bool("json", false, "defines JSON as the output format for the reported data. This is the DEFAULT")
	yamlPtr := flag.Bool("yaml", false, "defines YAML as the output format for the reported data")
	csvPtr := flag.Bool("csv", false, "defines CSV as the output format for the reported data")
	csvFilePtr := flag.String("csvfile", "", "defines CSV as the output format for the reported data and write the types to seperate csv files <name>-typ.csv")
	htmlPtr := flag.Bool("html", false, "defines HTML as the output format for the reported data")
	tablePtr := flag.Bool("table", false, "defines formated ASCI TABLE as the output format for the reported data")
	tabgroupPtr := flag.Bool("tabgroup", false, "defines formated ASCII TABLE WITH GROUPED ROWS as the output format for the reported data")

	// Filter flags
	isnamePtr := flag.String("isname", "", "filter output of one imageStream as json, eg. -is=wvv-service")
	istagnamePtr := flag.String("istagname", "", "filter output of one imageStreamTag")
	tagnamePtr := flag.String("tagname", "", "filter output all istags with this Tag")
	shanamePtr := flag.String("shaname", "", "filter output of a Image with this SHA")

	// Options
	ocClientPtr := flag.Bool("occlient", false, "use oc client instead of api call for cluster communication")
	noProxyPtr := flag.Bool("noproxy", false, "disable use of proxy for API http requests")
	socks5ProxyPtr := flag.String("sock5", "", "set socks5 proxy url and use it for API calls")
	profilerPtr := flag.Bool("profiler", false, "enable profiler support for debugging, http://localhost:6060/debug/pprof")

	flag.Parse()

	// define map with all flags
	flags := T_flags{
		Cluster:  string(*clusterPtr),
		Token:    string(*tokenPtr),
		Family:   string(*familyPtr),
		Json:     bool(*jsonPtr) || !(bool(*yamlPtr) || bool(*csvPtr) || string(*csvFilePtr) != "" || bool(*tablePtr) || bool(*tabgroupPtr)),
		Yaml:     bool(*yamlPtr) && !bool(*jsonPtr),
		Csv:      (bool(*csvPtr) || (string(*csvFilePtr) != "")) && !bool(*jsonPtr),
		CsvFile:  string(*csvFilePtr),
		Html:     bool(*htmlPtr) && !bool(*jsonPtr),
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
			Tagname:   *tagnamePtr,
			Imagename: *shanamePtr,
			Namespace: *namespacePtr,
		},
		Options: T_flagOpts{
			OcClient:    *ocClientPtr,
			NoProxy:     *noProxyPtr,
			Socks5Proxy: *socks5ProxyPtr,
			Profiler:    *profilerPtr,
		},
	}

	InfoLogger.Println(GetJsonFromMap(flags))

	if flags.Cluster == "" {
		exitWithError("a shortname for cluster must given like: '-cluster=cid'. Is now: ", flags.Cluster)
	}
	if flags.Family == "" {
		exitWithError("a name for family must given like: '-family=pkp'")
	}
	if FamilyNamespaces[flags.Family] == nil {
		exitWithError("Family", flags.Family, "is not defined")
	}

	foundNamespace := false
	for _, v := range FamilyNamespaces[flags.Family][flags.Cluster] {
		if flags.Filter.Namespace == v {
			foundNamespace = true
		}
	}
	if !foundNamespace && !(flags.Filter.Namespace == "") {
		exitWithError("Namespace", flags.Filter.Namespace, "is no image namespace for family", flags.Family)
	}

	if !(*isPtr || *istagPtr || *shaPtr || *allPtr || *usedPtr) {
		exitWithError("As least one of the output flags must set")
	}
	CmdParams = flags
}

func FilterAllIstags(list T_ResultExistingIstagsOverAllClusters) T_ResultExistingIstagsOverAllClusters {
	outputflags := CmdParams.Output
	if !outputflags.All {
		for _, cluster := range Clusters.Stages {
			x := list[cluster]
			if !outputflags.Is {
				x.Is = nil
			}
			if !outputflags.Istag {
				x.Istag = nil
			}
			if !outputflags.Image {
				x.Image = nil
			}
			list[cluster] = x
		}
	}
	return list
}
