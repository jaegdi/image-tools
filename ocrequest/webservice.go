package ocrequest

import (
	"encoding/json"
	"fmt"
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
		tagname := r.URL.Query().Get("tagname")

		if family == "" || tagname == "" {
			http.Error(w, "Missing parameters", http.StatusBadRequest)
			fmt.Println("Error: Missing parameters")
			return
		} else {
			fmt.Println("family: ", family, "tagname: ", tagname)
		}

		CmdParams.Family = T_familyName(family)
		CmdParams.Filter.Tagname = T_tagName(tagname)
		// CmdParams.Options.Debug = true
		// CmdParams.Options.Verify = true
		fmt.Println(GetJsonFromMap(CmdParams))
		result := CmdlineMode()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})
	http.ListenAndServe(":8080", nil)
}
