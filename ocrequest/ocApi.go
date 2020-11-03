// Package ocrequest provides primitoives to query an oc-cluster for istags
// calculate and filter them. And provides a json output primitive.
package ocrequest

import (
	"io/ioutil"
	"log"
	"net/http"
)

// ocApiCall requests an openshift-cluster via API and return the answer as string.
func ocAPiCall(cluster string, namespace string, typ string, name string) string {
	var url string
	var urlpath string

	// Create a Bearer string by appending string access token
	bearer := "Bearer " + ocGetToken(cluster)

	switch typ {
	case "imagestreamtags", "imagestreams", "deploymentconfigs":
		urlpath = "/oapi/v1/namespaces/"
	default:
		urlpath = "/api/v1/namespaces/"
	}

	if name != "" {
		url = "https://console." + cluster + ".sf-rz.de:8443" + urlpath +
			namespace + "/" + typ + "/" + name
	} else {
		url = "https://console." + cluster + ".sf-rz.de:8443" + urlpath +
			namespace + "/" + typ
	}

	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
	}

	// add header to the req
	req.Header.Set("Authorization", bearer)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	return string([]byte(body))
}
