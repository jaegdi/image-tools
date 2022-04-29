package ocrequest

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"os/exec"

	"gopkg.in/yaml.v2"

	// "reflect"
	"runtime"
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

// exitWithError write msg to StdErr, and logfile and exits to program
func exitWithError(errormsg ...interface{}) {
	ErrorLogger.Println(errormsg...)
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
		ErrorLogger.Println("dict: ", dict)
		ErrorLogger.Println("err: ", err)
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
			ErrorLogger.Println("dict: ", dict)
			ErrorLogger.Println("err: ", err)
		} else {
			s := string(jsonBytes)
			return s
		}
	} else {
		b := buffer.Bytes()
		if b1, err := UnescapeUtf8InJsonBytes(b); err != nil {
			ErrorLogger.Println("UnescapeUtf8InJsonBytes failed::", "in: ", string(b), "out: ", string(b1))
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
		ErrorLogger.Println("Convert map to Yaml failed", err)
	}
	return string(d)
}

// GetCsvFromMap generates CSV output from map
func GetCsvFromMap(list interface{}, family T_familyName) {
	if CmdParams.Output.Is || CmdParams.Output.All {
		output := T_csvDoc{}
		headline := T_csvLine{"Family", "DataRange", "DataType", "Imagestream", "Image", "ImagestreamTag", "Cluster"}
		output = append(output, headline)
		for _, cluster := range CmdParams.Cluster {
			for is, isMap := range list.(T_completeResults).AllIstags[cluster].Is {
				for image, shaMap := range isMap {
					for istag := range shaMap {
						line := T_csvLine{}
						line = append(line, string(family))
						line = append(line, "allIstags")
						line = append(line, "is")
						line = append(line, is.str())
						line = append(line, image.str())
						line = append(line, istag.str())
						line = append(line, cluster.str())
						output = append(output, line)
					}
				}
			}
		}
		output.csvDoc("imagestreams")
	}
	if CmdParams.Output.Istag || CmdParams.Output.All {
		output := T_csvDoc{}
		// headline := T_csvLine{"Family", "DataRange", "DataType", "istag"} //, "Imagestream", "Tagname", "Cluster", "Namespace", "Link", "Date", "AgeInDays", "Image", "CommitAuthor", "CommitDate", "CommitId", "CommitRef", "Commitversion", "IsProdImage", "BuildNName", "BuildNamespace"}
		headline := T_csvLine{"Family", "DataRange", "DataType", "istag"} //, "Imagestream", "Tagname", "Cluster", "Namespace", "Date", "AgeInDays", "Image", "CommitAuthor", "CommitDate", "CommitId", "CommitRef", "Commitversion", "IsProdImage", "BuildNName", "BuildNamespace"}
		headline = append(headline, "Cluster")
		headline = append(headline, toArrayString(T_istag{}.Names())...)
		output = append(output, headline)
		for _, cluster := range CmdParams.Cluster {
			for istagName, nsMap := range list.(T_completeResults).AllIstags[cluster].Istag {
				for _, istagMap := range nsMap {
					// InfoLogger.Println("namespace:", ns)
					line := T_csvLine{}
					line = append(line, string(family))
					line = append(line, "allIstags")
					line = append(line, "istag")
					line = append(line, istagName.str())
					line = append(line, cluster.str())
					line = append(line, toArrayString(istagMap.Values())...)
					output = append(output, line)
				}
			}
		}
		output.csvDoc("istags")
	}
	if CmdParams.Output.Image || CmdParams.Output.All {
		output := T_csvDoc{}
		// headline := T_csvLine{"Family", "DataRange", "DataType", "Image", "Istag", "Imagestream", "Cluster", "Namespace", "Link", "Date", "AgeInDays", "IsTagReferences"}
		headline := T_csvLine{"Family", "DataRange", "DataType", "Image", "Istag", "Imagestream", "Cluster", "Namespace", "Date", "AgeInDays", "IsTagReferences"}
		output = append(output, headline)
		for _, cluster := range CmdParams.Cluster {
			for shaName, shaMap := range list.(T_completeResults).AllIstags[cluster].Image {
				for istag, istagMap := range shaMap {
					line := T_csvLine{}
					line = append(line, string(family))
					line = append(line, "allIstags")
					line = append(line, "image")
					line = append(line, shaName.str())
					line = append(line, istag.str())
					line = append(line, istagMap.Imagestream.str())
					line = append(line, cluster.str())
					line = append(line, istagMap.Namespace.str())
					// line = append(line, istagMap.Link)
					line = append(line, istagMap.Date)
					line = append(line, strconv.Itoa(istagMap.AgeInDays))
					for tag := range istagMap.Istags {
						//  make a real copy of line !!!!
						copyOfLine := append([]string{}, line...)
						copyOfLine = append(copyOfLine, tag.str())
						output = append(output, copyOfLine)
					}
				}
			}
		}
		output.csvDoc("images")
	}
	if CmdParams.Output.Used || CmdParams.Output.All {
		output := T_csvDoc{}
		headline := T_csvLine{"Family", "DataRange", "DataType", "Imagestream", "Tag"} //, "UsedInNamespace", "Image", "UsedInCluster"}
		headline = append(headline, toArrayString(T_usedIstag{}.Names())...)
		output = append(output, headline)
		for is, isMap := range list.(T_completeResults).UsedIstags {
			for istag, istagArray := range isMap { //.(map[string][]map[string]string) {
				for _, istagMap := range istagArray {
					line := T_csvLine{}
					line = append(line, string(family))
					line = append(line, "usedistags")
					line = append(line, "is:tag")
					line = append(line, is.str())
					line = append(line, istag.str())
					line = append(line, toArrayString(istagMap.Values())...)
					output = append(output, line)
				}
			}
		}
		output.csvDoc("used-istags")
	}
	if CmdParams.Output.UnUsed || CmdParams.Output.All {
		output := T_csvDoc{}
		headline := T_csvLine{"unusedImagestreamtag"}
		headline = append(headline, toArrayString(T_unUsedIstag{}.Names())...)
		output = append(output, headline)
		for istag, istagMap := range list.(T_completeResults).UnUsedIstags {
			line := T_csvLine{}
			line = append(line, istag.str())
			line = append(line, toArrayString(istagMap.Values())...)
			output = append(output, line)
		}
		output.csvDoc("unused-istags")
	}
}

// GetTableFromMap generate ASCII table output from map
func GetTableFromMap(list interface{}, family T_familyName) {
	if CmdParams.Output.Is || CmdParams.Output.All {
		output := []table.Row{}
		headline := table.Row{"Imagestream " + string(family), "Image", "ImagestreamTag", "Cluster"}
		output = append(output, headline)
		for _, cluster := range CmdParams.Cluster {
			for is, isMap := range list.(T_completeResults).AllIstags[cluster].Is {
				for image, shaMap := range isMap {
					for istag := range shaMap {
						line := table.Row{}
						// line = append(line, "allIstags")
						// line = append(line, "is")
						line = append(line, is)
						line = append(line, image)
						line = append(line, istag)
						line = append(line, cluster)
						output = append(output, line)
					}
				}
			}
		}
		tablePrettyprint(output)
	}
	if CmdParams.Output.Istag || CmdParams.Output.All {
		output := []table.Row{}
		// headline := table.Row{"istag " + string(family)} //, "Cluster", "Imagestream", "Tagname", "Namespace", "Link", "Date", "AgeInDays", "Image", "CommitAuthor", "CommitDate", "CommitId", "CommitRef", "Commitversion", "IsProdImage", "BuildNName", "BuildNamespace"}
		headline := table.Row{"istag " + string(family)} //, "Cluster", "Imagestream", "Tagname", "Namespace", "Date", "AgeInDays", "Image", "CommitAuthor", "CommitDate", "CommitId", "CommitRef", "Commitversion", "IsProdImage", "BuildNName", "BuildNamespace"}
		headline = append(headline, "Cluster")
		headline = append(headline, toTableRow(T_istag{}.Names())...)
		output = append(output, headline)
		for _, cluster := range CmdParams.Cluster {
			for istagName, nsMap := range list.(T_completeResults).AllIstags[cluster].Istag {
				for _, istagMap := range nsMap {
					// InfoLogger.Println("namespace:", ns)
					line := table.Row{}
					// line = append(line, "allIstags")
					// line = append(line, "istag")
					line = append(line, istagName)
					line = append(line, cluster)
					line = append(line, toTableRow(istagMap.Values())...)
					output = append(output, line)
				}
			}
		}
		tablePrettyprint(output)
	}
	if CmdParams.Output.Image || CmdParams.Output.All {
		output := []table.Row{}
		// headline := table.Row{"Image " + string(family), "Istag", "Imagestream", "Cluster", "Namespace", "Link", "Date", "AgeInDays", "IsTagReferences"}
		headline := table.Row{"Image " + string(family), "Istag", "Imagestream", "Cluster", "Namespace", "Date", "AgeInDays", "IsTagReferences"}
		output = append(output, headline)
		for _, cluster := range CmdParams.Cluster {
			for shaName, shaMap := range list.(T_completeResults).AllIstags[cluster].Image {
				for istag, istagMap := range shaMap {
					line := table.Row{}
					// line = append(line, "allIstags")
					// line = append(line, "image")
					line = append(line, shaName)
					line = append(line, istag)
					line = append(line, istagMap.Imagestream)
					line = append(line, cluster)
					line = append(line, istagMap.Namespace)
					// line = append(line, istagMap.Link)
					line = append(line, istagMap.Date)
					line = append(line, strconv.Itoa(istagMap.AgeInDays))
					for tag := range istagMap.Istags {
						//  make a real copy of line !!!!
						copyOfLine := append([]interface{}{}, line...)
						copyOfLine = append(copyOfLine, tag)
						output = append(output, copyOfLine)
					}
				}
			}
		}
		tablePrettyprint(output)
	}
	if CmdParams.Output.Used || CmdParams.Output.All {
		output := []table.Row{}
		headline := table.Row{"Imagestream (used for " + string(family) + ")", "Tag (used)"}
		headline = append(headline, toTableRow(T_usedIstag{}.Names())...)
		output = append(output, headline)
		for is, isMap := range list.(T_completeResults).UsedIstags {
			for istag, istagArray := range isMap { //.(map[string][]map[string]string) {
				for _, istagMap := range istagArray {
					line := table.Row{}
					// line = append(line, "usedistags")
					// line = append(line, "is:tag")
					line = append(line, is)
					line = append(line, istag)
					line = append(line, toTableRow(istagMap.Values())...)
					output = append(output, line)
				}
			}
		}
		tablePrettyprint(output)
	}
	if CmdParams.Output.UnUsed || CmdParams.Output.All {
		output := []table.Row{}
		headline := table.Row{"unused Imagestreamtag"}
		headline = append(headline, toTableRow(T_unUsedIstag{}.Names())...)
		output = append(output, headline)
		for istag, istagMap := range list.(T_completeResults).UnUsedIstags {
			line := table.Row{}
			line = append(line, istag)
			line = append(line, toTableRow(istagMap.Values())...)
			output = append(output, line)
		}
		tablePrettyprint(output)
	}
}

// toTableRow convert a slice of interface{} to table.Row
func toTableRow(arr ...interface{}) table.Row {
	o := table.Row{}
	for _, v := range arr {
		for _, w := range v.([]interface{}) {
			o = append(o, w)
		}
	}
	return o
}

// toArrayString convert a slice of interface{} to slice of string
func toArrayString(arr ...interface{}) []string {
	o := []string{}
	for _, v := range arr {
		for _, w := range v.([]interface{}) {
			o = append(o, w.(string))
		}
	}
	return o
}

// tablePrettyprint print ASCII table
func tablePrettyprint(out []table.Row) {
	if len(out) == 0 {
		return
	}
	// get height of terminal
	// _, height, err := terminal.GetSize(0)
	// if err != nil {
	// 	InfoLogger.Println("failedt o get terminal size")
	// 	height = 60
	// }
	fd := int(os.Stdout.Fd())
	_, height, _ := terminal.GetSize(fd)

	// activate pager
	var cmd *exec.Cmd
	cmd, pager = runPager()
	defer func() {
		pager.Close()
		_ = cmd.Wait()
	}()
	// define table output
	t := table.NewWriter()
	t.SetOutputMirror(pager)
	t.AppendHeader(out[0])
	t.AppendFooter(table.Row{" ", " ", " "})
	t.AppendRows(out[1:])
	t.SetStyle(table.StyleColoredBright)
	t.SortBy([]table.SortBy{
		{Number: 1, Mode: table.Asc},
		{Number: 2, Mode: table.Asc},
		{Number: 3, Mode: table.Asc},
	})
	t.SetAutoIndex(true)
	// t.SetStyle(table.StyleLight)
	// t.Style().Options.SeparateRows = true
	// t.SetAllowedRowLength(450)
	t.SetPageSize(height - 4)
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
	if CmdParams.Html {
		t.Style().HTML = table.HTMLOptions{
			CSSClass:    "game-of-thrones",
			EmptyColumn: "&nbsp;",
			EscapeText:  true,
			Newline:     "<br/>",
		}
		t.RenderHTML()
	} else {
		t.Render()
	}
}

// runPager starts less or the standard pager of os and pipes output into its Stdin
func runPager() (*exec.Cmd, io.WriteCloser) {
	pager := os.Getenv("PAGER")
	if pager == "" {
		pager = "more"
	}
	var cmd *exec.Cmd
	if pager == "less" {
		cmd = exec.Command(pager, "-m", "-n", "-g", "-i", "-J", "-R", "-S", "--underline-special", "--SILENT")
	} else {
		cmd = exec.Command(pager)
	}
	out, err := cmd.StdinPipe()
	if err != nil {
		ErrorLogger.Println("ExecError", err)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		ErrorLogger.Println("ExecError", err)
	}
	return cmd, out
}

// openbrowser alternative to open browser for html
// TODO
func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		ErrorLogger.Println("unsupported platform")
	}
	if err != nil {
		ErrorLogger.Println("ExecError", err)
	}

}
