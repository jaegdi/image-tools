package ocrequest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
)

func getClusterToken(cluster string) string {
	token := ""
	if _, ok := Clusters.Config[cluster]; ok {
		token = Clusters.Config[cluster].Token
	}
	return token
}

func setClusterToken(cluster string, token string) {
	if v, ok := Clusters.Config[cluster]; ok {
		v.Token = token
		Clusters.Config[cluster] = v
	}
}

func ocGetToken(cluster string) string {
	InfoLogger.Println("Try to get cluster token for cluster: ", cluster)
	token := getClusterToken(cluster)
	if token != "" {
		return token
	} else {
		if cluster != CmdParams.Cluster || (cluster == CmdParams.Cluster && CmdParams.Token == "") {
			t, err := ocLogin(cluster)
			if err != nil {
				saveTokens(Clusters, "clusterconfig.json")
				fmt.Println("Automatic login into the openshift cluster did not work on your machine.")
				fmt.Println("Ask Dirk Jäger from SCP Plattform Team, how to konfigure that on your machine.")
				fmt.Println("------- as an alternative you can do the following steps to work -------")
				fmt.Println("Fill in the tokens in the generated clusterconfig.json file and try again")
				fmt.Println("-- HowTo:")
				fmt.Println("-- To get your login tokens, exec 'oc login -u <your ldap user> https://console.cid-apc0.sf-rz.de:8443'")
				fmt.Println("-- Then execute 'oc whoami -t' and put the token in the clusterconfig.json file.")
				fmt.Println("-- Repeat this for all clusters(int, ppr, pro) and put the pro token also to vpt.")
				fmt.Println("Then try again to exec this application.")
				exitWithError("Login failed to: " + cluster)
			}
			setClusterToken(cluster, t)
			return t
		} else {
			setClusterToken(cluster, CmdParams.Token)
			return CmdParams.Token
		}
	}
}

func ocLogin(cluster string) (string, error) {
	app := "ocl"
	InfoLogger.Println("Try to login: ", app, Clusters.Config[cluster].Name)
	cmd := exec.Command(app, Clusters.Config[cluster].Name)
	if stdout, err := cmd.Output(); err != nil {
		ErrorLogger.Println("cmd: ", app, Clusters.Config[cluster].Name, err.Error()+":"+string(stdout))

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

func saveTokens(clusterconfig T_ClusterConfig, filename string) {
	js, err := json.MarshalIndent(clusterconfig, "", "    ")
	if err != nil {
		log.Println("failed to serialize clusterconfig to json", err)
	} else {
		err := ioutil.WriteFile(filename, js, 0600)
		if err != nil {
			log.Println("failed to save serialized clusterconfig as json to file", err)
		}
	}
}

func readTokens(filename string) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println("failed to load configfile "+filename, err)
		return err
	} else {
		if err := json.Unmarshal([]byte(file), &Clusters); err != nil {
			log.Println("error unmarshal clusterconfig", err)
		}
		// js, err := json.MarshalIndent(Clusters.Config, "", "    ")
		// log.Println(string(js))
		return err
	}
}
