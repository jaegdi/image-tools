package ocrequest

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/template"
)

const tpl = `
<!DOCTYPE html>
<html>
<head>
    <title>.Title</title>
    <link rel="stylesheet" type="text/css" href="https://cdn.datatables.net/1.11.5/css/jquery.dataTables.css">
    <script type="text/javascript" charset="utf8" src="https://code.jquery.com/jquery-3.5.1.js"></script>
    <script type="text/javascript" charset="utf8" src="https://cdn.datatables.net/1.11.5/js/jquery.dataTables.js"></script>
    <style>
        tr:nth-child(even) {background-color: #f2f2f2;}
    </style>
    <script>
        $(document).ready(function() {
            $('#resultTable').DataTable();
        });
    </script>
</head>
<body>
    <h1>Result for Family: {{.Family}}</h1>
    <table id="resultTable" class="display">
        <thead>
            <tr>
                {{range .Headers}}
                <th>{{.}}</th>
                {{end}}
            </tr>
        </thead>
        <tbody>
            {{range .Data}}
            <tr>
                {{range $.Headers}}
                <td>{{index . $}}</td>
                {{end}}
            </tr>
            {{end}}
        </tbody>
    </table>
</body>
</html>
`

type T_htmlHeaders = []string
type T_htmlData = []map[string]interface{}
type T_htmlTitle = []byte

// GetHtmlTableFromMap generates HTML table output from a map based on the specified output parameters.
// It processes different types of data (ImageStreams, ImageStreamTags, Images, Used ImageStreamTags, and Unused ImageStreamTags)
// and generates corresponding HTML tables.
//
// Parameters:
// - list: The data structure containing the ImageStreamTags information.
// - family: The family name (T_familyName) to include in the table output.
func GetHtmlTableFromMap(list interface{}, family T_familyName) ([]byte, error) {
	// Define the HTML template with DataTables integration and alternating row colors

	var buf []byte
	writer := bytes.NewBuffer(buf)
	err := error(nil)

	// Extract headers from the data
	title := T_htmlTitle{}
	headers := T_htmlHeaders{}
	data := T_htmlData{}

	// Assuming list is a slice of maps
	switch v := list.(type) {
	case []map[string]interface{}:
		data = v
		if len(data) > 0 {
			for key := range data[0] {
				headers = append(headers, key)
			}
		}
	// case T_completeResults:
	// 	data := list.(T_completeResults).UsedIstags
	// 	if len(data) > 0 {
	// 		for key := range data {
	// 			headers = append(headers, key.str())
	// 		}
	// 	}
	default:
		fmt.Println("Type of list:", v)
		fmt.Println("Unsupported data type")
		err = fmt.Errorf("Unsupported data type")
		return buf, err
	}

	// Create a new template and parse the HTML into it
	tmpl, err := template.New("table").Parse(tpl)
	if err != nil {
		fmt.Println("Error creating template:", err)
		return buf, err
	}

	// Create the output file

	if CmdParams.Output.Is || CmdParams.Output.All {
		title, headers, data = generateHtmlTableForImageStreams(list, family)
	}
	if CmdParams.Output.Istag || CmdParams.Output.All {
		title, headers, data = generateHtmlTableForImageStreamTags(list, family)
	}
	if CmdParams.Output.Image || CmdParams.Output.All {
		title, headers, data = generateHtmlTableForImages(list, family)
	}
	if CmdParams.Output.Used || CmdParams.Output.All {
		title, headers, data = generateHtmlTableForUsedImageStreamTags(list, family)
	}
	if CmdParams.Output.UnUsed || CmdParams.Output.All {
		title, headers, data = generateHtmlTableForUnusedImageStreamTags(list, family)
	}

	// Execute the template with the result data
	err = tmpl.Execute(writer, struct {
		Family  T_familyName
		Title   T_htmlTitle
		Headers T_htmlHeaders
		Data    T_htmlData
	}{
		Family:  family,
		Title:   title,
		Headers: headers,
		Data:    data,
	})
	if err != nil {
		fmt.Println("Error executing template:", err)
		return buf, err
	}
	return buf, err
}

// generateHtmlTableForImageStreams generates HTML table output for ImageStreams.
// It iterates through the provided list of image streams and their associated tags, creating an HTML table with the relevant details.
//
// Parameters:
// - list: The data structure containing image streams and their tags.
// - family: The family name used for generating the HTML content.
//
// Returns:
// - A string containing the HTML table.
func generateHtmlTableForImageStreams(list interface{}, family T_familyName) (T_htmlTitle, T_htmlHeaders, T_htmlData) {
	// Initialize the HTML output with a table header
	title := T_htmlTitle("ImageStreams")
	headers := T_htmlHeaders{"ImageStreams", "Family", "DataRange", "DataType", "Imagestream", "Image", "ImagestreamTag", "Cluster"}
	data := T_htmlData{}
	// Iterate through each cluster specified in the command parameters
	for _, cluster := range CmdParams.Cluster {
		// Iterate through each image stream in the list for the current cluster
		for is, isMap := range list.(T_completeResults).AllIstags[cluster].Is {
			// Iterate through each image in the image stream
			for image, shaMap := range isMap {
				// Iterate through each image stream tag
				for istag := range shaMap {
					// Append a new row to the HTML table with the relevant details
					data = append(data, map[string]interface{}{
						"ImageStreams":   "ImageStreams",
						"Family":         family.str(),
						"DataRange":      "allIstags",
						"Imagestream":    is.str(),
						"Image":          image.str(),
						"ImagestreamTag": istag.str(),
						"Cluster":        cluster.str(),
					})
				}
			}
		}
	}
	// Close the HTML table
	return title, headers, data
}

// generateHtmlTableForImageStreamTags generates HTML table output for ImageStreamTags.
// It iterates through the provided list of image stream tags, creating an HTML table with the relevant details.
//
// Parameters:
// - list: The data structure containing image stream tags.
// - family: The family name used for generating the HTML content.
//
// Returns:
// - A string containing the HTML table.
func generateHtmlTableForImageStreamTags(list interface{}, family T_familyName) (T_htmlTitle, T_htmlHeaders, T_htmlData) {
	// Initialize the HTML output with a table header
	title := T_htmlTitle("ImageStreamTags")
	headers := T_htmlHeaders{"Family", "Istag", "Namespace", "Imagestream", "Date", "Age", "Cluster"}
	data := T_htmlData{}
	// Iterate through each cluster specified in the command parameters
	for _, cluster := range CmdParams.Cluster {
		// Iterate through each image stream tag name in the list for the current cluster
		for istagName, nsMap := range list.(T_completeResults).AllIstags[cluster].Istag {
			// Iterate through each namespace map for the current image stream tag name
			for _, istagMap := range nsMap {
				// Append a new row to the HTML table with the relevant details
				data = append(data, map[string]interface{}{
					"Family":      family.str(),
					"Istag":       istagName.str(),
					"Namespace":   istagMap.Namespace.str(),
					"Imagestream": istagMap.Imagestream.str(),
					"Date":        istagMap.Date,
					"Age":         strconv.Itoa(istagMap.AgeInDays),
					"Cluster":     cluster.str(),
				})
			}
		}
	}
	// Close the HTML table
	return title, headers, data
}

// generateHtmlTableForImages generates HTML table output for Images.
// It iterates through the provided list of images and their associated tags, creating an HTML table with the relevant details.
//
// Parameters:
// - list: The data structure containing images and their tags.
// - family: The family name used for generating the HTML content.
//
// Returns:
// - A string containing the HTML table.
func generateHtmlTableForImages(list interface{}, family T_familyName) (T_htmlTitle, T_htmlHeaders, T_htmlData) {
	// Initialize the HTML output with a table header
	title := T_htmlTitle("Images")
	headers := T_htmlHeaders{"Family", "Image", "ShaName", "Istag", "Imagestream", "Cluster", "Namespace", "Date", "Age", "Tags"}
	data := T_htmlData{}
	// Iterate through each cluster specified in the command parameters
	for _, cluster := range CmdParams.Cluster {
		// Iterate through each image in the list for the current cluster
		for shaName, shaMap := range list.(T_completeResults).AllIstags[cluster].Image {
			// Iterate through each image stream tag associated with the image
			for istag, istagMap := range shaMap {
				// Append a new row to the HTML table with the relevant details
				data = append(data, map[string]interface{}{
					"Family":      family.str(),
					"Image":       shaName.str(),
					"ShaName":     shaName.str(),
					"Istag":       istag.str(),
					"Imagestream": istagMap.Imagestream.str(),
					"Cluster":     cluster.str(),
					"Namespace":   istagMap.Namespace.str(),
					"Date":        istagMap.Date,
					"Age":         strconv.Itoa(istagMap.AgeInDays),
					"Tags":        strings.Join(istagMap.Istags.List(), ","),
				})
			}
		}
	}
	// Close the HTML table
	return title, headers, data
}

// generateHtmlTableForUsedImageStreamTags generates HTML table output for Used ImageStreamTags.
// It iterates through the provided list of used image stream tags, creating an HTML table with the relevant details.
//
// Parameters:
// - list: The data structure containing used image stream tags.
// - family: The family name used for generating the HTML content.
//
// Returns:
// - A string containing the HTML table.
func generateHtmlTableForUsedImageStreamTags(list interface{}, family T_familyName) (T_htmlTitle, T_htmlHeaders, T_htmlData) {
	// Initialize the HTML output with a table header
	title := T_htmlTitle("Used ImageStreamTags")
	headers := T_htmlHeaders{"Family", "Imagestream", "Is", "Tag", "Image", "Namespace", "UsedInNamespace", "Age", "Cluster"}
	data := T_htmlData{}
	// Iterate through each image stream in the list of used image stream tags
	for is, isMap := range list.(T_completeResults).UsedIstags {
		// Iterate through each image stream tag in the image stream
		for istag, istagArray := range isMap {
			// Iterate through each map of used image stream tag details
			for _, istagMap := range istagArray {
				data = append(data, map[string]interface{}{
					"Family":          family.str(),
					"Imagestream":     is.str(),
					"Is":              "is:tag",
					"Tag":             istag.str(),
					"Image":           istagMap.Image.str(),
					"Namespace":       istagMap.FromNamespace.str(),
					"UsedInNamespace": istagMap.UsedInNamespace.str(),
					"Age":             strconv.Itoa(istagMap.AgeInDays),
					"Cluster":         istagMap.Cluster.str(),
				})
			}
		}
	}
	// Close the HTML table
	return title, headers, data
}

// generateHtmlTableForUnusedImageStreamTags generates HTML table output for Unused ImageStreamTags.
// It iterates through the provided list of unused image stream tags, creating an HTML table with the relevant details.
//
// Parameters:
// - list: The data structure containing unused image stream tags.
// - family: The family name used for generating the HTML content.
//
// Returns:
// - A string containing the HTML table.
func generateHtmlTableForUnusedImageStreamTags(list interface{}, family T_familyName) (T_htmlTitle, T_htmlHeaders, T_htmlData) {
	// Initialize the HTML output with a table header
	title := T_htmlTitle("Unused ImageStreamTags")
	headers := T_htmlHeaders{"Family", "ImagestreamTag", "Cluster"}
	data := T_htmlData{}
	// Iterate through each unused image stream tag in the list
	for istag, istagMap := range list.(T_completeResults).UnUsedIstags {
		// Append the cells for the unused image stream tag map values
		data = append(data, map[string]interface{}{
			"Family":         family.str(),
			"ImagestreamTag": istag.str(),
			"Cluster":        istagMap.Cluster.str(),
		})
	}
	// Close the HTML table
	return title, headers, data
}
