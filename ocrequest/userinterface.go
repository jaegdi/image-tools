package ocrequest

import (
	"flag"
	"fmt"
	"os"
)

var CmdParams T_flags

// cmdUsage print the man page
func cmdUsage() {
	usageText := `

DESCRIPTION

  image-tools 
        - generate reports about images, imagestream(is), imagestreamtags(istag) 
		  with information like AgeInDays, Date, Namespace, Buildtags,..
          for a application family (eg. pkp, fpc, aps, ssp)
        or
        - generate shellscript output to delete istags when parameter -delete is set.
          The -delete parameter disables the report output, instead the delete script is generated as output
	
	image-tools only read information around images from the clusters and generate output. It never change or
	delete something in the clusters. Eg. for delete istags it only generates a script output, which than can 
	be executed by a cluster admin to really delete the istags.

	image-tools always write a log to 'image-tools.log' in the current directory.
	Additional it writes the log messages also to STDERR. To disable the log output to STDERR use parameter -nolog
	The report or delete-script output is written to STDOUT.

  - For reporting existing Images and delete istags it operates cluster and family specific. 
    For reporting used images it works over all clusters but family specific.
    It never works across different families.

  	For existing images, istags or imagestreams(is) that means it works for one cluster like 
        'cid, int, ppr, vpt or pro' 

	and for one families like 
		'pkp, sps, fpc, aps, ...'

	The cluster must be defined by the mandatory parameter 
		'-cluster=[cid|int|ppr|vpt|pro]'

	The family must be defined by the mandatory parameter 
		'-family=[aps|fpc|pkp|ssp]

  - Generate reports about imagestreamtags, imagestreams and images. 
		The content of the report can be defined by the mandatory parameter 
		'-output=[is|istag|image|used|unused|all]'. 
  
  - For used images it looks in all clusters and reports the istags used 
	  by any deploymentconfig, job, cronjob or pod of all namespaces that 
	  belong to the application family.

  - Variable output format: json(default), yaml, csv, csvfile, table and tabgroup 
	  (table with grouped rows for identical content). 
	  Output as table or tabgroup is automatically piped into less (or what is defined as PAGER)

  - filter data for reports. 
  	By specifying one of the parameters 
  		-isname=..., -istagname=..., tagname=... or -shaname=...
  	the report is filtered.

  - delete istags based on filters
	  The idea is to delete istags by filterpatterns like 
	  'older than n days' and/or 'istag name like pattern'
	  The image tool didn't delete the istags directly instead
	  it generate a shell-script that can be executed by a cluster admin
	  to delete the istag, they fit to the given filter parameters

	  To switch from reporting mode to delete mode, set the praameter -delete
	  But it needs further parameters:
	  -snapshot        delete istags with snapshot or PR-nn in the tag name.
	  -nonbuild        is specific for family pkp and delete istags fo all images, if they have no build tag.
	  -minage=int      defines the minimum age (in days) for istag to delete them. Default is 60 days.
	  -delpattern=str  define a regexp pattern for istags to delete
	  and
	  -isname=str
	  -tagname=str
	  -istagname=str
	  -namespace=str
	  can also be used to filter istags to delete them.
	  See examples in the EXAMPLES section

	For this reports the data is collected from the openshift cluster defined by 
	the mandatory parameters 
	     '-cluster=...' and the 
	     'family=...'
	For type '-used' (also included in type '-all') the data is collected
	from all clusters.

	For more speed a cache is build from the first run in  /tmp/tmp-report-istags/* 
	and used if not older than 5 minutes. If the cache is older or deleted, the
	data is fresh collected from the clusters.

INSTALLATION

	image-tols is a statically linked go programm and has no runtime dependencies. No installation is
	neccessary. Copy the binaryt into a directory, which is in the search path is enough.

EXAMPLES

  REPORTING

	Report all information for family pkp in cluster cid as json
	(which is the default output format)

		image-tools -cluster=cid -family=pkp -all
	
		or as table
		image-tools -cluster=cid -family=pkp -all -table
	
		or csv in different files for each type of information
		image-tools -cluster=cid -family=pkp -all -csvfile=prefix
		writes the output to different files 'prefix-type' in current directory
		
	Report only __used__ istags for family pkp as pretty printed table 
	(the output is paginated to fit your screen size and piped to 
		the pager define in the environment variable $PAGER/%PAGER%. 
		If $PAGER is not set, it try to use 'more')

		image-tools -cluster=cid -family=pkp -used -table
		or json
		image-tools -cluster=cid -family=pkp -used
		or yaml
		image-tools -cluster=cid -family=pkp -used -yaml
		or csv
		image-tools -cluster=cid -family=pkp -used -csv
		
	Report istags for family aps in cluster int as yaml report

		image-tools -cluster=int -family=aps -istag -yaml

	Report ImageStreams for family aps in cluster int as yaml report

		image-tools -cluster=int -family=aps -is -yaml

	Report Images for family aps in cluster int as yaml report

		image-tools -cluster=int -family=aps -image -yaml

  DELETE

	Generate a shell script to delete old istags(60 days, the default) for family pkp in cluster cid
    and all old snapshot istags and nonbuild istags and all istags of header-service, footer-service and zahlungsstoerung-service

		image-tools -family=pkp -cluster=cid -delete -snapshot -nonbuild -delpattern='(header|footer|zahlungsstoerung)-service'

	To use the script output to really delete the istags, you can use the following line:

		image-tools -family=pkp -cluster=cid -delete -snapshot -nonbuild -delpattern='(header|footer|zahlungsstoerung)-service'|xargs -n 1 -I{} bash -c "{}"
		
	To only generate a script to delete old snapshot istags:
	
		image-tools -family=pkp -cluster=cid -delete -snapshot

	To delete all not used images of family 'aps' in cluster cid
	
	    image-tools -family=aps -cluster=cid -delete  -minage=0 -delpattern='.'

	To delete all hybris istags of family pkp older than 45 days
	
		image-tools -family=pkp -cluster=cid -delete -isname=hybris -minage=45

CONNECTION

	If there are problems with the connection to the clusters,
	there is the option to disable the use of proxy with the
	parameter '-noproxy'.

	Or if a socks5 proxy can be the solution, eg. to run it from your notebook over VPN, then establish 
	a socks tunnel over the sprungserver and give the
	parameter '-socks5=ip:port' to the image-tools program.
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
	familyPtr := flag.String("family", "", "Mandatory: family name, eg. pkp, aps, ssp or fpc ")
	clusterPtr := flag.String("cluster", "", "Mandatory: shortname of cluster, eg. cid,int, ppr or pro")
	tokenPtr := flag.String("token", "", "Opt: token for cluster, its a alternative to login before exec")
	namespacePtr := flag.String("namespace", "", "Opt: namespace to look for istags")

	// Collect flags
	isPtr := flag.Bool("is", false, "Report: collect and report data of of existing imageStreams in the cluster")
	istagPtr := flag.Bool("istag", false, "Report: collect and report data of of existing imageStreamTags in the cluster")
	shaPtr := flag.Bool("image", false, "Report: collect and report data of existing Image's in the cluster")
	usedPtr := flag.Bool("used", false, "Report: collect and report data of used imageStreamTags from all clusters")
	unusedPtr := flag.Bool("unused", false, "Report: collect and report data of unused imageStreamTags from specified cluster")
	allPtr := flag.Bool("all", false, "Report: collect and report data all imageStreams, imageStreamTags, Image's and used imageStreamTags")

	// Output format of result data
	jsonPtr := flag.Bool("json", false, "Report: defines JSON as the output format for the reported data. This is the DEFAULT")
	yamlPtr := flag.Bool("yaml", false, "Report: defines YAML as the output format for the reported data")
	csvPtr := flag.Bool("csv", false, "Report: defines CSV as the output format for the reported data")
	csvFilePtr := flag.String("csvfile", "", "Report: defines CSV as the output format for the reported data and write the types to seperate csv files <name>-typ.csv")
	htmlPtr := flag.Bool("html", false, "Report: defines HTML as the output format for the reported data")
	tablePtr := flag.Bool("table", false, "Report: defines formated ASCI TABLE as the output format for the reported data")
	tabgroupPtr := flag.Bool("tabgroup", false, "Report: defines formated ASCII TABLE WITH GROUPED ROWS as the output format for the reported data")
	deletePtr := flag.Bool("delete", false, "Report: defines aktion to delete istags")

	// Filter flags
	isnamePtr := flag.String("isname", "", "Report and Delete: filter output for report or delete script of one imageStream")
	istagnamePtr := flag.String("istagname", "", "Report and Delete: filter output for report or delete script of one imageStreamTag")
	tagnamePtr := flag.String("tagname", "", "Report and Delete: filter output for report or delete script of all istags with this Tag")
	shanamePtr := flag.String("shaname", "", "Report and Delete: filter output for report or delete script of a Image with this SHA")

	// Delete flags
	deletePatternPtr := flag.String("delpattern", "", "Delete: filter for delete script all istags with this pattern")
	deleteMinAgePtr := flag.Int("minage", 60, "Delete: filter for delete script all istags, they are older or equal than minage")
	deleteNonBuildPtr := flag.Bool("nonbuild", false, "Delete: filter for delete script all istags with pure version number, where the referenced image has no build-tag and istag is minimum as old as minage")
	deleteSnapshotsPtr := flag.Bool("snapshot", false, "Delete: filter for delete script all istags, where the tag has SNAPSHOT or PR in the tagname and istag is minimum as old as minage")
	pruneImagesPtr := flag.Bool("prune", false, "Delete: prune images withpout istags")
	deleteConfirmPtr := flag.Bool("confirm", false, "Delete: confirm delete, if not set, it run in dry run mode")

	// Options
	ocClientPtr := flag.Bool("occlient", false, "TechOpt: use oc client instead of api call for cluster communication")
	noProxyPtr := flag.Bool("noproxy", false, "TechOpt: disable use of proxy for API http requests")
	socks5ProxyPtr := flag.String("socks5", "", "TechOpt: set socks5 proxy url and use it for API calls")
	profilerPtr := flag.Bool("profiler", false, "TechOpt: enable profiler support for debugging, http://localhost:6060/debug/pprof")
	noLogPtr := flag.Bool("nolog", false, "TechOpt: disable log output to screen. Logs to logfile is not disabled")
	debugPtr := flag.Bool("debug", false, "TechOpt: enable additional debug log output to screen")

	flag.Parse()

	// define map with all flags
	flags := T_flags{
		Cluster:  T_clName(*clusterPtr),
		Token:    string(*tokenPtr),
		Family:   T_family(*familyPtr),
		Json:     bool(*jsonPtr) || !(bool(*yamlPtr) || bool(*csvPtr) || bool(*deletePtr) || string(*csvFilePtr) != "" || bool(*tablePtr) || bool(*tabgroupPtr)),
		Yaml:     bool(*yamlPtr) && !bool(*jsonPtr),
		Csv:      (bool(*csvPtr) || (string(*csvFilePtr) != "")) && !bool(*jsonPtr),
		CsvFile:  string(*csvFilePtr),
		Delete:   bool(*deletePtr),
		Html:     bool(*htmlPtr) && !bool(*jsonPtr),
		Table:    bool(*tablePtr) && !bool(*jsonPtr),
		TabGroup: bool(*tabgroupPtr) && !bool(*jsonPtr),

		Output: T_flagOut{
			Is:     *isPtr,
			Istag:  *istagPtr,
			Image:  *shaPtr,
			Used:   *usedPtr,
			UnUsed: *unusedPtr,
			All:    *allPtr,
		},
		Filter: T_flagFilt{
			Isname:    T_isName(string(*isnamePtr)),
			Istagname: T_istagName(string(*istagnamePtr)),
			Tagname:   T_tagName(string(*tagnamePtr)),
			Imagename: T_shaName(string(*shanamePtr)),
			Namespace: T_nsName(string(*namespacePtr)),
		},
		DeleteOpts: T_flagDeleteOpts{
			Pattern:     *deletePatternPtr,
			MinAge:      *deleteMinAgePtr,
			NonBuild:    *deleteNonBuildPtr,
			Snapshots:   *deleteSnapshotsPtr,
			PruneImages: *pruneImagesPtr,
			Confirm:     *deleteConfirmPtr,
		},
		Options: T_flagOpts{
			OcClient:    *ocClientPtr,
			NoProxy:     *noProxyPtr,
			Socks5Proxy: *socks5ProxyPtr,
			Profiler:    *profilerPtr,
			NoLog:       *noLogPtr,
			Debug:       *debugPtr,
		},
	}

	CmdParams = flags
	LogMsg(GetJsonFromMap(flags))

	if flags.Family == "" {
		exitWithError("a name for family must given like: '-family=pkp'")
	}
	if !flags.Output.Used && (flags.Cluster == "") {
		exitWithError("a shortname for cluster must given like: '-cluster=cid'. Is now: ", flags.Cluster)
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

	if !(*isPtr || *istagPtr || *shaPtr || *allPtr || *usedPtr || *unusedPtr || *deletePtr) {
		exitWithError("As least one of the output flags must set")
	}
}
