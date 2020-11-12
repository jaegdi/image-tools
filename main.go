package main

import (
	"fmt"
	// "log"
	// "net/http"
	// _ "net/http/pprof"
	. "report-istags/ocrequest"
	// "sync"
)

func init() {
	Init()
}

func main() {
	// var wg sync.WaitGroup
	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()
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
	// wg.Add(1)
	// wg.Wait()
}
