package ocrequest

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/proxy"
)

func getHttpAnswer(url string) []byte {
	// Get the SystemCertPool, continue with an empty pool on error
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	// Append our cert to the system pool
	if ok := rootCAs.AppendCertsFromPEM([]byte(certs)); !ok {
		ErrorLogger.Println("No certs appended, using system certs only")
	}

	// Trust the augmented cert pool in our client
	config := &tls.Config{
		InsecureSkipVerify: CmdParams.Options.InsecureSSL,
		RootCAs:            rootCAs,
	}

	// Create a new request using http
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		ErrorLogger.Println("Get " + url + " failed. " + err.Error())
		return []byte("")
	}
	// add header to the req
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	// Send req using http Client
	var noproxyTransport http.RoundTripper = &http.Transport{Proxy: nil, TLSClientConfig: config}
	var defaultTransport = &http.Transport{TLSClientConfig: config}
	var client = &http.Client{Transport: defaultTransport}
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
		httpTransport := &http.Transport{TLSClientConfig: config}
		client = &http.Client{Transport: httpTransport}
		// set our socks5 as the dialer
		httpTransport.Dial = dialer.Dial

	}
	// client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ErrorLogger.Println("Error on sending request." + err.Error())
		return []byte("")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ErrorLogger.Println("Error on reading response." + err.Error())
		return []byte("")
	}
	return []byte(body)
}
