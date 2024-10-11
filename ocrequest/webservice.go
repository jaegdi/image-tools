package ocrequest

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "image-tool/docs" // Import the generated docs

	httpSwagger "github.com/swaggo/http-swagger"
)

var initialCmdParams T_flags

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

// @Description StartServer starts the HTTP server and handles the "/execute" endpoint.
// @Description It expects the "family" and "tagname" parameters to be provided in the URL query.
// @Description If any of the parameters is missing, it returns a HTTP 400 Bad Request error.
// @Description Otherwise, it sets the CmdParams.Family and CmdParams.Filter.Tagname based on the provided values.
// @Description It also sets CmdParams.Output.Used to true.
// @Description The server listens on port 8080 and returns the result of CmdlineMode as a JSON response.
// @Description The response content type is set to "application/json".
// @Description handleDocumentation serves the documentation page.
// @Description StartServer starts the HTTP server and handles incoming requests.
func StartServer() {
	InfoMsg("Starting server on port 8080")
	http.HandleFunc("/", handleDocumentation)
	http.HandleFunc("/query", handleQuery)
	http.HandleFunc("/is-tag-used", handleIsTagUsed)
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	http.ListenAndServe(":8080", nil)
}

// @Description handleDocumentation serves the documentation page for the webservice.
// @Description It provides information on how to use the webservice, including available endpoints,
// @Description required query parameters, and example usage.
// @Description
// @Description The documentation is served as an HTML page with the following structure:
// @Description - A title and introductory text
// @Description - A list of endpoints with descriptions
// @Description - A list of query parameters for each endpoint
// @Description - A list of possible responses
// @Description - An example usage of the endpoint
// @Description
// @Description The HTML content is written directly to the response writer.
// @Description
// @Description Example usage:
// @Description When a user navigates to the root URL ("/"), this function will be called
// @Description and the documentation page will be displayed.
// @Description
// @Description Parameters:
// @Description - w: The http.ResponseWriter to write the HTML content to.
// @Description - r: The http.Request object (not used in this function).
// @Summary Show documentation
// @Description
// @Tags documentation
// @Accept  html
// @Produce  html
// @Router / [get]
func handleDocumentation(w http.ResponseWriter, r *http.Request) {
	docPage := GetDocPage()
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(docPage))
}

// @Summary Execute a query
// @Description handleQuery handles the /query endpoint.
// @Description This function processes incoming HTTP requests to the /query endpoint,
// @Description executes the specified command based on the provided query parameters,
// @Description and returns the result in either HTML or JSON format.
// @Description
// @Description Query Parameters:
// @Description   - family:    The family parameter (required for "is_tag_used").
// @Description   - kind:      The kind of operation to perform. Valid values are "used", "is_tag_used",
// @Description                "unused", "istag", "is", "image", "all". Default is "is_tag_used".
// @Description   - tagname:   The tagname parameter to filter the istags by this tagname.
// @Description   - cluster:   The cluster parameter.
// @Description   - namespace: The namespace parameter.
// @Description
// @Description Responses:
// @Description   - 200 OK: The command was executed successfully. The response is in JSON format
// @Description             or HTML format based on the kind parameter.
// @Description   - 400 Bad Request: Missing or invalid parameters.
// @Description
// @Description Example usage:
// @Description GET /query?family=exampleFamily&kind=used&tagname=exampleTag
// @Description
// @Description The result is a HTML-Table with the queried items and theirs details
// @Description
// @Description This table can be downloaded by the "Download as Excel"
// @Description
// @Description The table is presented with a filter function and
// @Description a sort function for each column in the table.
// @Description
// @Tags query
// @Accept  json
// @Produce  text/html
// @Param   family     query    string     true        "The family parameter (required for 'is_tag_used')"
// @Param   kind       query    string     false       "The kind of operation to perform. Valid values are 'used', 'is_tag_used', 'unused', 'istag', 'is', 'image', 'all'. Default is 'is_tag_used'"
// @Param   tagname    query    string     false        "The tagname parameter (required for 'is_tag_used')"
// @Param   cluster    query    string     false       "The cluster parameter"
// @Param   namespace  query    string     false       "The namespace parameter"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {string} string "Missing or invalid parameters"
// @Router /query [get]
// Parameters:
// - w: The http.ResponseWriter to write the response to.
// - r: The http.Request object containing the query parameters.
func handleQuery(w http.ResponseWriter, r *http.Request) {
	// Extrahiere die Query-Parameter aus der URL
	family := r.URL.Query().Get("family")
	kind := r.URL.Query().Get("kind")
	cluster := r.URL.Query().Get("cluster")
	filter_namespace := r.URL.Query().Get("namespace")
	filter_tagname := r.URL.Query().Get("tagname")

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

// @Description Executes a query based on the provided query parameters.
// @Description handleIsTagUsed handles the /execute endpoint.
// @Description This function processes incoming HTTP requests to the /execute endpoint,
// @Description executes the specified command based on the provided query parameters,
// @Description and returns the result in either HTML or JSON format.
// @Description
// @Description Query Parameters:
// @Description   - family: The family parameter (required for "is_tag_used").
// @Description   - tagname: The tagname parameter (required for "is_tag_used").
// @Description
// @Description Responses:
// @Description   - 200 OK: The command was executed successfully. The response is in JSON format
// @Description     or HTML format based on the kind parameter.
// @Description   - 400 Bad Request: Missing or invalid parameters.
// @Description
// @Description Example usage:
// @Description GET /query?family=exampleFamily&tagname=exampleTag
// @Description
// @Summary Check if a image tag is used somewhere in the clusters.
// @Tags query
// @Accept  json
// @Produce  json
// @Param   family     query    string     true        "The family parameter (required for 'is_tag_used')"
// @Param   tagname    query    string     false        "The tagname parameter (required for 'is_tag_used')"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {string} string "Missing or invalid parameters"
// @Router /is-tag-used [get]
// Parameters:
// - w: The http.ResponseWriter to write the response to.
// - r: The http.Request object containing the query parameters.
func handleIsTagUsed(w http.ResponseWriter, r *http.Request) {
	// Extrahiere die Query-Parameter aus der URL
	family := r.URL.Query().Get("family")
	cluster := r.URL.Query().Get("cluster")
	filter_tagname := r.URL.Query().Get("tagname")
	filter_namespace := r.URL.Query().Get("namespace")
	kind := "is_tag_used"

	// Logge eine Informationsnachricht über die neue Anfrage
	InfoMsg("--------------  New request  --------------")

	// Validierung der erforderlichen Parameter
	if err := validateParams(family, kind, filter_tagname); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		ErrorMsg("Error:", err)
		return
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

// IsZero prüft, ob T_flags nicht initialisiert ist.
func (t T_flags) IsZero() bool {
	return t.Family == "" && len(t.Cluster) == 0 && t.Filter == (T_flagFilt{}) && t.Output == (T_flagOut{})
}

// initializeCmdParams initializes the command parameters based on the query parameters.
func initializeCmdParams(family, kind, cluster, tagname, namespace string) (T_flags, bool, error) {
	// cmdParams := T_flags{}
	if initialCmdParams.IsZero() {
		initialCmdParams = CmdParams
	} else {
		CmdParams = initialCmdParams
	}
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
