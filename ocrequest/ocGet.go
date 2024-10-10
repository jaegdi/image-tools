// Package ocrequest provides primitoives to query an oc-cluster for istags
// calculate and filter them. And provides a json output primitive.
package ocrequest

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"time"

	"golang.org/x/net/proxy"
)

// RZ Root CA
var (
	certs = `-----BEGIN CERTIFICATE-----
MIIHNjCCBR6gAwIBAgIJAN58FKMMGU5JMA0GCSqGSIb3DQEBCwUAMIHCMQswCQYD
VQQGEwJERTEMMAoGA1UECBMDSEVTMRIwEAYDVQQHEwlXaWVzYmFkZW4xGjAYBgNV
BAoTEVNDSFVGQSBIb2xkaW5nIEFHMQswCQYDVQQLEwJJVDFAMD4GA1UEAxM3U0NI
VUZBIEhvbGRpbmcgQUcgLSBaZXJ0aWZpemllcnVuZ3NzdGVsbGUgUmVjaGVuemVu
dHJ1bTEmMCQGCSqGSIb3DQEJARYXcnotemVydGlmaWthdEBzY2h1ZmEuZGUwHhcN
MTQxMDI4MTMxMjI1WhcNMzQxMDIzMTMxMjI1WjCBwjELMAkGA1UEBhMCREUxDDAK
BgNVBAgTA0hFUzESMBAGA1UEBxMJV2llc2JhZGVuMRowGAYDVQQKExFTQ0hVRkEg
SG9sZGluZyBBRzELMAkGA1UECxMCSVQxQDA+BgNVBAMTN1NDSFVGQSBIb2xkaW5n
IEFHIC0gWmVydGlmaXppZXJ1bmdzc3RlbGxlIFJlY2hlbnplbnRydW0xJjAkBgkq
hkiG9w0BCQEWF3J6LXplcnRpZmlrYXRAc2NodWZhLmRlMIICIjANBgkqhkiG9w0B
AQEFAAOCAg8AMIICCgKCAgEA0pEyG/HMVfWznJ4mC1mRERMTcLBnhuKix5ViyVxm
x6QLFzwFjrAyKqdQ9E32L7Zu089v5iDX9IvYu0Spj5bIbLcfmy+jsN5QHKELvTZ0
AYjN+mHIygQiIP0/8hEVcmpGNC6hjcpOnd7b7xVIY6YGR67hLgwZjiFg69ln5Eep
wqHYORQsq0iqvQVpg6SzOtxZLvhCDIrhQMAaTcVNIbbsJ2Lg5TVxQEOYIXQEMlQH
d5rMa0SAXf/tv4grg4FrT5cV0byKplS/Kn/CfxPMdBuK27ISH4eP97/aZv2NeLmT
MwB7Sse7ZATb8+G2p5lJxE1iov8KqEFCd9F0w9Wu7Nlw7V9+982/Jrfut7VLT02l
e38wCHT+AJe7d1whuVZ3WEzPx/EX1RBSe/A8CYmaw1agtGCc/Ue7QwxHbiZXqs0D
mKgN/GhKfVQubUKEZyuHki6MDJLSRiDB4ANXy/xvWwBeoovNgcwX/CNT4wO42sLL
+S+GyqGAQfmJn9mMmbSeGSIu5+/aYzWBL7LhIZuIJeQaQjZ4Okmjc+ktV+a21Qj6
AfdGZKgS62JscyECZVunCL7PjUGIzFJ9vIlEhxdItL3+5JQH7AmOU9RiaY97YLHq
fSDxuIbI0Hl0qA9XrLxQ6C85oLNJfZNIhL4LtP5nbHnJa7M2m+OkW6eM5tvCswyf
cfECAwEAAaOCASswggEnMB0GA1UdDgQWBBSwoe2F3UMUVG60Zqq18Rv+vHQLxDCB
9wYDVR0jBIHvMIHsgBSwoe2F3UMUVG60Zqq18Rv+vHQLxKGByKSBxTCBwjELMAkG
A1UEBhMCREUxDDAKBgNVBAgTA0hFUzESMBAGA1UEBxMJV2llc2JhZGVuMRowGAYD
VQQKExFTQ0hVRkEgSG9sZGluZyBBRzELMAkGA1UECxMCSVQxQDA+BgNVBAMTN1ND
SFVGQSBIb2xkaW5nIEFHIC0gWmVydGlmaXppZXJ1bmdzc3RlbGxlIFJlY2hlbnpl
bnRydW0xJjAkBgkqhkiG9w0BCQEWF3J6LXplcnRpZmlrYXRAc2NodWZhLmRlggkA
3nwUowwZTkkwDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAgEAXA4pUTv3
OMNUs3nWp5x0GmpX8DrZj7/GN0iKAR9+gJNlPbLqHqF1T/wqS935hPoQI+ffhMua
92jQ3P1cy5RL/8514joEVnGWQUs9iccOSfSRag0nqBRSbNoFcBnI8gDDx9jDZDsx
bEsQb1juL8WT3EfHIIhzEK7XfMF2ikzVDbsDq9HV/ZX6ofA9w6XTfeq7rQ0FRSuZ
3S14RTeirJdM3ngJ7N6/9/U1QI8suqDktx9fL0vLCspBSj99N0iS0lK8GYfwKB6t
fT/f4xdjTQyfVbu8zTDZhZQerHAUOkrNaK+vzZ+vmkfKy4DxvXJwzglCFTVTLnSW
KTPQ4AOD+GotfDPwPoFEGX1lSzW9KYOzcn7G5KODhsB1Cjc1k095z3zJqoGq2F0q
F4dpWs5n0TKbuxbTgV6IJX8yHZ4mthQjUX0vItZs9DB5SoWMBQwbA+WUQ22R1fat
UH2kwb57+6FiSFUyAhgPM78t3gTinYg7cmQwxmzos82bu1MW01375ZWmtTKA2Lvu
taFKcEs51Fa5N9wywWqRcmj236PwnEwhbBRNujnB/vsuZBoqve6cBSgjzOP/B7OB
/c06p8Yc4ROSQV2xPXh1Fw6hP9OKqkY1FE2r4xF8fEhxsnCRKFfctxSParmXkPvt
rf0DbWmzK8OwXWu+x28rj88N8pO1wceTRcY=
-----END CERTIFICATE-----`
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
	DebugMsg(cmd)
	jsonstr, err := cmd.Output()
	if err != nil {
		ExitWithError("oc get failed:", string(jsonstr), "Error:", err)
	}
	return []byte(jsonstr)
}

// ocApiCall requests an OpenShift cluster via API and returns the response as a byte slice.
//
// Parameters:
// - cluster: The cluster name.
// - namespace: The namespace name.
// - resourceType: The type of the resource to request.
// - resourceName: The name of the resource to request.
//
// Returns:
// - A byte slice containing the API response.
func ocApiCall(cluster T_clName, namespace T_nsName, resourceType string, resourceName string) []byte {
	baseURL := string(Clusters.Config[cluster].Url)
	url := buildURL(cluster, baseURL, namespace, resourceType, resourceName)

	config := createTLSConfig()
	bearerToken := "Bearer " + ocGetToken(cluster)

	req, err := createRequest(url, bearerToken)
	if err != nil {
		ErrorMsg("Get " + url + " failed. " + err.Error())
		return []byte("")
	}

	client := createHTTPClient(config)
	resp, err := client.Do(req)
	if err != nil {
		ErrorMsg("Error on sending request. " + err.Error())
		return []byte("")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ErrorMsg("Error on reading response. " + err.Error())
		return []byte("")
	}
	DebugMsg("body:", string(body))
	return body
}

// buildURL constructs the URL for the API request based on the resource type and name.
//
// Parameters:
// - baseURL: The base URL of the cluster.
// - namespace: The namespace name.
// - resourceType: The type of the resource to request.
// - resourceName: The name of the resource to request.
//
// Returns:
// - A string containing the constructed URL.
func buildURL(cluster T_clName, baseURL string, namespace T_nsName, resourceType string, resourceName string) string {
	var url string
	switch resourceType {
	case "images", "is", "imagestreams", "istags", "imagestreamtags":
		url = baseURL + "/apis/image.openshift.io/v1/" + resourceType
	case "dc", "deploymentconfigs":
		url = baseURL + "/apis/apps.openshift.io/v1/namespaces/" + namespace.str() + "/deploymentconfigs"
	case "deployments":
		url = baseURL + "/apis/apps/v1/namespaces/" + namespace.str() + "/deployments"
	case "jobs", "cronjobs":
		url = baseURL + "/apis/batch/v1/namespaces/" + namespace.str() + "/" + resourceType
	case "builds", "buildconfigs":
		url = baseURL + "/apis/build.openshift.io/v1/namespaces/" + namespace.str() + "/" + resourceType
	case "daemonsets", "replicasets", "pods":
		url = baseURL + "/apis/apps/v1/namespaces/" + namespace.str() + "/" + resourceType
	case "pipelines", "tasks", "clustertasks":
		url = baseURL + "/apis/tekton.dev/v1beta1/namespaces/" + namespace.str() + "/" + resourceType
	default:
		if namespace.str() != "" {
			url = baseURL + "/apis/apps/v1/namespaces/" + namespace.str() + "/" + resourceType
		} else {
			url = baseURL + "/api/v1/" + resourceType
		}
	}

	if resourceName != "" {
		url += "/" + resourceName
	}

	DebugMsg("call API to cluster: ", cluster, "with: ", url, "to get: ", resourceType, resourceName, ".")
	return url
}

// createTLSConfig creates a TLS configuration with the system cert pool and additional certs.
//
// Returns:
// - A pointer to a tls.Config containing the TLS configuration.
func createTLSConfig() *tls.Config {
	// Get the SystemCertPool, continue with an empty pool on error
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	// Append our cert to the system pool
	if ok := rootCAs.AppendCertsFromPEM([]byte(certs)); !ok {
		ErrorMsg("No certs appended, using system certs only")
	}

	// Trust the augmented cert pool in our client
	return &tls.Config{
		InsecureSkipVerify: CmdParams.Options.InsecureSSL,
		RootCAs:            rootCAs,
	}
}

// createRequest creates an HTTP GET request with the specified URL and bearer token.
//
// Parameters:
// - url: The URL for the request.
// - bearerToken: The bearer token for authorization.
//
// Returns:
// - A pointer to an http.Request containing the created request.
// - An error if the request creation fails.
func createRequest(url string, bearerToken string) (*http.Request, error) {
	// Create a new request using http
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	// Add headers to the request
	req.Header.Set("Authorization", bearerToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

// createHTTPClient creates an HTTP client with the specified TLS configuration and proxy settings.
//
// Parameters:
// - config: A pointer to a tls.Config containing the TLS configuration.
//
// Returns:
// - A pointer to an http.Client containing the created client.
func createHTTPClient(config *tls.Config) *http.Client {
	// Initialize the HTTP transport with the provided TLS configuration and a default timeout
	var transport http.RoundTripper = &http.Transport{
		TLSClientConfig: config,
		DialContext: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).DialContext,
	}

	// Check if the NoProxy option is enabled
	if CmdParams.Options.NoProxy {
		// If NoProxy is enabled, create a transport without a proxy and with a shorter timeout
		transport = &http.Transport{
			Proxy:           nil,
			TLSClientConfig: config,
			DialContext: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).DialContext,
		}
	}

	// Check if a SOCKS5 proxy is specified
	if CmdParams.Options.Socks5Proxy != "" {
		// Create a SOCKS5 dialer
		dialer, err := proxy.SOCKS5("tcp", CmdParams.Options.Socks5Proxy, nil, proxy.Direct)
		if err != nil {
			// If there is an error creating the SOCKS5 dialer, return a client with a short timeout
			return &http.Client{Timeout: 5 * time.Second}
		}
		// Define a dial context function using the SOCKS5 dialer
		dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		}
		// Create a transport with the SOCKS5 dial context and shorter timeouts
		transport = &http.Transport{
			DialContext:           dialContext,
			TLSClientConfig:       config,
			TLSHandshakeTimeout:   5 * time.Second,
			ExpectContinueTimeout: 5 * time.Second,
		}
	}

	// Return the HTTP client with the configured transport and a short timeout
	return &http.Client{
		Transport: transport,
		Timeout:   5 * time.Second,
	}
}

// checkCache checks if the cache file exists and is not older than 1 minute
func checkCache(tmpdir string, cluster T_clName, namespace T_nsName, typ string, name string) (string, bool) {
	filename := tmpdir + "/" + "cache_" + string(cluster) + "_" + string(namespace) + "_" + typ + "_" + name + ".tmp"
	//  no cache in serverMode
	if CmdParams.Options.ServerMode {
		return filename, false
	}
	//  dir not exist
	if _, err := os.Stat(tmpdir); os.IsNotExist(err) {
		err := os.MkdirAll(tmpdir, 0755)
		if err != nil {
			DebugMsg("failed to create cache dir", err)
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
		DebugMsg("Cache Age:", duration.Minutes())
		return filename, false
	}
	return filename, true
}

// writeCache writes the connntent to the cache file
func writeCache(tmpdir string, filename string, content []byte) {
	// no cache in serverMode
	if !CmdParams.Options.ServerMode {
		err := os.WriteFile(tmpdir+"/"+filename, content, 0644)
		if err != nil {
			DebugMsg("Writing cache file failed", err)
		}
	}
}

// readCache read from the cache file
func readCache(filename string) []byte {
	content, _ := os.ReadFile(filename)
	return content
}

// ocGetCall requests an openshift-cluster via API or oc-client and return the answer as string.
func ocGetCall(cluster T_clName, namespace T_nsName, typ string, name string) string {
	tmpdir := "/tmp/tmp-report-istags"
	VerifyMsg("cluster:", cluster, "| namespace:", namespace, "| type:", typ, "| name:", name)
	var content []byte
	filename, cacheOk := checkCache(tmpdir, cluster, namespace, typ, name)
	if !cacheOk {
		DebugMsg("Request Openshift for:", filename)
		if CmdParams.Options.OcClient {
			content = ocClientCall(cluster, namespace, typ, name)
		} else {
			content = ocApiCall(cluster, namespace, typ, name)
		}
		writeCache(tmpdir, filename, content)
	} else {
		DebugMsg("Use Cache for:", filename)
		content = readCache(filename)
	}
	return string(content)
}
