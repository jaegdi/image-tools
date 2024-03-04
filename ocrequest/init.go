package ocrequest

import (
	"fmt"
	"log"
	"os"
	"regexp"
)

type T_DebugLogger log.Logger

// Global Vars
var (
	WarningLogger       *log.Logger
	InfoLogger          *log.Logger
	ErrorLogger         *log.Logger
	DebugLogger         *log.Logger
	Multiproc           bool
	regexValidNamespace *regexp.Regexp
	LogfileName         string
)

// Init is the intialization routine
func Init() {
	err := os.Remove(LogFileName)
	logfile, err := os.OpenFile(LogFileName, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	os.Setenv("HTTP_PROXY", "")
	InfoLogger = log.New(logfile, "INFO: ", log.Ldate|log.Ltime|log.Lmsgprefix|log.Llongfile)
	WarningLogger = log.New(logfile, "WARNING: ", log.Ldate|log.Ltime|log.Lmsgprefix|log.Llongfile)
	ErrorLogger = log.New(logfile, "ERROR: ", log.Ldate|log.Ltime|log.Lmsgprefix|log.Llongfile)
	DebugLogger = log.New(logfile, "DEBUG: ", log.Ldate|log.Ltime|log.Lmsgprefix|log.Llongfile)

	// FamilyNamespaces = FamilyNamespacesStat
	EvalFlags()

	InfoLogger.Println("------------------------------------------------------------")
	var currCluster T_clName

	if len(CmdParams.Cluster) > 0 {
		InfoLogger.Println("Get cluster from parameter -cluster")
		currCluster = CmdParams.Cluster[0]
	} else {
		if cl := os.Getenv("CLUSTER"); cl != "" {
			InfoLogger.Println("Set cluster to env var CLUSTER", cl)
			currCluster = T_clName(cl)
		} else {
			InfoLogger.Println("Set cluster to default value -cluster=cid-scp0")
			currCluster = "cid-scp0"
		}
	}
	InfoLogger.Println("Starting reading config from", currCluster, "config-tools")
	clustersConfig := GetClusters()
	familiesConfig := GetFamilies()
	environmentsConfig := GetEnvironments()
	namespacesConfig := GetNamespaces()
	pipelinesConfig := GetPipelines()
	// InfoLogger.Println("Cluster Configs", clustersConfig)
	// InfoLogger.Println("Environment Configs": environmentsConfig)
	// InfoLogger.Println("NAmespace Configs": namespacesConfig)
	// InfoLogger.Println("Pipeline Configs": pipelinesConfig)
	// cfg := genClusterConfig(clustersConfig)
	// InfoLogger.Println("ClusterConfig", cfg)

	// use statig config if cmdparam statcfg is true
	var fns T_famNsList
	if CmdParams.Options.StaticConfig || true {
		FamilyNamespaces = FamilyNamespacesStat
	} else {
		fns = genFamilyNamespacesConfig(clustersConfig, familiesConfig, environmentsConfig, namespacesConfig, pipelinesConfig)
		FamilyNamespaces = fns
	}

	InfoLogger.Println("------------------------------------------------------------")
	InfoLogger.Println("dynamic Config", GetJsonOneliner(fns))
	InfoLogger.Println("------------------------------------------------------------")
	InfoLogger.Println("Static Config", GetJsonOneliner(FamilyNamespacesStat))
	InfoLogger.Println("------------------------------------------------------------")

	InfoLogger.Println("############################################################")
	InfoLogger.Println("Starting execution of image-tools")

	Multiproc = true
	InfoLogger.Println("disable proxy: " + fmt.Sprint(CmdParams.Options.NoProxy))
	InfoLogger.Println("Multithreading: " + fmt.Sprint(Multiproc))
	if CmdParams.Options.Socks5Proxy == "no" {
		CmdParams.Options.Socks5Proxy = ""
	}
	InfoLogger.Println("Socks5Proxy: " + fmt.Sprint(CmdParams.Options.Socks5Proxy))
	InfoLogger.Println("StaticConfig: " + fmt.Sprint(CmdParams.Options.StaticConfig))

	regexValidNamespace = regexp.MustCompile(
		`^` + string(CmdParams.Family) + `$` + `|` +
			`^` + string(CmdParams.Family) + `-.*` + `|` +
			`.*?-` + string(CmdParams.Family) + `-.*` + `|` +
			`.*?-` + string(CmdParams.Family) + `$`)

	for _, cluster := range CmdParams.Cluster {
		if len(Clusters.Config[cluster].Token) < 10 {
			InfoLogger.Println("Try to read clusterconfig.json")
			if err := readTokens("clusterconfig.json"); err != nil {
				InfoLogger.Println("Read Clusterconfig is failed, try to get the tokens from clusters with oc login")
				for _, cluster := range FamilyNamespaces[CmdParams.Family].Stages {
					ocGetToken(cluster)
				}
				saveTokens(Clusters, "clusterconfig.json")
			} else {
				InfoLogger.Println("Clusterconfig and Tokens loaded from clusterconfig.json")
			}
		}
	}
	InitIsNamesForFamily(CmdParams.Family)
}
