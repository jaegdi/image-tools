package ocrequest

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// getClusterToken fetch login token from Clusters config
func getClusterToken(cluster T_clName) string {
	token := ""
	if _, ok := Clusters.Config[cluster]; ok {
		token = Clusters.Config[cluster].Token
		if token != "" {
			DebugMsg("Got token for cluster ", cluster, " from Clusters.Token")
		} else {
			ErrorMsg("Failed to get token for cluster ", cluster, " from Clusters.Token")
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
	DebugMsg("Try to get cluster token for cluster:", cluster)
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
				ExitWithError("Login failed to:", cluster)
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
	DebugMsg("Try to login: ", app, cluster, Clusters.Config[cluster].Name)
	cmd := exec.Command(app, Clusters.Config[cluster].Name)
	if stdout, err := cmd.Output(); err != nil {
		ErrorMsg("cmd: ", app, Clusters.Config[cluster].Name, err.Error()+":"+string(stdout))
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
		ErrorMsg("failed to serialize clusterconfig to json", err)
	} else {
		err := os.WriteFile(filePath, js, 0600)
		if err != nil {
			ErrorMsg("failed to save serialized clusterconfig as json to file", err)
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
	file, err := os.ReadFile(filePath)
	if err != nil {
		ErrorMsg("failed to load configfile "+filePath, err)
		return err
	} else {
		if err := json.Unmarshal([]byte(file), &Clusters); err != nil {
			ErrorMsg("error unmarshal clusterconfig from", filePath, err)
		} else {
			VerifyMsg("Token read from", filePath)
		}
		js, err := json.MarshalIndent(Clusters.Config, "", "    ")
		DebugMsg(string(js))
		return err
	}
}

// readTokens reads the token from the config file or a Kubernetes secret if running in ServerMode
func readServerTokens(filename string) error {
	if CmdParams.Options.ServerMode {
		// Read tokens from Kubernetes secret
		config, err := rest.InClusterConfig()
		if err != nil {
			ErrorMsg("failed to load in-cluster config", err)
			return err
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			ErrorMsg("failed to create Kubernetes client", err)
			return err
		}

		// Get the namespace the Pod is running in
		namespace, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
		if err != nil {
			ErrorMsg("failed to read namespace", err)
			return err
		}

		secret, err := clientset.CoreV1().Secrets(string(namespace)).Get(context.TODO(), "your-secret-name", metav1.GetOptions{})
		if err != nil {
			ErrorMsg("failed to get secret", err)
			return err
		}

		tokenData, ok := secret.Data["token"]
		if !ok {
			err := fmt.Errorf("token not found in secret")
			ErrorMsg("token not found in secret", err)
			return err
		}

		if err := json.Unmarshal(tokenData, &Clusters); err != nil {
			ErrorMsg("error unmarshal clusterconfig from secret", err)
			return err
		}

		VerifyMsg("Token read from Kubernetes secret")
		js, err := json.MarshalIndent(Clusters.Config, "", "    ")
		DebugMsg(string(js))
		return err
	} else {
		// Read tokens from config file
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		exPath := filepath.Dir(ex)
		filePath := exPath + "/" + filename
		file, err := os.ReadFile(filePath)
		if err != nil {
			ErrorMsg("failed to load configfile "+filePath, err)
			return err
		} else {
			if err := json.Unmarshal([]byte(file), &Clusters); err != nil {
				ErrorMsg("error unmarshal clusterconfig from", filePath, err)
			} else {
				VerifyMsg("Token read from", filePath)
			}
			js, err := json.MarshalIndent(Clusters.Config, "", "    ")
			DebugMsg(string(js))
			return err
		}
	}
}
