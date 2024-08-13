package ocrequest

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
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


  Per defaul image-tools is executed as a cmdline tool but with parameter -server it can be startet as a webservice.

  image-tools as web-service

    In serverMode image-tools looks for used images in all clusters filtered by family and tagname and returns the a
    JSON list as HTTP response

    image-tools starts a webserver to listen on port 8080.
    To request the webservice by curl, use the following pattern:

        curl "http://localhost:8080/execute?family=exampleFamily&tagname=exampleTagname" | jq

  image-tools as cmdline tool

    image-tools only read information around images from the clusters and generate output. It never change or
    delete something in the clusters. Eg. for delete istags it only generates a script output, which than can
    be executed by a cluster admin to really delete the istags.

    image-tools always write a log to 'image-tools.log' in the current directory.
    Additional it writes the log messages also to STDERR. To disable the log output to STDERR use parameter -nolog
    The report or delete-script output is written to STDOUT.

  - For reporting existing Images and delete istags it operates cluster and family specific.
    For reporting used images it works over all clusters but family specific.
    It never works across different families.

      For existing images, istags or imagestreams(is) that means it works for one or more clusters like
        cid-scp0,  ppr-scp0, vpt-scp0, pro-scp0
        or more than one like 'cid-scp0,cid-scp0,ppr-scp0,ppr-scp0'

    and for one families like
        'pkp, sps, fpc, aps, ...'

    The cluster must be defined by the mandatory parameter
        '-cluster=[cid-scp0|int-scp0|ppr-scp0|vpt-scp0|pro-scp0|dev-scp0|cid-scp0|ppr-scp0|vpt-scp0|pro-scp0]'

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
      -delminage=int      defines the minimum age (in days) for istag to delete them. Default is 60 days.
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

    image-tools is a statically linked go programm and has no runtime dependencies. No installation is
    neccessary. Copy the binary from artifactory
        - linux:  https://artifactory-pro.sf-rz.de:8443/artifactory/scpas-bin-develop/istag_and_image_management/image-tools-linux/image-tools
        - windows:  https://artifactory-pro.sf-rz.de:8443/artifactory/scpas-bin-develop/istag_and_image_management/image-tools-windows/image-tools.exe
    into a directory, which is in the search path is enough.


EXAMPLES

  REPORTING
  |
  |        Report all information for family pkp in cluster cid as json
  |        (which is the default output format)
  |
  |            image-tools -cluster=cid-scp0 -family=pkp -all
  |
  |            or as table
  |            image-tools -cluster=cid-scp0 -family=pkp -all -table
  |
  |            or csv in different files for each type of information
  |            image-tools -cluster=cid-scp0 -family=pkp -all -csvfile=prefix
  |            writes the output to different files 'prefix-type' in current directory
  |
  |        Report only __used__ istags for family pkp as pretty printed table
  |        (the output is paginated to fit your screen size and piped to
  |            the pager define in the environment variable $PAGER/%PAGER%.
  |            If $PAGER is not set, it try to use 'more')
  |
  |            image-tools -cluster=cid-scp0 -family=pkp -used -table
  |            or json
  |            image-tools -cluster=cid-scp0 -family=pkp -used
  |            or yaml
  |            image-tools -cluster=cid-scp0 -family=pkp -used -yaml
  |            or csv
  |            image-tools -cluster=cid-scp0 -family=pkp -used -csv
  |
  |        Report istags for family aps in cluster int as yaml report
  |
  |            image-tools -cluster=int-scp0 -family=aps -istag -yaml
  |
  |        Report ImageStreams for family aps in cluster int as yaml report
  |
  |            image-tools -cluster=int-scp0 -family=aps -is -yaml
  |
  |        Report Images for family aps in cluster int as yaml report
  |
  |            image-tools -cluster=int-scp0 -family=aps -image -yaml
  |
  |        Report combined with pc(print columns) tool
  |
  |            image-tools -socks5=localhost:65022 -family=pkp -cluster=cid-scp0,ppr-scp0,pro-scp0 -istag -csv | pc -sep=, -sortcol=4  1 5 8 6 7



  DELETE
  |
  |        Generate a shell script to delete old istags(60 days, the default) for family pkp in cluster cid
  |        and all old snapshot istags and nonbuild istags and all istags of header-service, footer-service
  |        and zahlungsstoerung-service
  |
  |            image-tools -family=pkp -cluster=cid-scp0 -delete -snapshot \
  |                        -nonbuild -delpattern='(header|footer|zahlungsstoerung)-service'
  |
  |        To use the script output to really delete the istags, you can use the following line:
  |
  |            image-tools -family=pkp -cluster=cid-scp0 -delete -snapshot -nonbuild \
  |                        -delpattern='(header|footer|zahlungsstoerung)-service'      | xargs -n 1 -I{} bash -c "{}"
  |
  |        To only generate a script to delete old snapshot istags:
  |
  |            image-tools -family=pkp -cluster=cid-scp0 -delete -snapshot
  |
  |        To delete all not used images of family 'aps' in cluster cid
  |
  |            image-tools -family=aps -cluster=cid-scp0 -delete  -delminage=0 -delpattern='.'
  |
  |        To delete all hybris istags of family pkp older than 45 days
  |
  |            image-tools -family=pkp -cluster=cid-scp0 -delete -isname=hybris -delminage=45


  HINT
  |
  |        To directly delete the istags, that reportet by 'image-tools -delete ...', make shure, you are
  |        logged in into the correct cluster, because the output is executed with oc client and work on the
  |        currently logged in cluster. And append the following to the end of the image-tools - command:
  |
  |
  |            | xargs -n 1 -I{} bash -c "{}"
  |
  |        After deleting the istags, the images must removed from the registry by executing a command similar
  |        to this example:
  |
  |            oc login ..... to the cluster
  |            registry_url="$(oc -n default get route|grep docker-registry|awk '{print $2}')"
  |            oc adm prune images --registry-url=$registry_url --keep-tag-revisions=3 --keep-younger-than=60m --confirm


CONNECTION

    As default the sock5 proxy to localhost:65022 is enabled becaus the api of the upper clusters is only reacheable
	over the sprungserver. To disable SOCKS5 set the parameter -socks5=no
	If your socks5 jumpserver config listens on a different port set the parameter -socks5=<host>:<port>

    If there are problems with the connection to the clusters,
    there is the option to disable the use of web proxy with the
    parameter '-noproxy'.

    A socks5 proxy can be the solution, eg. to run it from your notebook over VPN, then establish
    a socks tunnel over the sprungserver and give the
    parameter '-socks5=host:port' to the image-tools program.
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
	familyPtr := flag.String("family", "", "Mandatory: family name, eg.: "+FamilyNamespaces.familyListStr())
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
	flags := T_flags{
		Cluster:  T_clName(*clusterPtr).list(),
		Token:    string(*tokenPtr),
		Family:   T_familyName(*familyPtr),
		App:      T_appName(*appPtr),
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
			Minage:    *minAgePtr,
			Maxage:    *maxAgePtr,
		},
		FilterReg: T_flagFiltRegexp{
			Isname:    is_r,
			Istagname: istag_r,
			Tagname:   tag_r,
			Namespace: ns_r,
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
		},
	}

	CmdParams = flags
	InfoMsg(GetJsonFromMap(flags))

	if flags.Family == "" && !flags.Options.ServerMode {
		exitWithError("a name for family must given like: '-family=pkp'")
	}
	if !flags.Output.Used && (len(flags.Cluster) == 0) {
		exitWithError("a shortname for cluster must given like: '-cluster=cid-scp0'. Is now: ", flags.Cluster)
	}
	// for _, cluster := range flags.Cluster {
	// 	_, clusterDefined := Clusters.Config[cluster]
	// 	if !clusterDefined && !flags.Output.Used {
	// 		clusterlist := []string{}
	// 		for clname := range Clusters.Config {
	// 			clusterlist = append(clusterlist, string(clname))
	// 		}
	// 		clusters := strings.Join(clusterlist, ",")
	// 		exitWithError("The clustername given as -cluster= is not defined: Given: ", cluster, " valid names: ", clusters)
	// 	}
	// }
	// if FamilyNamespaces[flags.Family].ImageNamespaces == nil {
	// 	exitWithError("Family", flags.Family, "is not defined")
	// }

	// for _, cluster := range flags.Cluster {
	// 	foundNamespace := false
	// 	for _, v := range FamilyNamespaces[flags.Family].ImageNamespaces[cluster] {
	// 		if flags.Filter.Namespace == v {
	// 			foundNamespace = true
	// 		}
	// 	}
	// 	if !foundNamespace && !(flags.Filter.Namespace == "") && !flags.Output.Used {
	// 		exitWithError("Namespace", flags.Filter.Namespace, "is no image namespace for family", flags.Family)
	// 	}
	// }

	if !flags.Options.ServerMode && !(*isPtr || *istagPtr || *shaPtr || *allPtr || *usedPtr || *unusedPtr || *deletePtr) {
		exitWithError("As least one of the output flags must set")
	}
	if flags.Options.ServerMode {
		CmdParams.Output.Used = true
	}
}
