// Package ocrequest provides primitoives to query an oc-cluster for istags
// calculate and filter them. And provides a json output primitive.
package ocrequest

import (
	"golang.org/x/net/proxy"
	"io/ioutil"
	"net/http"
	"os/exec"
)

// ocApiCall requests an openshift-cluster via API and return the answer as string.
func ocGetCall(cluster string, namespace string, typ string, name string) string {
	if CmdParams.Options.OcClient {
		var cmd *exec.Cmd
		token := ocGetToken(cluster)
		ns := namespace != ""
		na := name != ""
		switch {
		case ns && na:
			cmd = exec.Command("oc", "--token", token, "-n", namespace, "get", typ, name, "-o", "json")
		case ns && !na:
			cmd = exec.Command("oc", "--token", token, "-n", namespace, "get", typ, "-o", "json")
		case !ns && na:
			cmd = exec.Command("oc", "--token", token, "get", typ, name, "-o", "json")
		default:
			cmd = exec.Command("oc", "--token", token, "get", typ, "-o", "json")
		}
		jsonstr, err := cmd.Output()
		// InfoLogger.Println("JsonStr:" + string(jsonstr))
		if err != nil {
			LogError("oc get failed: " + string(jsonstr) + "Error:" + err.Error())
			exitWithError(err.Error())
		}
		return string([]byte(jsonstr))
	} else {
		var url string
		var urlpath string
		// Create a Bearer string by appending string access token
		bearer := "Bearer " + ocGetToken(cluster)
		calltyp := typ
		switch typ {
		case "images":
			urlpath = "/apis/image.openshift.io/v1/"
		case "imagestreamtags", "imagestreams", "deploymentconfigs", "namespace":
			urlpath = "/oapi/v1/namespaces/" + namespace + "/"
		case "jobs":
			urlpath = "/apis/batch/v1/namespaces/" + namespace + "/"
		case "cronjobs":
			urlpath = "/apis/batch/v1beta1/namespaces/" + namespace + "/"
		case "namespaces":
			urlpath = "/api/v1/namespaces"
			typ = ""
		default:
			urlpath = "/api/v1/namespaces/"
		}
		switch {
		case typ != "" && name != "":
			url = Clusters.Config[cluster].Url + urlpath + typ + name
		case typ != "" && name == "":
			url = Clusters.Config[cluster].Url + urlpath + typ
		default:
			url = Clusters.Config[cluster].Url + urlpath
		}
		InfoLogger.Println("call API to cluster: ", cluster, "with: ", url, "to get: ", calltyp, name, ".")
		// Create a new request using http
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			LogError("Get " + url + " failed. " + err.Error())
			return ""
		}
		// add header to the req
		req.Header.Set("Authorization", bearer)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
		// Send req using http Client
		var defaultTransport http.RoundTripper = &http.Transport{Proxy: nil}
		var client = &http.Client{}
		if CmdParams.Options.NoProxy {
			client = &http.Client{Transport: defaultTransport}
		}
		if CmdParams.Options.Socks5Proxy != "" {
			dialer, err := proxy.SOCKS5("tcp", CmdParams.Options.Socks5Proxy, nil, proxy.Direct)
			if err != nil {
				LogError("can't connect to the proxy:", err)
				// os.Exit(1)
			}
			httpTransport := &http.Transport{}
			client = &http.Client{Transport: httpTransport}
			// set our socks5 as the dialer
			httpTransport.Dial = dialer.Dial

		}
		// client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			LogError("Error on sending request.\n[ERROR] -" + err.Error())
			return ""
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			LogError("Error on reading response.\n[ERROR] -" + err.Error())
			return ""
		}
		return string([]byte(body))
	}
}
