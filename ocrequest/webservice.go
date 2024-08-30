package ocrequest

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// T_flags, T_familyName, T_clName, T_tagName, T_nsName, CmdParams, InitServerMode, CmdlineMode, InfoMsg, ErrorMsg, VerifyMsg, GetJsonFromMap, GetHtmlTableFromMap
// sollten hier definiert oder importiert werden

// StartServer starts the HTTP server and handles the "/execute" endpoint.
// It expects the "family" and "tagname" parameters to be provided in the URL query.
// If any of the parameters is missing, it returns a HTTP 400 Bad Request error.
// Otherwise, it sets the CmdParams.Family and CmdParams.Filter.Tagname based on the provided values.
// It also sets CmdParams.Output.Used to true.
// The server listens on port 8080 and returns the result of CmdlineMode as a JSON response.
// The response content type is set to "application/json".
// handleDocumentation serves the documentation page.
// StartServer starts the HTTP server and handles incoming requests.
func StartServer() {
	http.HandleFunc("/", handleDocumentation)
	http.HandleFunc("/execute", handleExecute)
	http.ListenAndServe(":8080", nil)
}

// handleDocumentation serves the documentation page for the webservice.
// It provides information on how to use the webservice, including available endpoints,
// required query parameters, and example usage.
//
// The documentation is served as an HTML page with the following structure:
// - A title and introductory text
// - A list of endpoints with descriptions
// - A list of query parameters for each endpoint
// - A list of possible responses
// - An example usage of the endpoint
//
// The HTML content is written directly to the response writer.
//
// Example usage:
// When a user navigates to the root URL ("/"), this function will be called
// and the documentation page will be displayed.
//
// Parameters:
// - w: The http.ResponseWriter to write the HTML content to.
// - r: The http.Request object (not used in this function).
func handleDocumentation(w http.ResponseWriter, r *http.Request) {
	docPage := `
<!DOCTYPE html>
<html>
<head>
    <title>Webservice IMAGE-TOOL Documentation</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        h1 { color: #333; }
        p { margin: 10px 0; }
        code { background-color: #f4f4f4; padding: 2px 4px; border-radius: 4px; }
    </style>
</head>
<body>
    <h1>Webservice IMAGE-TOOL Documentation</h1>
    <p>Welcome to the webservice documentation page of the image-tool. Below you will find information on how to use the webservice.</p>
    <h2>Endpoints</h2>
	<h3>GET /</h3>
	<p>This endpoint show this documentation</p>
    <h3>GET /execute</h3>
    <p>This endpoint executes a command based on the provided query parameters.</p>
    <p><strong>Query Parameters:</strong></p>
    <ul>
        <li><code>family</code> (required for <code>is_tag_used</code>): The family parameter.</li>
        <li><code>kind</code>: The kind of operation to perform. Valid values are <code>is_tag_used</code>.
		<br><pre>       The default is <code>is_tag_used</code></pre></li>
        <li><code>cluster</code>: The cluster parameter (is not necessary for kind <code>is_tag_used</code>): eg. cid-scp0, ... or pro-scp0</li>
        <li><code>tagname</code> (required for <code>is_tag_used</code>): The tagname parameter.</li>
        <li><code>namespace</code>: The namespace parameter.</li>
    </ul>
    <p><strong>Responses:</strong></p>
    <ul>
        <li><code>200 OK</code>: The command was executed successfully. The response is in JSON format.
		<br><pre>       eg.: <code>{"TagIsUsed":true,"TagName":"pkp-3.19.0-build-3"}</code>
		<br>       eg.: <code>{"TagIsUsed":false,"TagName":"pkp-x-not-there"}</code></pre></li>

        <li><code>400 Bad Request</code>: Missing or invalid parameters.</li>
    </ul>
    <p>Example usage:</p>
    <pre><code>GET /execute?family=pkp&kind=is_tag_used&tagname=pkp-3.19.0-build-3</code></pre>
</body>
</html>
`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(docPage))
}

// handleExecute handles the /execute endpoint.
// This function processes incoming HTTP requests to the /execute endpoint,
// executes the specified command based on the provided query parameters,
// and returns the result in either HTML or JSON format.
//
// Query Parameters:
//   - family: The family parameter (required for "is_tag_used").
//   - kind: The kind of operation to perform. Valid values are "used", "is_tag_used",
//     "unused", "istag", "is", "image", "all". Default is "is_tag_used".
//   - cluster: The cluster parameter.
//   - tagname: The tagname parameter (required for "is_tag_used").
//   - namespace: The namespace parameter.
//
// Responses:
//   - 200 OK: The command was executed successfully. The response is in JSON format
//     or HTML format based on the kind parameter.
//   - 400 Bad Request: Missing or invalid parameters.
//
// Example usage:
// GET /execute?family=exampleFamily&kind=used&tagname=exampleTag
//
// Parameters:
// - w: The http.ResponseWriter to write the response to.
// - r: The http.Request object containing the query parameters.
func handleExecute(w http.ResponseWriter, r *http.Request) {
	// Extrahiere die Query-Parameter aus der URL
	family := r.URL.Query().Get("family")
	kind := r.URL.Query().Get("kind")
	cluster := r.URL.Query().Get("cluster")
	filter_tagname := r.URL.Query().Get("tagname")
	filter_namespace := r.URL.Query().Get("namespace")

	// Logge eine Informationsnachricht über die neue Anfrage
	InfoMsg("--------------  New request  --------------")

	// Validierung der erforderlichen Parameter
	if err := validateParams(family, kind, filter_tagname); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		ErrorMsg("Error:", err)
		return
	}

	// Setze den Standardwert für 'kind', falls nicht angegeben
	if kind == "" {
		kind = "is_tag_used"
	}

	// Initialisiere die Kommando-Parameter basierend auf den Abfrageparametern
	cmdParams, html, err := initializeCmdParams(family, kind, cluster, filter_tagname, filter_namespace)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		ErrorMsg("Error:", err)
		return
	}

	// Initialisiere den Servermodus mit den Kommando-Parametern
	InitServerMode(cmdParams)

	// Logge die Details der Anfrage
	InfoMsg("family:", family, "| kind:", kind, "| tagname:", filter_tagname)
	VerifyMsg("\nCmdParams Options:", GetJsonFromMap(CmdParams.Options), "Output:", GetJsonFromMap(CmdParams.Output))

	// Führe den Befehl aus und erhalte das Ergebnis
	result := CmdlineMode()
	VerifyMsg("\nCmdParams Result:", GetJsonFromMap(result))

	// Verarbeite die Ergebnisse und schreibe die Antwort
	processResults(w, result, family, html, kind, filter_tagname)
}

// validateParams validates the required query parameters.
func validateParams(family, kind, tagname string) error {
	if kind == "is_tag_used" && (family == "" || tagname == "") {
		return fmt.Errorf("Missing parameters: family and tagname are required for kind 'is_tag_used'")
	}
	return nil
}

// initializeCmdParams initializes the command parameters based on the query parameters.
func initializeCmdParams(family, kind, cluster, tagname, namespace string) (T_flags, bool, error) {
	cmdParams := T_flags{}
	cmdParams.Family = T_familyName(family)
	cmdParams.Cluster = T_clName(cluster).list()
	html := true

	switch kind {
	case "used":
		cmdParams.Output.Used = true
	case "is_tag_used":
		cmdParams.Output.Used = true
		html = false
	case "unused":
		cmdParams.Output.UnUsed = true
	case "istag":
		cmdParams.Output.Istag = true
	case "is":
		cmdParams.Output.Is = true
	case "image":
		cmdParams.Output.Image = true
	case "all":
		cmdParams.Output.All = true
	default:
		return cmdParams, html, fmt.Errorf("Invalid kind parameter")
	}

	cmdParams.Filter.Tagname = T_tagName(tagname)
	cmdParams.Filter.Namespace = T_nsName(namespace)
	return cmdParams, html, nil
}

// processResults processes the results and writes the response.
func processResults(w http.ResponseWriter, result T_completeResults, family string, html bool, kind, tagname string) {
	if html {
		htmldata, err := GetHtmlTableFromMap(result, T_familyName(family))
		if err != nil {
			ErrorMsg("Error creating template:", err)
		}
		w.Write([]byte(htmldata))
	} else {
		w.Header().Set("Content-Type", "application/json")
		if kind == "is_tag_used" {
			tagIsUsed := len(result.UsedIstags) > 0
			response := map[string]interface{}{
				"TagName":   tagname,
				"TagIsUsed": tagIsUsed,
			}
			json.NewEncoder(w).Encode(response)
		} else {
			json.NewEncoder(w).Encode(result)
		}
	}
}
