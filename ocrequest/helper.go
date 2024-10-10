package ocrequest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"

	"gopkg.in/yaml.v2"

	// "reflect"

	"strconv"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"golang.org/x/crypto/ssh/terminal"
)

var today time.Time = time.Now()
var pager io.WriteCloser

// ageInDays calculates the no of days between a date and current date
func ageInDays(date string) int {
	t, _ := time.Parse(time.RFC3339, date)
	return int(today.Sub(t).Hours()) / 24
}

// ExitWithError write msg to StdErr, and logfile and exits to program
func ExitWithError(errormsg ...interface{}) {
	ErrorMsg(errormsg...)
	os.Stderr.WriteString(fmt.Sprint(errormsg...) + "\n")
	os.Exit(1)
}

// UnescapeUtf8InJsonBytes removes escape sign from JSON byte
func UnescapeUtf8InJsonBytes(_jsonRaw json.RawMessage) (json.RawMessage, error) {
	var patterns = [][2]string{
		{`\\u`, `\u`},
		{`%3a`, `@`},
		{`%3A`, `:`},
		{`\<`, `<`},
		{`\>`, `>`},
		{`\&`, `&`},
	}
	var patterns2 = [][2]string{
		{`%3a`, `@`},
		{`%3A`, `:`},
		{`\<`, `<`},
		{`\>`, `>`},
		{`\&`, `&`},
		{`"\"`, `"`},
		{`\""`, `"`},
		{`\"`, ``},
	}
	str := strconv.Quote(string(_jsonRaw))
	for _, p := range patterns {
		str = strings.Replace(str, p[0], p[1], -1)
	}
	if str1, err := strconv.Unquote(str); err != nil {
		return _jsonRaw, err
	} else {
		for _, p := range patterns2 {
			str1 = strings.Replace(str1, p[0], p[1], -1)
		}
		return []byte(str1), nil
	}
}

func GetJsonOneliner(dict interface{}) string {
	j, err := json.Marshal(dict)
	if err != nil {
		ErrorMsg("dict: ", dict)
		ErrorMsg("err: ", err)
	}
	return string(j)
}

// GetJsonFromMap generate formated json output depending on the commadline flags
func GetJsonFromMap(dict interface{}) string {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "   ")
	if err := encoder.Encode(dict); err != nil {
		if jsonBytes, err := json.MarshalIndent(dict, "", "  "); err != nil {
			ErrorMsg("dict: ", dict)
			ErrorMsg("err: ", err)
		} else {
			return string(jsonBytes)
		}
	} else {
		b := buffer.Bytes()
		if b1, err := UnescapeUtf8InJsonBytes(b); err != nil {
			ErrorMsg("UnescapeUtf8InJsonBytes failed::", "in: ", string(b), "out: ", string(b1))
			return string(b1)
		} else {
			return string(b1)
		}
	}
	return ""
}

// Read a multidoc yaml file and generate data in out as []interface{}
func UnmarshalMultidocYaml(in []byte, out *([]interface{})) error {
	r := bytes.NewReader(in)
	decoder := yaml.NewDecoder(r)
	for {
		var data map[string]interface{}

		if err := decoder.Decode(&data); err != nil {
			// Break when there are no more documents to decode
			if err != io.EOF {
				return err
			}
			break
		}
		*out = append(*out, data)
	}
	return nil
}

// func getFieldFromStruct(v reflect.Type, field string) string {
// 	r := reflect.ValueOf(v)
// 	f := reflect.Indirect(r).FieldByName(field)
// 	return string(f.String())
// }

// GetYamlFromMap generate YAML output from map
func GetYamlFromMap(list interface{}) string {
	d, err := yaml.Marshal(&list)
	if err != nil {
		ErrorMsg("Convert map to Yaml failed", err)
	}
	return string(d)
}

// GetCsvFromMap generates CSV output from a map based on the specified output parameters.
// It processes different types of data (ImageStreams, ImageStreamTags, Images, Used ImageStreamTags, and Unused ImageStreamTags)
// and generates corresponding CSV files.
//
// Parameters:
// - list: The data structure containing the ImageStreamTags information.
// - family: The family name (T_familyName) to include in the CSV output.
func GetCsvFromMap(list interface{}, family T_familyName) {
	if CmdParams.Output.Is || CmdParams.Output.All {
		generateCsvForImageStreams(list, family)
	}
	if CmdParams.Output.Istag || CmdParams.Output.All {
		generateCsvForImageStreamTags(list, family)
	}
	if CmdParams.Output.Image || CmdParams.Output.All {
		generateCsvForImages(list, family)
	}
	if CmdParams.Output.Used || CmdParams.Output.All {
		generateCsvForUsedImageStreamTags(list, family)
	}
	if CmdParams.Output.UnUsed || CmdParams.Output.All {
		generateCsvForUnusedImageStreamTags(list, family)
	}
}

// generateCsvForImageStreams generates CSV output for ImageStreams.
// It iterates through the provided list of image streams and their associated tags, creating a CSV document with the relevant details.
//
// Parameters:
// - list: The data structure containing image streams and their tags.
// - family: The family name used for generating the CSV content.
func generateCsvForImageStreams(list interface{}, family T_familyName) {
	// Initialize the CSV document with a headline
	output := T_csvDoc{}
	headline := T_csvLine{"Family", "DataRange", "DataType", "Imagestream", "Image", "ImagestreamTag", "Cluster"}
	output = append(output, headline)

	// Iterate through each cluster specified in the command parameters
	for _, cluster := range CmdParams.Cluster {
		// Iterate through each image stream in the list for the current cluster
		for is, isMap := range list.(T_completeResults).AllIstags[cluster].Is {
			// Iterate through each image in the image stream
			for image, shaMap := range isMap {
				// Iterate through each image stream tag
				for istag := range shaMap {
					if strings.HasPrefix(istag.str(), family.str()) {
						// Create a new CSV line with the relevant details
						line := T_csvLine{}
						line = append(line, string(family))
						line = append(line, "allIstags")
						line = append(line, "is")
						line = append(line, is.str())
						line = append(line, image.str())
						line = append(line, istag.str())
						line = append(line, cluster.str())
						// Append the line to the CSV document
						output = append(output, line)
					}
				}
			}
		}
	}
	// Generate the CSV file with the name "imagestreams"
	output.csvDoc("imagestreams")
}

// generateCsvForImageStreamTags generates CSV output for ImageStreamTags.
// It iterates through the provided list of image stream tags, creating a CSV document with the relevant details.
//
// Parameters:
// - list: The data structure containing image stream tags.
// - family: The family name used for generating the CSV content.
func generateCsvForImageStreamTags(list interface{}, family T_familyName) {
	// Initialize the CSV document with a headline
	output := T_csvDoc{}
	headline := T_csvLine{"Family", "DataRange", "DataType", "istag", "Cluster"}
	// Append the names of the istag fields to the headline
	headline = append(headline, toArrayString(T_istag{}.Names())...)
	output = append(output, headline)

	// Iterate through each cluster specified in the command parameters
	for _, cluster := range CmdParams.Cluster {
		// Iterate through each image stream tag name in the list for the current cluster
		for istagName, nsMap := range list.(T_completeResults).AllIstags[cluster].Istag {
			// Iterate through each namespace map for the current image stream tag name
			for _, istagMap := range nsMap {
				// Create a new CSV line with the relevant details
				line := T_csvLine{}
				line = append(line, string(family))
				line = append(line, "allIstags")
				line = append(line, "istag")
				line = append(line, istagName.str())
				line = append(line, cluster.str())
				// Append the values of the istag map to the line
				line = append(line, toArrayString(istagMap.Values())...)
				// Append the line to the CSV document
				output = append(output, line)
			}
		}
	}
	// Generate the CSV file with the name "istags"
	output.csvDoc("istags")
}

// generateCsvForImages generates CSV output for Images.
// It iterates through the provided list of images and their associated tags, creating a CSV document with the relevant details.
//
// Parameters:
// - list: The data structure containing images and their tags.
// - family: The family name used for generating the CSV content.
func generateCsvForImages(list interface{}, family T_familyName) {
	// Initialize the CSV document with a headline
	output := T_csvDoc{}
	headline := T_csvLine{"Family", "DataRange", "DataType", "Image", "Istag", "Imagestream", "Cluster", "Namespace", "Date", "AgeInDays", "IsTagReferences"}
	output = append(output, headline)

	// Iterate through each cluster specified in the command parameters
	for _, cluster := range CmdParams.Cluster {
		// Iterate through each image in the list for the current cluster
		for shaName, shaMap := range list.(T_completeResults).AllIstags[cluster].Image {
			// Iterate through each image stream tag associated with the image
			for istag, istagMap := range shaMap {
				// Create a new CSV line with the relevant details
				line := T_csvLine{}
				line = append(line, string(family))
				line = append(line, "allIstags")
				line = append(line, "image")
				line = append(line, shaName.str())
				line = append(line, istag.str())
				line = append(line, istagMap.Imagestream.str())
				line = append(line, cluster.str())
				line = append(line, istagMap.Namespace.str())
				line = append(line, istagMap.Date)
				line = append(line, strconv.Itoa(istagMap.AgeInDays))
				// Iterate through each tag reference in the image stream tag map
				for tag := range istagMap.Istags {
					// Create a copy of the line and append the tag reference
					copyOfLine := append([]string{}, line...)
					copyOfLine = append(copyOfLine, tag.str())
					// Append the line to the CSV document
					output = append(output, copyOfLine)
				}
			}
		}
	}
	// Generate the CSV file with the name "images"
	output.csvDoc("images")
}

// generateCsvForUsedImageStreamTags generates CSV output for Used ImageStreamTags.
// It iterates through the provided list of used image stream tags, creating a CSV document with the relevant details.
//
// Parameters:
// - list: The data structure containing used image stream tags.
// - family: The family name used for generating the CSV content.
func generateCsvForUsedImageStreamTags(list interface{}, family T_familyName) {
	// Initialize the CSV document with a headline
	output := T_csvDoc{}
	headline := T_csvLine{"Family", "DataRange", "DataType", "Imagestream", "Tag"}
	// Append the names of the used image stream tag fields to the headline
	headline = append(headline, toArrayString(T_usedIstag{}.Names())...)
	output = append(output, headline)

	// Iterate through each image stream in the list of used image stream tags
	for is, isMap := range list.(T_completeResults).UsedIstags {
		// Iterate through each image stream tag in the image stream
		for istag, istagArray := range isMap {
			// Iterate through each map of used image stream tag details
			for _, istagMap := range istagArray {
				if nsBelongsToFamily(istagMap.UsedInNamespace, family) || nsBelongsToFamily(istagMap.FromNamespace, family) {
					// Create a new CSV line with the relevant details
					line := T_csvLine{}
					line = append(line, string(family))
					line = append(line, "usedistags")
					line = append(line, "is:tag")
					line = append(line, is.str())
					line = append(line, istag.str())
					// Append the values of the used image stream tag map to the line
					line = append(line, toArrayString(istagMap.Values())...)
					// Append the line to the CSV document
					output = append(output, line)
				}
			}
		}
	}
	// Generate the CSV file with the name "used-istags"
	output.csvDoc("used-istags")
}

// generateCsvForUnusedImageStreamTags generates CSV output for Unused ImageStreamTags.
// It iterates through the provided list of unused image stream tags, creating a CSV document with the relevant details.
//
// Parameters:
// - list: The data structure containing unused image stream tags.
// - family: The family name used for generating the CSV content.
func generateCsvForUnusedImageStreamTags(list interface{}, family T_familyName) {
	// Initialize the CSV document with a headline
	InfoMsg(family)
	output := T_csvDoc{}
	headline := T_csvLine{"unusedImagestreamtag"}
	// Append the names of the unused image stream tag fields to the headline
	headline = append(headline, toArrayString(T_unUsedIstag{}.Names())...)
	output = append(output, headline)

	// Iterate through each unused image stream tag in the list
	for istag, istagMap := range list.(T_completeResults).UnUsedIstags {
		// Create a new CSV line with the relevant details
		line := T_csvLine{}
		line = append(line, istag.str())
		// Append the values of the unused image stream tag map to the line
		line = append(line, toArrayString(istagMap.Values())...)
		// Append the line to the CSV document
		output = append(output, line)
	}
	// Generate the CSV file with the name "unused-istags"
	output.csvDoc("unused-istags")
}

// GetTextTableFromMap generates ASCII table output from a map based on the specified output parameters.
// It processes different types of data (ImageStreams, ImageStreamTags, Images, Used ImageStreamTags, and Unused ImageStreamTags)
// and generates corresponding ASCII tables.
//
// Parameters:
// - list: The data structure containing the ImageStreamTags information.
// - family: The family name (T_familyName) to include in the table output.
func GetTextTableFromMap(list interface{}, family T_familyName) string {
	result := ""
	if CmdParams.Output.Is || CmdParams.Output.All {
		result = generateTableForImageStreams(list, family)
	}
	if CmdParams.Output.Istag || CmdParams.Output.All {
		result = generateTableForImageStreamTags(list, family)
	}
	if CmdParams.Output.Image || CmdParams.Output.All {
		result = generateTableForImages(list, family)
	}
	if CmdParams.Output.Used || CmdParams.Output.All {
		result = generateTableForUsedImageStreamTags(list, family)
	}
	if CmdParams.Output.UnUsed || CmdParams.Output.All {
		result = generateTableForUnusedImageStreamTags(list, family)
	}
	return result
}

// generateTableForImageStreams generates ASCII table output for ImageStreams.
// It iterates through the provided list of image streams and their associated tags, creating an ASCII table with the relevant details.
//
// Parameters:
// - list: The data structure containing image streams and their tags.
// - family: The family name used for generating the table content.
func generateTableForImageStreams(list interface{}, family T_familyName) string {
	// Initialize the table output with a headline
	output := []table.Row{}
	headline := table.Row{"Imagestream " + string(family), "Image", "ImagestreamTag", "Cluster"}
	output = append(output, headline)

	// Iterate through each cluster specified in the command parameters
	for _, cluster := range CmdParams.Cluster {
		// Iterate through each image stream in the list for the current cluster
		for is, isMap := range list.(T_completeResults).AllIstags[cluster].Is {
			// Iterate through each image in the image stream
			for image, shaMap := range isMap {
				// Iterate through each image stream tag
				for istag := range shaMap {
					if strings.HasPrefix(istag.str(), family.str()) {
						// Create a new table row with the relevant details
						line := table.Row{is, image, istag, cluster}
						// Append the row to the table output
						output = append(output, line)
					}
				}
			}
		}
	}
	// Pretty print the table output
	return tablePrettyprint(output)
}

// Beispiel für slice.Contains Funktion
func sliceContains(slice []T_nsName, item T_nsName) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func testFamilyNamespaces(ns T_nsName, family T_familyName) bool {
	if sliceContains(FamilyNamespaces[family].nsList(), ns) {
		// Der Namespace ist in der nsList enthalten
		return true
	}
	return false
}

func nsBelongsToFamily(ns T_nsName, family T_familyName) bool {
	if family.str() == "" || strings.HasPrefix(ns.str(), string(family)) || testFamilyNamespaces(ns, family) {
		return true
	}
	return false
}

// generateTableForImageStreamTags generates ASCII table output for ImageStreamTags.
// It iterates through the provided list of image stream tags, creating an ASCII table with the relevant details.
//
// Parameters:
// - list: The data structure containing image stream tags.
// - family: The family name used for generating the table content.
func generateTableForImageStreamTags(list interface{}, family T_familyName) string {
	// Initialize the table output with a headline
	output := []table.Row{}
	headline := table.Row{"istag " + string(family), "Cluster"}
	// Append the names of the istag fields to the headline
	headline = append(headline, toTableRow(T_istag{}.Names())...)
	output = append(output, headline)

	// Iterate through each cluster specified in the command parameters
	for _, cluster := range CmdParams.Cluster {
		// Iterate through each image stream tag name in the list for the current cluster
		for istagName, nsMap := range list.(T_completeResults).AllIstags[cluster].Istag {
			// Iterate through each namespace map for the current image stream tag name
			for ns, istagMap := range nsMap {
				if nsBelongsToFamily(ns, family) {
					// Create a new table row with the relevant details
					line := table.Row{istagName, cluster}
					// Append the values of the istag map to the row
					line = append(line, toTableRow(istagMap.Values())...)
					// Append the row to the table output
					output = append(output, line)
				}
			}
		}
	}
	// Pretty print the table output
	return tablePrettyprint(output)
}

// generateTableForImages generates ASCII table output for Images.
// It iterates through the provided list of images and their associated tags, creating an ASCII table with the relevant details.
//
// Parameters:
// - list: The data structure containing images and their tags.
// - family: The family name used for generating the table content.
func generateTableForImages(list interface{}, family T_familyName) string {
	// Initialize the table output with a headline
	output := []table.Row{}
	headline := table.Row{"Image " + string(family), "Istag", "Imagestream", "Cluster", "Namespace", "Date", "AgeInDays", "IsTagReferences"}
	output = append(output, headline)

	// Iterate through each cluster specified in the command parameters
	for _, cluster := range CmdParams.Cluster {
		// Iterate through each image in the list for the current cluster
		for shaName, shaMap := range list.(T_completeResults).AllIstags[cluster].Image {
			// Iterate through each image stream tag associated with the image
			for istag, istagMap := range shaMap {
				// Create a new table row with the relevant details
				line := table.Row{shaName, istag, istagMap.Imagestream, cluster, istagMap.Namespace, istagMap.Date, strconv.Itoa(istagMap.AgeInDays)}
				// Iterate through each tag reference in the image stream tag map
				for tag := range istagMap.Istags {
					// Create a copy of the line and append the tag reference
					copyOfLine := append([]interface{}{}, line...)
					copyOfLine = append(copyOfLine, tag)
					// Append the row to the table output
					output = append(output, copyOfLine)
				}
			}
		}
	}
	// Pretty print the table output
	return tablePrettyprint(output)
}

// generateTableForUsedImageStreamTags generates ASCII table output for Used ImageStreamTags.
// It iterates through the provided list of used image stream tags, creating an ASCII table with the relevant details.
//
// Parameters:
// - list: The data structure containing used image stream tags.
// - family: The family name used for generating the table content.
func generateTableForUsedImageStreamTags(list interface{}, family T_familyName) string {
	// Initialize the table output with a headline
	output := []table.Row{}
	headline := table.Row{"Imagestream (used for " + string(family) + ")", "Tag (used)"}
	// Append the names of the used image stream tag fields to the headline
	headline = append(headline, toTableRow(T_usedIstag{}.Names())...)
	output = append(output, headline)

	// Iterate through each image stream in the list of used image stream tags
	for is, isMap := range list.(T_completeResults).UsedIstags {
		// Iterate through each image stream tag in the image stream
		for istag, istagArray := range isMap {
			// Iterate through each map of used image stream tag details
			for _, istagMap := range istagArray {
				// Create a new table row with the relevant details
				line := table.Row{is, istag}
				// Append the values of the used image stream tag map to the row
				line = append(line, toTableRow(istagMap.Values())...)
				// Append the row to the table output
				output = append(output, line)
			}
		}
	}
	// Pretty print the table output
	return tablePrettyprint(output)
}

// generateTableForUnusedImageStreamTags generates ASCII table output for Unused ImageStreamTags.
// It iterates through the provided list of unused image stream tags, creating an ASCII table with the relevant details.
//
// Parameters:
// - list: The data structure containing unused image stream tags.
// - family: The family name used for generating the table content.
func generateTableForUnusedImageStreamTags(list interface{}, family T_familyName) string {
	InfoMsg(family)
	// Initialize the table output with a headline
	output := []table.Row{}
	headline := table.Row{"unused Imagestreamtag"}
	// Append the names of the unused image stream tag fields to the headline
	headline = append(headline, toTableRow(T_unUsedIstag{}.Names())...)
	output = append(output, headline)

	// Iterate through each unused image stream tag in the list
	for istag, istagMap := range list.(T_completeResults).UnUsedIstags {
		// Create a new table row with the relevant details
		line := table.Row{istag}
		// Append the values of the unused image stream tag map to the row
		line = append(line, toTableRow(istagMap.Values())...)
		// Append the row to the table output
		output = append(output, line)
	}
	// Pretty print the table output
	return tablePrettyprint(output)
}

// toTableRow converts a variadic slice of interface{} to a table.Row.
// It iterates through the provided slice of interface{} and flattens any nested slices.
//
// Parameters:
// - arr: A variadic slice of interface{} representing the table row values.
//
// Returns:
// - A table.Row containing the flattened values.
func toTableRow(arr ...interface{}) table.Row {
	o := table.Row{}
	// Iterate through each element in the variadic slice
	for _, v := range arr {
		// Type assert the element to a slice of interface{}
		for _, w := range v.([]interface{}) {
			// Append each element of the nested slice to the table.Row
			o = append(o, w)
		}
	}
	return o
}

// toArrayString converts a variadic slice of interface{} to a slice of strings.
// It iterates through the provided slice of interface{} and flattens any nested slices, converting each element to a string.
//
// Parameters:
// - arr: A variadic slice of interface{} representing the input values.
//
// Returns:
// - A slice of strings containing the flattened and converted values.
func toArrayString(arr ...interface{}) []string {
	o := []string{}
	// Iterate through each element in the variadic slice
	for _, v := range arr {
		// Type assert the element to a slice of interface{}
		for _, w := range v.([]interface{}) {
			// Append each element of the nested slice to the output slice after converting to string
			o = append(o, w.(string))
		}
	}
	return o
}

// tablePrettyprint prints an ASCII table.
// It formats and displays the provided table rows in a paginated and styled ASCII table.
//
// Parameters:
// - out: A slice of table.Row representing the rows to be printed.
func tablePrettyprint(out []table.Row) string {
	// If there are no rows to print, return immediately
	if len(out) == 0 {
		return ""
	}

	// Get the terminal height
	height := getTerminalHeight()

	// Activate the pager
	cmd, pager, pagername := activatePager()
	defer func() {
		// Ensure the pager is closed and the command waits for completion
		if pager != nil {
			pager.Close()
			_ = cmd.Wait()
		}
	}()

	// Define and render the table output
	return renderTable(out, pager, height, pagername)
}

// getTerminalHeight retrieves the terminal height.
// If it fails to get the terminal size, it sets a default height.
//
// Returns:
// - An integer representing the terminal height.
func getTerminalHeight() int {
	if !CmdParams.Html {
		fd := int(os.Stdout.Fd())
		_, height, terr := terminal.GetSize(fd)
		if terr != nil {
			// If getting the terminal size fails, set a default height
			height = 60
			ErrorMsg("failed to get terminal size, set it to", height)
		}
		return height
	}
	return 0
}

// activatePager activates the pager for table output.
//
// Returns:
// - A pointer to exec.Cmd representing the pager command.
// - A pipe writer for the pager.
func activatePager() (*exec.Cmd, io.WriteCloser, string) {
	if CmdParams.Html {
		return nil, nil, ""
	}
	cmd, pager, pagername := runPager()
	return cmd, pager, pagername
}

// renderTable defines and renders the table output.
//
// Parameters:
// - out: A slice of table.Row representing the rows to be printed.
// - pager: A pipe writer for the pager.
// - height: An integer representing the terminal height.
func renderTable(out []table.Row, pager io.WriteCloser, height int, pagername string) string {
	// Define the table output
	t := table.NewWriter()
	if !CmdParams.Html {
		if pager != nil {
			t.SetOutputMirror(pager) // Set the output to the pager
			if pagername == "ov" {   // If the pager is ov, set the page size to zero
				t.SetPageSize(0)
			} else {
				t.SetPageSize(height - 4) // Set the page size based on terminal height
			}
		}
		// t.SetStyle(table.StyleColoredBright) // Set the table style
		t.SetStyle(table.StyleColoredBlackOnGreenWhite)
	}
	t.AppendHeader(out[0])                   // Append the header row
	t.AppendFooter(table.Row{" ", " ", " "}) // Append a footer row
	t.AppendRows(out[1:])                    // Append the data rows
	t.SortBy([]table.SortBy{                 // Define sorting for the table
		{Number: 1, Mode: table.Asc},
		{Number: 2, Mode: table.Asc},
		{Number: 3, Mode: table.Asc},
	})
	t.SetAutoIndex(true) // Enable automatic indexing of rows

	// If TabGroup parameter is set, configure column merging
	if CmdParams.TabGroup {
		t.SetColumnConfigs([]table.ColumnConfig{
			{Number: 1, AutoMerge: true},
			{Number: 2, AutoMerge: true},
			{Number: 3, AutoMerge: true},
			{Number: 4, AutoMerge: true},
			{Number: 5, AutoMerge: true},
			{Number: 6, AutoMerge: true},
			{Number: 7, AutoMerge: true},
			{Number: 8, AutoMerge: true},
		})
	}

	// If Html parameter is set, render the table as HTML
	if CmdParams.Html {
		t.Style().HTML = table.HTMLOptions{
			CSSClass:    "game-of-thrones",
			EmptyColumn: "&nbsp;",
			EscapeText:  true,
			Newline:     "<br/>",
		}
		htmlTable := t.RenderHTML()
		downloadButton := `
            <button onclick="downloadTableAsExcel()">Download as Excel</button>
        `
		downloadButton2 := `
            <script>
                function downloadTableAsExcel() {
                    var table = document.querySelector('.dataTables_scrollBody > table');
                    var html = table.outerHTML;
                    var url = 'data:application/vnd.ms-excel,' + escape(html);
                    var link = document.createElement('a');
                    link.href = url;
                    link.setAttribute('download', 'table.xls');
                    link.click();
                }
            </script>
        `
		dataTableScript := `
            <link rel="stylesheet" type="text/css" href="https://cdn.datatables.net/1.11.5/css/jquery.dataTables.css">
            <script type="text/javascript" charset="utf8" src="https://code.jquery.com/jquery-3.5.1.js"></script>
            <script type="text/javascript" charset="utf8" src="https://cdn.datatables.net/1.11.5/js/jquery.dataTables.js"></script>
            <style>
                html, body {
                    height: 100%;
                    margin: 0;
                    padding: 0;
                }
                .dataTables_wrapper {
                    height: 100%;
                    display: flex;
                    flex-direction: column;
                }
                .dataTables_wrapper .dataTables_filter {
                    display: flex;
                    justify-content: space-between;
                    align-items: center;
                }
                table {
                    width: 100%;
                }
                .dataTables_scrollBody {
                    flex: 1 1 auto;
                    overflow: auto;
                }
                .dataTables_scrollHead {
                    flex: 0 0 auto;
                }
            </style>
            <script>
                $(document).ready(function() {
                    var table = $('table').DataTable({
                        "scrollY": "calc(100vh - 150px)",
                        "scrollCollapse": true,
                        "paging": false,
                        "scrollX": true
                    });
                    $('.dataTables_filter').append(` + "`" + downloadButton + "`" + `);
                });
            </script>
        `
		return dataTableScript + htmlTable + downloadButton2
	} else {
		t.Render()
		return ""
	}
}

// runPager starts less or the standard pager of os and pipes output into its Stdin
func runPager() (*exec.Cmd, io.WriteCloser, string) {
	var pager string
	var cmd *exec.Cmd

	if CmdParams.Html {
		return nil, nil, ""
	}
	if _, err := exec.LookPath("ov"); err == nil {
		pager = "ov"
	} else {
		pager = os.Getenv("PAGER")
	}
	switch pager {
	case "less":
		{
			cmd = exec.Command(pager, "-m", "-n", "-g", "-i", "-J", "-R", "-S", "--underline-special", "--SILENT")
		}
	case "ov":
		{
			cmd = exec.Command("ov", "-w=false", "-c", "-H", "1", "-d", "' +'")
		}
	default:
		{
			cmd = exec.Command(pager)
		}
	}

	out, err := cmd.StdinPipe()
	if err != nil {
		ErrorMsg("ExecError", err)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		ErrorMsg("ExecError", err)
	}
	return cmd, out, pager
}
