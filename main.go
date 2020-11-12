package main

import (
	"fmt"
	. "report-istags/ocrequest"
	// _ "github.com/rakyll/gom/http"
	// "net/http"
	// _ "net/http/pprof"
)

func init() {
	Init()
}

func main() {
	// go func() {
	// 	InfoLogger.Println(http.ListenAndServe("localhost:6060", nil))
	// }()
	// var str string
	// Scanln(&str)
	// InfoLogger.Println(str)
	result := T_completeResults{}
	if !CmdParams.Output.Used {
		InitAllImagesOfCluster(CmdParams.Cluster)
		allIsTags := (GetAllIstagsForFamilyInCluster())
		filteredIsTags := FilterAllIstags(allIsTags)
		result.AllIstags = filteredIsTags
	}
	if CmdParams.Output.All || CmdParams.Output.Used {
		usedIsTags := (GetUsedIstagsForFamily())
		result.UsedIstags = usedIsTags
	}
	resultFamilies := T_completeResultsFamilies{}
	resultFamilies[CmdParams.Family] = result
	switch {
	case CmdParams.Json:
		fmt.Println(GetJsonFromMap(resultFamilies))
	case CmdParams.Yaml:
		fmt.Println(GetYamlFromMap(resultFamilies))
	case CmdParams.Csv:
		GetCsvFromMap(result, CmdParams.Family)
	case CmdParams.Table || CmdParams.TabGroup:
		GetTableFromMap(result, CmdParams.Family)
	}
	// Scanln(&str)
	// InfoLogger.Println(str)
	// Test_MergeNestedMaps()
}
