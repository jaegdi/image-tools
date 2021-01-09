// Package ocrequest provides primitoives to query an oc-cluster for istags
// calculate and filter them. And provides a json output primitive.
package ocrequest

import (
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"

	"golang.org/x/net/proxy"
)

// ocClientCall requests an openshift-cluster via oc clioent and return the answer as string.
func ocClientCall(cluster T_clName, namespace T_nsName, typ string, name string) []byte {
	var cmd *exec.Cmd
	token := ocGetToken(cluster)
	ns := namespace != ""
	na := name != ""
	switch {
	case ns && na:
		cmd = exec.Command("oc", "--token", token, "-n", namespace.str(), "get", typ, name, "-o", "json")
	case ns && !na:
		cmd = exec.Command("oc", "--token", token, "-n", namespace.str(), "get", typ, "-o", "json")
	case !ns && na:
		cmd = exec.Command("oc", "--token", token, "get", typ, name, "-o", "json")
	default:
		cmd = exec.Command("oc", "--token", token, "get", typ, "-o", "json")
	}
	LogDebug(cmd)
	jsonstr, err := cmd.Output()
	if err != nil {
		exitWithError("oc get failed:", string(jsonstr), "Error:", err)
	}
	return []byte(jsonstr)
}

// ocApiCall requests an openshift-cluster via API and return the answer as string.
func ocApiCall(cluster T_clName, namespace T_nsName, typ string, name string) []byte {
	var url string
	var urlpath string
	// Create a Bearer string by appending string access token
	bearer := "Bearer " + ocGetToken(cluster)
	calltyp := typ
	switch typ {
	case "images":
		urlpath = "/apis/image.openshift.io/v1/"
	case "imagestreamtags", "imagestreams", "deploymentconfigs", "namespace":
		urlpath = "/oapi/v1/namespaces/" + namespace.str() + "/"
	case "jobs":
		urlpath = "/apis/batch/v1/namespaces/" + namespace.str() + "/"
	case "cronjobs":
		urlpath = "/apis/batch/v1beta1/namespaces/" + namespace.str() + "/"
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
	LogDebug("call API to cluster: ", cluster, "with: ", url, "to get: ", calltyp, name, ".")
	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		LogError("Get " + url + " failed. " + err.Error())
		return []byte("")
	}
	// add header to the req
	req.Header.Set("Authorization", bearer)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	// Send req using http Client
	var defaultTransport http.RoundTripper = &http.Transport{Proxy: nil}
	var client = &http.Client{}
	// NO_PROXY handling
	if CmdParams.Options.NoProxy {
		client = &http.Client{Transport: defaultTransport}
	}
	// Socks5 proxy handling
	if CmdParams.Options.Socks5Proxy != "" {
		dialer, err := proxy.SOCKS5("tcp", CmdParams.Options.Socks5Proxy, nil, proxy.Direct)
		if err != nil {
			exitWithError("can't connect to the proxy:", err)
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
		return []byte("")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		LogError("Error on reading response.\n[ERROR] -" + err.Error())
		return []byte("")
	}
	return []byte(body)
}

func checkCache(tmpdir string, cluster T_clName, namespace T_nsName, typ string, name string) (string, bool) {
	filename := tmpdir + "/" + "cache_" + string(cluster) + "_" + string(namespace) + "_" + typ + "_" + name + ".tmp"
	//  dir not exist
	if _, err := os.Stat(tmpdir); os.IsNotExist(err) {
		err := os.MkdirAll(tmpdir, 0755)
		if err != nil {
			LogDebug("failed to create cache dir", err)
		}
		return filename, false
	}
	info, err := os.Stat(filename)
	// file not exist
	if err != nil {
		return filename, false
	}
	duration := time.Since(info.ModTime())
	// file too old
	if duration.Minutes() > float64(1.0) {
		LogDebug("Cache Age:", duration.Minutes())
		return filename, false
	}
	return filename, true
}

func writeCache(tmpdir string, filename string, content []byte) {
	err := ioutil.WriteFile(filename, content, 0644)
	if err != nil {
		LogDebug("Writing cache file failed", err)
	}
}

func readCache(filename string) []byte {
	content, _ := ioutil.ReadFile(filename)
	return content
}

// ocGetCall requests an openshift-cluster via API or oc-client and return the answer as string.
func ocGetCall(cluster T_clName, namespace T_nsName, typ string, name string) string {
	tmpdir := "/tmp/tmp-report-istags"
	var content []byte
	filename, cacheOk := checkCache(tmpdir, cluster, namespace, typ, name)
	if !cacheOk {
		LogDebug("Request Openshift for:", filename)
		if CmdParams.Options.OcClient {
			content = ocClientCall(cluster, namespace, typ, name)
		} else {
			content = ocApiCall(cluster, namespace, typ, name)
		}
		writeCache(tmpdir, filename, content)
	} else {
		LogDebug("Use Cache for:", filename)
		content = readCache(filename)
	}
	return string(content)
}
