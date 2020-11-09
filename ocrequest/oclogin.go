package ocrequest

import (
	"os/exec"
)

func getClusterToken(cluster string) string {
	token := ""
	if _, ok := Clusters[cluster]; ok {
		token = Clusters[cluster].(T_Cluster).Token
	}
	return token
}

func setClusterToken(cluster string, token string) {
	if v, ok := Clusters[cluster].(T_Cluster); ok {
		v.Token = token
		Clusters[cluster] = v
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
	InfoLogger.Println("Try to login: ", app, Clusters[cluster].(T_Cluster).Name)
	cmd := exec.Command(app, Clusters[cluster].(T_Cluster).Name)
	if stdout, err := cmd.Output(); err != nil {
		ErrorLogger.Println("cmd: ", app, Clusters[cluster].(T_Cluster).Name, err.Error()+":"+string(stdout))
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
