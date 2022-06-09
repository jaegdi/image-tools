package ocrequest

import (
	"fmt"
	"log"
	"os"
	"regexp"
)

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

	InfoLogger.Println("------------------------------------------------------------")
	var currCluster T_clName
	if cl := os.Getenv("CLUSTER"); cl != "" {
		currCluster = T_clName(cl)
	} else {
		if len(CmdParams.Cluster) > 0 {
			currCluster = CmdParams.Cluster[0]
		} else {
			currCluster = "cid-apc0"
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
	fns := genFamilyNamespacesConfig(clustersConfig, familiesConfig, environmentsConfig, namespacesConfig, pipelinesConfig)
	FamilyNamespaces = fns
	InfoLogger.Println("------------------------------------------------------------")
	InfoLogger.Println("dynamic Config", GetJsonOneliner(fns))
	InfoLogger.Println("------------------------------------------------------------")
	InfoLogger.Println("Static Config", GetJsonOneliner(FamilyNamespacesStat))
	InfoLogger.Println("------------------------------------------------------------")

	EvalFlags()

	InfoLogger.Println("------------------------------------------------------------")
	InfoLogger.Println("Starting execution of image-tools")

	Multiproc = true
	InfoLogger.Println("disable proxy: " + fmt.Sprint(CmdParams.Options.NoProxy))
	InfoLogger.Println("Multithreading: " + fmt.Sprint(Multiproc))

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
