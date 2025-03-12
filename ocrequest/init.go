package ocrequest

import (
	"fmt"
	"os"
	"regexp"
)

// Init is the initialization routine for setting up the environment and configurations.
// It sets environment variables, initializes logging, evaluates command-line flags, and reads configurations.
func Init() {
	// Clear the HTTP_PROXY environment variable
	os.Setenv("HTTP_PROXY", "")

	// Initialize logging
	EvalFlags()
	InitLogging()

	InfoMsg("------------------------------------------------------------")
	var currCluster T_clName

	// Determine the current cluster based on command-line parameters or environment variables
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

	// Load static or dynamic configuration based on command-line options
	if CmdParams.Options.StaticConfig {
		// Use static config if cmdparam statcfg is true
		FamilyNamespaces = FamilyNamespacesStat
		// Log the static configurations
		InfoMsg("------------------------------------------------------------")
		InfoMsg("Static Config", GetJsonOneliner(FamilyNamespacesStat))
	} else {
		InfoMsg("Starting reading config from", currCluster, "config-tools")
		// Read various configurations
		clustersConfig := GetClusters()
		familiesConfig := GetFamilies()
		environmentsConfig := GetEnvironments()
		namespacesConfig := GetNamespaces()
		pipelinesConfig := GetPipelines()

		// Log the configurations
		VerifyMsg("Cluster Configs", clustersConfig)
		VerifyMsg("Environment Configs", environmentsConfig)
		VerifyMsg("Namespace Configs", namespacesConfig)
		VerifyMsg("Pipeline Configs", pipelinesConfig)

		// Generate family namespaces configuration
		var fns T_famNsList
		fns = genFamilyNamespacesConfig(clustersConfig, familiesConfig, environmentsConfig, namespacesConfig, pipelinesConfig)
		FamilyNamespaces = fns
		// Log the static configurations
		InfoMsg("------------------------------------------------------------")
		InfoMsg("dynamic Config", GetJsonFromMap(fns))
	}
	// check some parameters
	if FamilyNamespaces[CmdParams.Family].ImageNamespaces == nil && !(CmdParams.Family == "all" || CmdParams.Options.ServerMode) {
		ExitWithError("Family ", CmdParams.Family, " is not defined")
	}
	for _, cluster := range CmdParams.Cluster {
		foundNamespace := false
		for _, v := range FamilyNamespaces[CmdParams.Family].ImageNamespaces[cluster] {
			if CmdParams.Filter.Namespace == v {
				foundNamespace = true
			}
		}
		if !foundNamespace && !(CmdParams.Filter.Namespace == "") && !CmdParams.Output.Used {
			ExitWithError("Namespace ", CmdParams.Filter.Namespace, " is no image namespace for family ", CmdParams.Family)
		}
	}

	// End of Log the dynamic and static configurations
	InfoMsg("------------------------------------------------------------")

	InfoMsg("############################################################")
	InfoMsg("Starting execution of image-tools")

	// Enable multiprocessing
	Multiproc = true

	// Compile a regular expression for validating namespaces
	regexValidNamespace = regexp.MustCompile(`^` + string(CmdParams.Family) + `(?:-.*)?$`)

	// If not in server mode, read cluster configurations
	if !CmdParams.Options.ServerMode {
		for _, cluster := range CmdParams.Cluster {
			// Check if the token length is less than 10
			if len(Clusters.Config[cluster].Token) < 10 {
				InfoMsg("Try to read clusterconfig.json")
				// Attempt to read tokens from clusterconfig.json
				if err := readTokens("clusterconfig.json"); err != nil {
					InfoMsg("Read Clusterconfig is failed, try to get the tokens from clusters with oc login")
					// If reading tokens fails, get tokens from clusters using oc login
					for _, cluster := range FamilyNamespaces[CmdParams.Family].Stages {
						ocGetToken(cluster)
					}
					// Save the tokens to clusterconfig.json
					saveTokens(Clusters, "clusterconfig.json")
				} else {
					InfoMsg("Clusterconfig and Tokens loaded from clusterconfig.json")
				}
			}
		}
	}
	// Initialize image stream names for the family
	InitIsNamesForFamily(CmdParams.Family)
	// If server mode is enabled, initialize server mode
	if CmdParams.Options.ServerMode {
		InitServerMode(CmdParams)
		InfoMsg("ServerMode is enabled")
	}
	// Log various options
	InfoMsg("disable proxy: " + fmt.Sprint(CmdParams.Options.NoProxy))
	InfoMsg("Multithreading: " + fmt.Sprint(Multiproc))
	InfoMsg("StaticConfig: " + fmt.Sprint(CmdParams.Options.StaticConfig))
	InfoMsg("Socks5Proxy: " + fmt.Sprint(CmdParams.Options.Socks5Proxy))
	if CmdParams.Options.Socks5Proxy == "no" {
		CmdParams.Options.Socks5Proxy = ""
	}
}

// InitServerMode initializes the server mode with the provided command-line parameters.
// It sets up regular expressions for filtering namespaces, tag names, image stream names, and image stream tag names.
//
// Parameters:
// - cp: The command-line parameters to use for initializing the server mode.
func InitServerMode(cp T_flags) {
	// Set command parameters from the provided flags
	// CmdParams.Family = cp.Family
	// CmdParams.Filter = cp.Filter
	// CmdParams.Output = cp.Output

	// Compile a regular expression for validating namespaces
	if CmdParams.Family == "scp" {
		regexValidNamespace = regexp.MustCompile(`^(?:` + string(CmdParams.Family) + `|ocp)(?:-.*)?$`)
	} else {
		regexValidNamespace = regexp.MustCompile(`^` + string(CmdParams.Family) + `(?:-.*)?$`)
	}

	// Compile regular expressions for filtering tag names, image stream names, and image stream tag names
	if CmdParams.Filter.Namespace != "" {
		CmdParams.FilterReg.Namespace = regexp.MustCompile(CmdParams.Filter.Namespace.str())
	} else {
		CmdParams.FilterReg.Namespace = regexValidNamespace
	}

	if CmdParams.Filter.Tagname != "" {
		CmdParams.FilterReg.Tagname = regexp.MustCompile(CmdParams.Filter.Tagname.str())
	}
	if CmdParams.Filter.Isname != "" {
		CmdParams.FilterReg.Isname = regexp.MustCompile(CmdParams.Filter.Isname.str())
	}
	if CmdParams.Filter.Istagname != "" {
		CmdParams.FilterReg.Istagname = regexp.MustCompile(CmdParams.Filter.Istagname.str())
	}

	// Enable verification option
	// CmdParams.Options.Verify = true
}
