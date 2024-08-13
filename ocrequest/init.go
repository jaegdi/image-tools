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

	InfoMsg("------------------------------------------------------------")
	var currCluster T_clName

	if len(CmdParams.Cluster) > 0 {
		InfoMsg("Get cluster from parameter -cluster")
		currCluster = CmdParams.Cluster[0]
	} else {
		if cl := os.Getenv("CLUSTER"); cl != "" {
			InfoMsg("Set cluster to env var CLUSTER", cl)
			currCluster = T_clName(cl)
		} else {
			InfoMsg("Set cluster to default value -cluster=cid-scp0")
			currCluster = "cid-scp0"
		}
	}
	InfoMsg("Starting reading config from", currCluster, "config-tools")
	clustersConfig := GetClusters()
	familiesConfig := GetFamilies()
	environmentsConfig := GetEnvironments()
	namespacesConfig := GetNamespaces()
	pipelinesConfig := GetPipelines()
	// InfoMsg("Cluster Configs", clustersConfig)
	// InfoMsg("Environment Configs": environmentsConfig)
	// InfoMsg("NAmespace Configs": namespacesConfig)
	// InfoMsg("Pipeline Configs": pipelinesConfig)
	// cfg := genClusterConfig(clustersConfig)
	// InfoMsg("ClusterConfig", cfg)

	// use static config if cmdparam statcfg is true
	var fns T_famNsList
	if CmdParams.Options.StaticConfig || true {
		FamilyNamespaces = FamilyNamespacesStat
	} else {
		fns = genFamilyNamespacesConfig(clustersConfig, familiesConfig, environmentsConfig, namespacesConfig, pipelinesConfig)
		FamilyNamespaces = fns
	}

	InfoMsg("------------------------------------------------------------")
	InfoMsg("dynamic Config", GetJsonOneliner(fns))
	InfoMsg("------------------------------------------------------------")
	InfoMsg("Static Config", GetJsonOneliner(FamilyNamespacesStat))
	InfoMsg("------------------------------------------------------------")

	InfoMsg("############################################################")
	InfoMsg("Starting execution of image-tools")

	Multiproc = true
	InfoMsg("disable proxy: " + fmt.Sprint(CmdParams.Options.NoProxy))
	InfoMsg("Multithreading: " + fmt.Sprint(Multiproc))
	if CmdParams.Options.Socks5Proxy == "no" {
		CmdParams.Options.Socks5Proxy = ""
	}
	InfoMsg("Socks5Proxy: " + fmt.Sprint(CmdParams.Options.Socks5Proxy))
	InfoMsg("StaticConfig: " + fmt.Sprint(CmdParams.Options.StaticConfig))

	regexValidNamespace = regexp.MustCompile(`^` + string(CmdParams.Family) + `(?:-.*)?$`)
	// + `|` +
	// `^` + string(CmdParams.Family) + `-.*` + `|` +
	// `^.*?-` + string(CmdParams.Family) + `-.*` + `|` +
	// `^.*?-` + string(CmdParams.Family) + `$`)

	for _, cluster := range CmdParams.Cluster {
		if len(Clusters.Config[cluster].Token) < 10 {
			InfoMsg("Try to read clusterconfig.json")
			if err := readTokens("clusterconfig.json"); err != nil {
				InfoMsg("Read Clusterconfig is failed, try to get the tokens from clusters with oc login")
				for _, cluster := range FamilyNamespaces[CmdParams.Family].Stages {
					ocGetToken(cluster)
				}
				saveTokens(Clusters, "clusterconfig.json")
			} else {
				InfoMsg("Clusterconfig and Tokens loaded from clusterconfig.json")
			}
		}
	}
	InitIsNamesForFamily(CmdParams.Family)
}

func DebugMsg(p ...interface{}) {
	DebugLogger.Println(p...)
}

func InfoMsg(p ...interface{}) {
	InfoLogger.Println(p...)
}

func ErrorMsg(p ...interface{}) {
	ErrorLogger.Println(p...)
}
