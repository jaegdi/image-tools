package ocrequest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	// "log"
	"os/exec"
)

// getClusterToken fetch login token from Clusters config
func getClusterToken(cluster T_clName) string {
	token := ""
	if _, ok := Clusters.Config[cluster]; ok {
		token = Clusters.Config[cluster].Token
		if token != "" {
			DebugLogger.Print("Got token for cluster ", cluster, " from Clusters.Token")
		} else {
			ErrorLogger.Print("Failed to get token for cluster ", cluster, " from Clusters.Token")
		}
	}
	// fmt.Println(token)
	// os.Exit(1)
	return token
}

// setClusterToken set login token onto Cluster config
func setClusterToken(cluster T_clName, token string) {
	if v, ok := Clusters.Config[cluster]; ok {
		v.Token = token
		Clusters.Config[cluster] = v
	}
}

// ocGetToken tries to get the oc token from config, if it is not defined in config, it requests it from command line parameter
func ocGetToken(cluster T_clName) string {
	LogDebug("Try to get cluster token for cluster:", cluster)
	token := getClusterToken(cluster)
	if token != "" {
		return token
	} else {
		if !CmdParams.Cluster.contains(cluster) || CmdParams.Cluster.contains(cluster) && CmdParams.Token == "" {
			t, err := ocLogin(cluster)
			if err != nil {
				saveTokens(Clusters, "clusterconfig.json")
				fmt.Println("Automatic login into the openshift cluster did not work on your machine.")
				fmt.Println("Ask Dirk Jäger from SCP Plattform Team, how to configure that on your machine.")
				fmt.Println("------- as an alternative you can do the following steps to work -------")
				fmt.Println("Fill in the tokens in the generated clusterconfig.json file and try again")
				fmt.Println("-- HowTo:")
				fmt.Println("-- To get your login tokens, exec 'oc login -u <your ldap user> https://console.<cluster-name>.sf-rz.de:8443'")
				fmt.Println("-- Then execute 'oc whoami -t' and put the token in the clusterconfig.json file.")
				fmt.Println("-- Repeat this for all clusters(int, ppr, pro) and put the pro token also to vpt.")
				fmt.Println("Then try again to exec this application.")
				exitWithError("Login failed to:", cluster)
			}
			setClusterToken(cluster, t)
			return t
		} else {
			setClusterToken(cluster, CmdParams.Token)
			return CmdParams.Token
		}
	}
}

// ocLogin tries to login with the token into the cluster
func ocLogin(cluster T_clName) (string, error) {
	app := "ocl"
	LogDebug("Try to login: ", app, Clusters.Config[cluster].Name)
	cmd := exec.Command(app, Clusters.Config[cluster].Name)
	if stdout, err := cmd.Output(); err != nil {
		LogError("cmd: ", app, Clusters.Config[cluster].Name, err.Error()+":"+string(stdout))

		return "login failed", err
	} else {
		cmd := exec.Command("oc", "whoami", "-t")
		token, err := cmd.Output()
		if err != nil {
			return "get token with 'oc whoami -t' failed:" + string(token), err
		}
		// token has a linefeed, that must removed
		return string(token[:len(token)-1]), err
	}
}

// saveTokens save the login token in the config file
func saveTokens(clusterconfig T_ClusterConfig, filename string) {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	filePath := exPath + "/" + filename
	js, err := json.MarshalIndent(clusterconfig, "", "    ")
	if err != nil {
		LogError("failed to serialize clusterconfig to json", err)
	} else {
		err := ioutil.WriteFile(filePath, js, 0600)
		if err != nil {
			LogError("failed to save serialized clusterconfig as json to file", err)
		}
	}
}

// readTokens reads the token from the config file
func readTokens(filename string) error {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	filePath := exPath + "/" + filename
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		LogError("failed to load configfile "+filePath, err)
		return err
	} else {
		if err := json.Unmarshal([]byte(file), &Clusters); err != nil {
			LogError("error unmarshal clusterconfig from", filePath, err)
		} else {
			InfoLogger.Println("Token read from", filePath)
		}
		// js, err := json.MarshalIndent(Clusters.Config, "", "    ")
		// LogDebug(string(js))
		return err
	}
}
