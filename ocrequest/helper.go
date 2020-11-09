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

func GetCsvFromMap(list interface{}) {
	output := T_csvDoc{}
	headline := T_csvLine{"DataRange", "DataType", "Imagestream", "Image", "ImagestreamTag"}
	output = append(output, headline)
	for is, isMap := range list.(T_completeResults).AllIstags.Is {
		for sha, shaMap := range isMap {
			for istag := range shaMap {
				line := T_csvLine{}
				line = append(line, "allIstags")
				line = append(line, "is")
				line = append(line, is)
				line = append(line, sha)
				line = append(line, istag)
				output = append(output, line)
			}
		}
	}
	// tmp := T_istag{}
	headline = T_csvLine{"DataRange", "DataType", "istag", "Imagestream", "Tagname", "Namespace", "Link", "Date", "AgeInDays", "Image", "CommitAuthor", "CommitDate", "CommitId", "CommitRef", "Commitversion", "IsProdImage", "BuildNName", "BuildNamespace"}
	// v := reflect.ValueOf(tmp)
	// typeOfS := v.Type()
	// istagHeadline := []string{}
	// for i := 0; i < v.NumField(); i++ {
	// 	if typeOfS.Field(i).Name != "Build" {
	// 		istagHeadline = append(istagHeadline, typeOfS.Field(i).Name)
	// 	}
	// }
	// b := reflect.ValueOf(tmp.Build)
	// typeOfB := b.Type()
	// buildHeadline := []string{}
	// for i := 0; i < b.NumField(); i++ {
	// 	buildHeadline = append(buildHeadline, typeOfB.Field(i).Name)
	// }
	// headline = append(headline, istagHeadline...)
	// headline = append(headline, buildHeadline...)
	output = append(output, headline)
	for istagName, istagMap := range list.(T_completeResults).AllIstags.Istag {
		line := T_csvLine{}
		line = append(line, "allIstags")
		line = append(line, "istag")
		line = append(line, istagName)
		line = append(line, istagMap.List()...)
		// line = append(line, istagMap.Imagestream)
		// line = append(line, istagMap.Tagname)
		// line = append(line, istagMap.Namespace)
		// line = append(line, istagMap.Link)
		// line = append(line, istagMap.Date)
		// line = append(line, istagMap.AgeInDays)
		// line = append(line, istagMap.Sha)
		// v := reflect.ValueOf(istagMap)
		// for i := 0; i < v.NumField(); i++ {
		// 	line = append(line, v.Field(i).String())
		// }
		// for _, v := range istagHeadline {
		// 	line = append(line, getFieldFromStruct(&istagMap, v))
		// }
		// line = append(line, istagMap.Build.CommitAuthor)
		// line = append(line, istagMap.Build.CommitDate)
		// line = append(line, istagMap.Build.CommitId)
		// line = append(line, istagMap.Build.CommitRef)
		// line = append(line, istagMap.Build.CommitVersion)
		// line = append(line, istagMap.Build.IsProdImage)
		// line = append(line, istagMap.Build.Name)
		// line = append(line, istagMap.Build.Namespace)
		// build := reflect.ValueOf(istagMap.Build)
		// for i := 0; i < build.NumField(); i++ {
		// 	line = append(line, build.Field(i).String())
		// }
		// for _, v := range buildHeadline {
		// 	line = append(line, getFieldFromStruct(&istagMap.Build, v))
		// }
		output = append(output, line)
	}
	headline = T_csvLine{"DataRange", "DataType", "Image", "Istag", "Imagestream", "Namespace", "Link", "Date", "AgeInDays", "IsTagReferences"}
	output = append(output, headline)
	for shaName, shaMap := range list.(T_completeResults).AllIstags.Sha {
		for istag, istagMap := range shaMap {
			line := T_csvLine{}
			line = append(line, "allIstags")
			line = append(line, "sha")
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
	headline = T_csvLine{"DataType", "Imagestream", "Tag", "UsedInNamespace", "Image", "UsedInCluster"}
	output = append(output, headline)
	for is, isMap := range list.(T_completeResults).UsedIstags {
		for istag, istagArray := range isMap { //.(map[string][]map[string]string) {
			for _, istagMap := range istagArray {
				line := T_csvLine{}
				line = append(line, "usedistags")
				line = append(line, is)
				line = append(line, istag)
				line = append(line, istagMap.UsedInNamespace)
				line = append(line, istagMap.Sha)
				line = append(line, istagMap.Cluster)
				output = append(output, line)
			}
		}
	}
	w := csv.NewWriter(os.Stdout)
	// for _, line := range output {
	// 	w.Write(line)
	// }
	if err := w.WriteAll(output.csvDoc()); err != nil {
		ErrorLogger.Println("writing csv failed" + err.Error())
	}
}
