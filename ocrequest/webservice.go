package ocrequest

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "image-tool/docs" // Import the generated docs

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Webservice IMAGE-TOOL API
// @version 1.0
// @description This is a sample server for IMAGE-TOOL.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @query.collection.format multi

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
	InfoMsg("Starting server on port 8080")
	http.HandleFunc("/", handleDocumentation)
	http.HandleFunc("/execute", handleExecute)
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)
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
	docPage := GetDocPage()
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
// @Summary Execute a command
// @Description Executes a command based on the provided query parameters.
// @Tags execute
// @Accept  json
// @Produce  json
// @Param   family     query    string     true        "The family parameter (required for 'is_tag_used')"
// @Param   kind       query    string     false       "The kind of operation to perform. Valid values are 'used', 'is_tag_used', 'unused', 'istag', 'is', 'image', 'all'. Default is 'is_tag_used'"
// @Param   cluster    query    string     false       "The cluster parameter"
// @Param   tagname    query    string     false        "The tagname parameter (required for 'is_tag_used')"
// @Param   namespace  query    string     false       "The namespace parameter"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {string} string "Missing or invalid parameters"
// @Router /execute [get]
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
	_, html, err := initializeCmdParams(family, kind, cluster, filter_tagname, filter_namespace)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		ErrorMsg("Error:", err)
		return
	}
	InfoMsg("Parameters:", GetJsonFromMap(CmdParams))
	// Initialisiere den Servermodus mit den Kommando-Parametern
	InitServerMode(CmdParams)

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
	// cmdParams := T_flags{}
	CmdParams.Family = T_familyName(family)
	CmdParams.Cluster = T_clName(cluster).list()
	CmdParams.Filter.Tagname = T_tagName(tagname)
	CmdParams.Filter.Namespace = T_nsName(namespace)
	html := true
	CmdParams.Output.All = false
	CmdParams.Output.Image = false
	CmdParams.Output.Is = false
	CmdParams.Output.Istag = false
	CmdParams.Output.UnUsed = false
	CmdParams.Output.Used = false

	switch kind {
	case "used":
		CmdParams.Output.Used = true
	case "is_tag_used":
		CmdParams.Output.Used = true
		html = false
	case "unused":
		CmdParams.Output.UnUsed = true
	case "istag":
		CmdParams.Output.Istag = true
	case "is":
		CmdParams.Output.Is = true
	case "image":
		CmdParams.Output.Image = true
	case "all":
		CmdParams.Output.All = true
	default:
		return CmdParams, html, fmt.Errorf("Invalid kind parameter")
	}
	CmdParams.Html = html
	return CmdParams, html, nil
}

// processResults processes the results and writes the response.
func processResults(w http.ResponseWriter, result T_completeResults, family string, html bool, kind, tagname string) {
	switch kind {
	case "is_tag_used":
		{
			tagIsUsed := len(result.UsedIstags) > 0
			response := map[string]interface{}{
				"TagName":   tagname,
				"TagIsUsed": tagIsUsed,
			}
			json.NewEncoder(w).Encode(response)
		}
	default:
		htmldata := GetTextTableFromMap(result, T_familyName(family))
		if html {
			w.Header().Set("Content-Type", "text/html")
		} else {
			w.Header().Set("Content-Type", "application/json")
		}
		w.Write([]byte(htmldata))
	}
}
