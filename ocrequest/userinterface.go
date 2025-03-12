package ocrequest

import (
	"flag"
	"os"
	"regexp"
	"strings"
)

var CmdParams T_flags

// FamilyNamespaces T_famNs

// EvalFlags evaluate all command line flags and set a struct with their values
func EvalFlags() {
	flag.Usage = CmdUsage
	manPtr := flag.Bool("man", false, "Print the ManPage")

	// Global Flags
	familyPtr := flag.String("family", "", "Mandatory: family or appgroup name, eg.: "+FamilyNamespaces.familyListStr())
	appgroupPtr := flag.String("appgroup", "", "Mandatory: appgroup or family name, eg.: "+FamilyNamespaces.familyListStr())
	appPtr := flag.String("app", "", "Mandatory: app name, eg.: "+AppNamespaces.appListStr())
	clusterPtr := flag.String("cluster", "", "Mandatory: name of one or more cluster, eg.: "+FamilyNamespaces[T_familyName("base")].clusterListStr()+" or more than one: cid-scp0,cid-scp0,ppr-acp0")
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
	csvFilePtr := flag.String("csvfile", "", "Report: defines CSV as the output format for the reported data and\n"+
		"write the types to seperate csv files <name>-typ.csv")
	htmlPtr := flag.Bool("html", false, "Report: defines HTML as the output format for the reported data")
	tablePtr := flag.Bool("table", false, "Report: defines formated ASCI TABLE as the output format for the reported data")
	tabgroupPtr := flag.Bool("tabgroup", false, "Report: defines formated ASCII TABLE WITH GROUPED ROWS as the output format for the reported data")
	deletePtr := flag.Bool("delete", false, "Report: defines aktion to delete istags")

	// Filter flags
	isnamePtr := flag.String("isname", "", "Report and Delete: filter output for report or delete script of one imageStream")
	istagnamePtr := flag.String("istagname", "", "Report and Delete: filter output for report or delete script of one imageStreamTag")
	tagnamePtr := flag.String("tagname", "", "Report and Delete: filter output for report or delete script of all istags with this Tag")
	shanamePtr := flag.String("shaname", "", "Report and Delete: filter output for report or delete script of a Image with this SHA")
	minAgePtr := flag.Int("minage", -1, "Report and Delete: filter for all istags, they are older or equal than minage")
	maxAgePtr := flag.Int("maxage", -1, "Report and Delete: filter for all istags, they are younger or equal than minage")

	// Delete flags
	deletePatternPtr := flag.String("delpattern", "", "Delete: filter for delete script all istags with this pattern")
	deleteMinAgePtr := flag.Int("delminage", 60, "Delete: filter for delete script all istags, they are older or equal than minage")
	deleteNonBuildPtr := flag.Bool("nonbuild", false, "Delete: filter for delete script all istags with pure version number,\n"+
		"where the referenced image has no build-tag and istag is minimum as old as minage")
	deleteSnapshotsPtr := flag.Bool("snapshot", false, "Delete: filter for delete script all istags,\n"+
		"where the tag has SNAPSHOT or PR in the tagname and istag is minimum as old as minage")
	pruneImagesPtr := flag.Bool("prune", false, "Delete: prune images withpout istags")
	deleteConfirmPtr := flag.Bool("confirm", false, "Delete: confirm delete, if not set, it run in dry run mode")

	// Options
	insecurePtr := flag.Bool("insecure-ssl", false, "Accept/Ignore all server SSL certificates")
	ocClientPtr := flag.Bool("occlient", false, "TechOpt: use oc client instead of api call for cluster communication")
	noProxyPtr := flag.Bool("noproxy", false, "TechOpt: disable use of proxy for API http requests")
	statCfgPtr := flag.Bool("statcfg", false, "TechOpt: use the static defined config instead of dynamic generated config from config-tool repos")

	socks5proxy := os.Getenv("SOCKS5PROXY")
	if len(strings.TrimSpace(socks5proxy)) == 0 {
		socks5proxy = "127.0.0.1:65022"
	}
	if socks5proxy == "no" {
		socks5proxy = ""
	}
	socks5ProxyPtr := flag.String("socks5", socks5proxy, "TechOpt: set socks5 proxy url and use it for API calls.\n eg. -socks5=127.0.0.1:65022 .  If env var SOCKS5PROXY is defined, it uses this as default, otherwise '127.0.0.1:65022' string")
	profilerPtr := flag.Bool("profiler", false,
		"TechOpt: enable profiler support for debugging, http://localhost:6060/debug/pprof\n"+
			" or: ~/go/bin/pprof -http localhost:8080 http://localhost:6060/debug/pprof/goroutine")
	noLogPtr := flag.Bool("nolog", false, "TechOpt: disable log output to screen. Logs to logfile is not disabled")
	debugPtr := flag.Bool("debug", false, "TechOpt: enable additional debug log output to screen")
	verifyPtr := flag.Bool("verify", false, "TechOpt: enable additional log output to screen")
	serverPtr := flag.Bool("server", false, "TechOpt: start in server mode and provide a rest api")

	flag.Parse()

	//  Print manpage and exit it param -man is set in cmdline
	if bool(*manPtr) {
		ManPage()
		os.Exit(0)
	}

	var is_r *regexp.Regexp
	var istag_r *regexp.Regexp
	var tag_r *regexp.Regexp
	var ns_r *regexp.Regexp
	if string(*isnamePtr) != "" {
		is_r = regexp.MustCompile(string(*isnamePtr))
	}
	if string(*istagnamePtr) != "" {
		istag_r = regexp.MustCompile(string(*istagnamePtr))
	}
	if string(*tagnamePtr) != "" {
		tag_r = regexp.MustCompile(string(*tagnamePtr))
	}
	if string(*namespacePtr) != "" {
		ns_r = regexp.MustCompile(string(*namespacePtr))
	}

	// define map with all flags
	flags := T_flags{}
	flags.Cluster = T_clName(*clusterPtr).list()
	flags.Token = string(*tokenPtr)
	if *familyPtr != "" {
		flags.Family = T_familyName(*familyPtr)
	} else {
		flags.Family = T_familyName(*appgroupPtr)
	}
	flags.App = T_appName(*appPtr)
	flags.Json = bool(*jsonPtr) || !(bool(*yamlPtr) || bool(*csvPtr) || bool(*deletePtr) || string(*csvFilePtr) != "" || bool(*tablePtr) || bool(*tabgroupPtr))
	flags.Yaml = bool(*yamlPtr) && !bool(*jsonPtr)
	flags.Csv = (bool(*csvPtr) || (string(*csvFilePtr) != "")) && !bool(*jsonPtr)
	flags.CsvFile = string(*csvFilePtr)
	flags.Delete = bool(*deletePtr)
	flags.Html = bool(*htmlPtr) && !bool(*jsonPtr)
	flags.Table = bool(*tablePtr) && !bool(*jsonPtr)
	flags.TabGroup = bool(*tabgroupPtr) && !bool(*jsonPtr)

	flags.Output = T_flagOut{
		Is:     *isPtr,
		Istag:  *istagPtr,
		Image:  *shaPtr,
		Used:   *usedPtr,
		UnUsed: *unusedPtr,
		All:    *allPtr,
	}
	flags.Filter = T_flagFilt{
		Isname:    T_isName(string(*isnamePtr)),
		Istagname: T_istagName(string(*istagnamePtr)),
		Tagname:   T_tagName(string(*tagnamePtr)),
		Imagename: T_shaName(string(*shanamePtr)),
		Namespace: T_nsName(string(*namespacePtr)),
		Minage:    *minAgePtr,
		Maxage:    *maxAgePtr,
	}
	flags.FilterReg = T_flagFiltRegexp{
		Isname:    is_r,
		Istagname: istag_r,
		Tagname:   tag_r,
		Namespace: ns_r,
	}
	flags.DeleteOpts = T_flagDeleteOpts{
		Pattern:     *deletePatternPtr,
		MinAge:      *deleteMinAgePtr,
		NonBuild:    *deleteNonBuildPtr,
		Snapshots:   *deleteSnapshotsPtr,
		PruneImages: *pruneImagesPtr,
		Confirm:     *deleteConfirmPtr,
	}
	flags.Options = T_flagOpts{
		InsecureSSL:  *insecurePtr,
		OcClient:     *ocClientPtr,
		NoProxy:      *noProxyPtr,
		Socks5Proxy:  *socks5ProxyPtr,
		Profiler:     *profilerPtr,
		NoLog:        *noLogPtr,
		Debug:        *debugPtr,
		Verify:       *verifyPtr,
		ServerMode:   *serverPtr,
		StaticConfig: *statCfgPtr,
	}

	CmdParams = flags
	// VerifyMsg("-- Test --")
	// VerifyMsg(GetJsonFromMap(flags))

	// If CmdParams.Cluster is empty, set it to all clusters
	if !CmdParams.Delete &&
		(len(CmdParams.Cluster[0]) == 0 ||
			CmdParams.Cluster[0] == "all") {
		CmdParams.Cluster = Clusters.getClusterList()
	}

	if CmdParams.Delete &&
		(len(CmdParams.Cluster[0]) == 0 ||
			CmdParams.Cluster[0] == "all") {
		envCluster := os.Getenv("CLUSTER")
		if envCluster != "" {
			CmdParams.Cluster = T_clName(envCluster).list()
		} else {
			ExitWithError("\nA name for cluster must given like: '-cluster=cid-scp0'. Is now: ", CmdParams.Cluster[0])
		}
	}

	if CmdParams.Family == "" && !CmdParams.Options.ServerMode {
		ExitWithError("\nA name for family must given like: '-family=pkp'")
	}

	// if !CmdParams.Options.ServerMode && !CmdParams.Output.Used && len(CmdParams.Cluster[0]) == 0 {
	// 	ExitWithError("\nA shortname for cluster must given like: '-cluster=cid-scp0' or -cluster=cid-scp0,ppr-scp0. Is now: ", flags.Cluster[0])
	// }
	if !CmdParams.Options.ServerMode {
		for _, cluster := range CmdParams.Cluster {
			_, clusterDefined := Clusters.Config[cluster]
			if !clusterDefined && !CmdParams.Output.Used {
				clusterlist := []string{}
				for clname := range Clusters.Config {
					clusterlist = append(clusterlist, string(clname))
				}
				// ExitWithError("The clustername given as -cluster= is not defined: Given: ", cluster, " valid names: ", strings.Join(clusterlist, ","))
			}
		}
	}
	if !CmdParams.Options.ServerMode && !(*isPtr || *istagPtr || *shaPtr || *allPtr || *usedPtr || *unusedPtr || *deletePtr) {
		ExitWithError("\nAs least one of the output flags must set")
	}
}
