package ocrequest

import (
	// "fmt"
	"log"
	"os/exec"
)

type T_clusterTokens struct {
	Cid string
	Int string
	Ppr string
	Vpt string
	Pro string
}

var clusterTokens = T_clusterTokens{}

func getClusterToken(cluster string) string {
	token := ""
	switch cluster {
	case "cid-apc0", "cid":
		token = clusterTokens.Cid
	case "int-apc0", "int":
		token = clusterTokens.Int
	case "ppr-apc0", "ppr":
		token = clusterTokens.Ppr
	case "vpt-apc0", "vpt":
		token = clusterTokens.Vpt
	case "pro-apc0", "pro":
		token = clusterTokens.Ppr
	}
	return token
}

func setClusterToken(cluster string, token string) {
	switch cluster {
	case "cid-apc0", "cid":
		clusterTokens.Cid = token
	case "int-apc0", "int":
		clusterTokens.Int = token
	case "ppr-apc0", "ppr":
		clusterTokens.Ppr = token
	case "vpt-apc0", "vpt":
		clusterTokens.Vpt = token
	case "pro-apc0", "pro":
		clusterTokens.Ppr = token
	}
}

func ocGetToken(cluster string) string {
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
	cmd := exec.Command(app, cluster)

	if stdout, err := cmd.Output(); err != nil {
		log.Println(err.Error() + ":" + string(stdout))
		return "login failed", err
	} else {
		cmd := exec.Command("oc", "whoami", "-t")
		token, err := cmd.Output()
		if err != nil {
			return "token failed:" + string(token), err
		}
		// token has a linefeed, that must removed
		return string(token[:len(token)-1]), err
	}
}
