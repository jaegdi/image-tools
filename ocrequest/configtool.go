package ocrequest

import (
	"io/ioutil"
	"net/http"

	"golang.org/x/net/proxy"
)

func GetClusters(configtoolUrl string) interface{} {
	url := configtoolUrl + "/clusters?validate=true"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		LogError("Get " + url + " failed. " + err.Error())
		return []byte("")
	}
	// add header to the req
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	// Send req using http Client
	var noproxyTransport http.RoundTripper = &http.Transport{Proxy: nil}
	var client = &http.Client{}
	// NO_PROXY handling
	if CmdParams.Options.NoProxy {
		client = &http.Client{Transport: noproxyTransport}
	}
	// Socks5 proxy handling
	if CmdParams.Options.Socks5Proxy != "" {
		dialer, err := proxy.SOCKS5("tcp", CmdParams.Options.Socks5Proxy, nil, proxy.Direct)
		if err != nil {
			exitWithError("can't connect to the proxy:", err)
		}
		httpTransport := &http.Transport{}
		client = &http.Client{}
		// set our socks5 as the dialer
		httpTransport.Dial = dialer.Dial
	}

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
