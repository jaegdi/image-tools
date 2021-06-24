package ocrequest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	// "log"
	"os/exec"
)

// getClusterToken fetch login token from Clusters config
func getClusterToken(cluster T_clName) string {
	token := ""
	if _, ok := Clusters.Config[cluster]; ok {
		token = Clusters.Config[cluster].Token
	}
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
		if cluster != CmdParams.Cluster || (cluster == CmdParams.Cluster && CmdParams.Token == "") {
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
	js, err := json.MarshalIndent(clusterconfig, "", "    ")
	if err != nil {
		LogError("failed to serialize clusterconfig to json", err)
	} else {
		err := ioutil.WriteFile(filename, js, 0600)
		if err != nil {
			LogError("failed to save serialized clusterconfig as json to file", err)
		}
	}
}

// readTokens reads the token from the config file
func readTokens(filename string) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		LogError("failed to load configfile "+filename, err)
		return err
	} else {
		if err := json.Unmarshal([]byte(file), &Clusters); err != nil {
			LogError("error unmarshal clusterconfig", err)
		}
		// js, err := json.MarshalIndent(Clusters.Config, "", "    ")
		// LogDebug(string(js))
		return err
	}
}
