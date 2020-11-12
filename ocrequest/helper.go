package ocrequest

import (
	// "github.com/imdario/mergo"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	// "reflect"
	"github.com/jedib0t/go-pretty/v6/table"
	"golang.org/x/crypto/ssh/terminal"
	"strconv"
	"strings"
	"time"
)

var today time.Time = time.Now()

func ageInDays(date string) int {
	t, _ := time.Parse(time.RFC3339, date)
	return int(today.Sub(t).Hours()) / 24
}

func exitWithError(errormsg string) {
	ErrorLogger.Println(errormsg)
	log.Println(errormsg)
	os.Exit(1)
}

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

// Generate json output depending on the commadline flags
func GetJsonFromMap(list interface{}) string {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "   ")
	if err := encoder.Encode(list); err != nil {
		if jsonBytes, err := json.MarshalIndent(list, "", "  "); err != nil {
			ErrorLogger.Println(err)
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
		// return buffer.String()
	}
	return ""
}

// func getFieldFromStruct(v reflect.Type, field string) string {
// 	r := reflect.ValueOf(v)
// 	f := reflect.Indirect(r).FieldByName(field)
// 	return string(f.String())
// }

func GetYamlFromMap(list interface{}) string {
	d, err := yaml.Marshal(&list)
	if err != nil {
		ErrorLogger.Println("Convert map to Yaml failed", err)
	}
	return string(d)
}

func GetCsvFromMap(list interface{}, family string) {
	output := T_csvDoc{}
	if CmdParams.Output.Is || CmdParams.Output.All {
		headline := T_csvLine{"Family", "DataRange", "DataType", "Imagestream", "Image", "ImagestreamTag"}
		output = append(output, headline)
		for is, isMap := range list.(T_completeResults).AllIstags.Is {
			for image, shaMap := range isMap {
				for istag := range shaMap {
					line := T_csvLine{}
					line = append(line, family)
					line = append(line, "allIstags")
					line = append(line, "is")
					line = append(line, is)
					line = append(line, image)
					line = append(line, istag)
					output = append(output, line)
				}
			}
		}
	}
	if CmdParams.Output.Istag || CmdParams.Output.All {
		headline := T_csvLine{"Family", "DataRange", "DataType", "istag"} //, "Imagestream", "Tagname", "Namespace", "Link", "Date", "AgeInDays", "Image", "CommitAuthor", "CommitDate", "CommitId", "CommitRef", "Commitversion", "IsProdImage", "BuildNName", "BuildNamespace"}
		headline = append(headline, toArrayString(T_istag{}.Names())...)
		output = append(output, headline)
		for istagName, istagMap := range list.(T_completeResults).AllIstags.Istag {
			line := T_csvLine{}
			line = append(line, family)
			line = append(line, "allIstags")
			line = append(line, "istag")
			line = append(line, istagName)
			line = append(line, toArrayString(istagMap.Values())...)
			output = append(output, line)
		}
	}
	if CmdParams.Output.Image || CmdParams.Output.All {
		headline := T_csvLine{"Family", "DataRange", "DataType", "Image", "Istag", "Imagestream", "Namespace", "Link", "Date", "AgeInDays", "IsTagReferences"}
		output = append(output, headline)
		for shaName, shaMap := range list.(T_completeResults).AllIstags.Image {
			for istag, istagMap := range shaMap {
				line := T_csvLine{}
				line = append(line, family)
				line = append(line, "allIstags")
				line = append(line, "image")
				line = append(line, shaName)
				line = append(line, istag)
				line = append(line, istagMap.Imagestream)
				line = append(line, istagMap.Namespace)
				line = append(line, istagMap.Link)
				line = append(line, istagMap.Date)
				line = append(line, istagMap.AgeInDays)
				for tag := range istagMap.Istags {
					//  make a real copy of line !!!!
					copyOfLine := append([]string{}, line...)
					copyOfLine = append(copyOfLine, tag)
					output = append(output, copyOfLine)
				}
			}
		}
	}
	if CmdParams.Output.Used || CmdParams.Output.All {
		headline := T_csvLine{"Family", "DataRange", "DataType", "Imagestream", "Tag"} //, "UsedInNamespace", "Image", "UsedInCluster"}
		headline = append(headline, toArrayString(T_usedIstag{}.Names())...)
		output = append(output, headline)
		for is, isMap := range list.(T_completeResults).UsedIstags {
			for istag, istagArray := range isMap { //.(map[string][]map[string]string) {
				for _, istagMap := range istagArray {
					line := T_csvLine{}
					line = append(line, family)
					line = append(line, "usedistags")
					line = append(line, "is:tag")
					line = append(line, is)
					line = append(line, istag)
					line = append(line, toArrayString(istagMap.Values())...)
					output = append(output, line)
				}
			}
		}
	}
	w := csv.NewWriter(os.Stdout)
	if err := w.WriteAll(output.csvDoc()); err != nil {
		ErrorLogger.Println("writing csv failed" + err.Error())
	}
	// tablePrettyprint(output)
}

func GetTableFromMap(list interface{}, family string) {
	if CmdParams.Output.Is || CmdParams.Output.All {
		output := []table.Row{}
		headline := table.Row{"Imagestream " + family, "Image", "ImagestreamTag"}
		output = append(output, headline)
		for is, isMap := range list.(T_completeResults).AllIstags.Is {
			for image, shaMap := range isMap {
				for istag := range shaMap {
					line := table.Row{}
					// line = append(line, "allIstags")
					// line = append(line, "is")
					line = append(line, is)
					line = append(line, image)
					line = append(line, istag)
					output = append(output, line)
				}
			}
		}
		tablePrettyprint(output)
	}
	if CmdParams.Output.Istag || CmdParams.Output.All {
		output := []table.Row{}
		headline := table.Row{"istag " + family} //, "Imagestream", "Tagname", "Namespace", "Link", "Date", "AgeInDays", "Image", "CommitAuthor", "CommitDate", "CommitId", "CommitRef", "Commitversion", "IsProdImage", "BuildNName", "BuildNamespace"}
		headline = append(headline, toTableRow(T_istag{}.Names())...)
		output = append(output, headline)
		for istagName, istagMap := range list.(T_completeResults).AllIstags.Istag {
			line := table.Row{}
			// line = append(line, "allIstags")
			// line = append(line, "istag")
			line = append(line, istagName)
			line = append(line, toTableRow(istagMap.Values())...)
			output = append(output, line)
		}
		tablePrettyprint(output)
	}
	if CmdParams.Output.Image || CmdParams.Output.All {
		output := []table.Row{}
		headline := table.Row{"Image " + family, "Istag", "Imagestream", "Namespace", "Link", "Date", "AgeInDays", "IsTagReferences"}
		output = append(output, headline)
		for shaName, shaMap := range list.(T_completeResults).AllIstags.Image {
			for istag, istagMap := range shaMap {
				line := table.Row{}
				// line = append(line, "allIstags")
				// line = append(line, "image")
				line = append(line, shaName)
				line = append(line, istag)
				line = append(line, istagMap.Imagestream)
				line = append(line, istagMap.Namespace)
				line = append(line, istagMap.Link)
				line = append(line, istagMap.Date)
				line = append(line, istagMap.AgeInDays)
				for tag := range istagMap.Istags {
					//  make a real copy of line !!!!
					copyOfLine := append([]interface{}{}, line...)
					copyOfLine = append(copyOfLine, tag)
					output = append(output, copyOfLine)
				}
			}
		}
		tablePrettyprint(output)
	}
	if CmdParams.Output.Used || CmdParams.Output.All {
		output := []table.Row{}
		headline := table.Row{"Imagestream (used for " + family + ")", "Tag (used)"} //, "UsedInNamespace", "Image", "UsedInCluster"}
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
}

func toTableRow(arr ...interface{}) table.Row {
	o := table.Row{}
	for _, v := range arr {
		for _, w := range v.([]interface{}) {
			o = append(o, w)
		}
	}
	return o
}

func toArrayString(arr ...interface{}) []string {
	o := []string{}
	for _, v := range arr {
		for _, w := range v.([]interface{}) {
			o = append(o, w.(string))
		}
	}
	return o
}

func tablePrettyprint(out []table.Row) {
	if len(out) == 0 {
		return
	}
	_, height, err := terminal.GetSize(0)
	if err != nil {
		log.Println("failedt o get terminal size")
		height = 60
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
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
			// {Number: 4, AutoMerge: true},
		})
	}
	t.Render()
}
