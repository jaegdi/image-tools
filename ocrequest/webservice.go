package ocrequest

import (
	"encoding/json"
	"net/http"
)

// StartServer starts the HTTP server and handles the "/execute" endpoint.
// It expects the "family" and "tagname" parameters to be provided in the URL query.
// If any of the parameters is missing, it returns a HTTP 400 Bad Request error.
// Otherwise, it sets the CmdParams.Family and CmdParams.Filter.Tagname based on the provided values.
// It also sets CmdParams.Output.Used to true.
// The server listens on port 8080 and returns the result of CmdlineMode as a JSON response.
// The response content type is set to "application/json".
func StartServer() {
	http.HandleFunc("/execute", func(w http.ResponseWriter, r *http.Request) {
		family := r.URL.Query().Get("family")
		kind := r.URL.Query().Get("kind")
		cluster := r.URL.Query().Get("cluster")
		filter_tagname := r.URL.Query().Get("tagname")
		filter_namespace := r.URL.Query().Get("namespace")

		InfoMsg("--------------  New request  --------------")

		if kind == "is_tag_used " && (family == "" || filter_tagname == "") {
			http.Error(w, "Missing parameters", http.StatusBadRequest)
			ErrorMsg("Error: Missing parameters")
			ErrorMsg("family:", family, "| kind:", kind, "| tagname:", filter_tagname)
			return
		}

		if kind == "" {
			kind = "is_tag_used"
		}
		html := true

		cmdParams := T_flags{}
		cmdParams.Family = T_familyName(family)
		cmdParams.Cluster = T_clName(cluster).list()
		// Switch block to handle different kinds
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
			// Handle unknown kind
			http.Error(w, "Invalid kind parameter", http.StatusBadRequest)
			ErrorMsg("Error: Invalid kind parameter")
			ErrorMsg("family:", cmdParams.Family, "| kind:", kind, "| tagname:", filter_tagname)
			return
		}
		cmdParams.Filter.Tagname = T_tagName(filter_tagname)
		cmdParams.Filter.Namespace = T_nsName(filter_namespace)
		InitServerMode(cmdParams)

		InfoMsg("family:", family, "| kind:", kind, "| tagname:", filter_tagname)
		VerifyMsg("\nCmdParams Options:", GetJsonFromMap(CmdParams.Options), "Output:", GetJsonFromMap(CmdParams.Output))
		result := CmdlineMode()
		VerifyMsg("\nCmdParams Result:", GetJsonFromMap(result))
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
					// "Result":    result,
					"TagName":   filter_tagname,
					"TagIsUsed": tagIsUsed,
				}
				json.NewEncoder(w).Encode(response)
			} else {
				json.NewEncoder(w).Encode(result)
			}
		}
	})
	http.ListenAndServe(":8080", nil)
}
