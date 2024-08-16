package ocrequest

import (
	"fmt"
	"os"
	"regexp"
)

// Init is the intialization routine
func Init() {
	os.Setenv("HTTP_PROXY", "")
	// FamilyNamespaces = FamilyNamespacesStat

	EvalFlags()
	InitLogging()

	InfoMsg("------------------------------------------------------------")
	var currCluster T_clName

	if len(CmdParams.Cluster) > 0 && !CmdParams.Options.ServerMode {
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
	VerifyMsg("Cluster Configs", clustersConfig)
	VerifyMsg("Environment Configs", environmentsConfig)
	VerifyMsg("NAmespace Configs", namespacesConfig)
	VerifyMsg("Pipeline Configs", pipelinesConfig)
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

	regexValidNamespace = regexp.MustCompile(`^` + string(CmdParams.Family) + `(?:-.*)?$`)

	if !CmdParams.Options.ServerMode {
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
	}
	InitIsNamesForFamily(CmdParams.Family)
	if CmdParams.Options.ServerMode {
		InitServerMode(CmdParams)
		InfoMsg("ServerMode is enabled")
	}
	InfoMsg("disable proxy: " + fmt.Sprint(CmdParams.Options.NoProxy))
	InfoMsg("Multithreading: " + fmt.Sprint(Multiproc))
	InfoMsg("StaticConfig: " + fmt.Sprint(CmdParams.Options.StaticConfig))
	InfoMsg("Socks5Proxy: " + fmt.Sprint(CmdParams.Options.Socks5Proxy))
	if CmdParams.Options.Socks5Proxy == "no" {
		CmdParams.Options.Socks5Proxy = ""
	}
}

func InitServerMode(cp T_flags) {
	CmdParams.Family = cp.Family
	CmdParams.Filter = cp.Filter
	CmdParams.Output = cp.Output

	regexValidNamespace = regexp.MustCompile(`^` + string(CmdParams.Family) + `(?:-.*)?$`)
	CmdParams.FilterReg.Namespace = regexValidNamespace

	if CmdParams.Filter.Tagname != "" {
		CmdParams.FilterReg.Tagname = regexp.MustCompile(CmdParams.Filter.Tagname.str())
	}
	if CmdParams.Filter.Isname != "" {
		CmdParams.FilterReg.Isname = regexp.MustCompile(CmdParams.Filter.Isname.str())
	}
	if CmdParams.Filter.Istagname != "" {
		CmdParams.FilterReg.Istagname = regexp.MustCompile(CmdParams.Filter.Istagname.str())
	}

	// CmdParams.Options.Debug = true
	CmdParams.Options.Verify = true
}
