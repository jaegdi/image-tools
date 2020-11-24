package main

import (
	"fmt"
	// "net/http"
	// _ "net/http/pprof"
	. "report-istags/ocrequest"
	// "sync"
)

func init() {
	Init()
}

var chanAllIsTags = make(chan T_ResultExistingIstagsOverAllClusters, 1)
var chanUsedIsTags = make(chan T_usedIstagsResult, 1)
var chanCompleteResults = make(chan T_completeResults, 1)
var chanInitAllImages = make(chan string, 1)
var LogfileName string

func main() {
	// var wg sync.WaitGroup
	// if CmdParams.Options.Profiler {
	// 	go func() {
	// 		LogMsg(http.ListenAndServe("localhost:6060", nil))
	// 	}()
	// }
	result := T_completeResults{}
	go InitAllImages(chanInitAllImages)
	go GetUsedIstagsForFamily(chanUsedIsTags)

	LogMsg("Wait for chanInitAllImages")
	LogMsg(<-chanInitAllImages)

	go GetAllIstagsForFamily(chanAllIsTags)

	LogMsg("Wait for chanAllIsTags")
	result.AllIstags = <-chanAllIsTags

	LogMsg("Wait for chanUsedIsTags")
	result.UsedIstags = <-chanUsedIsTags

	go PutShaIntoUsedIstags(chanCompleteResults, result)

	// filteredIsTags := T_ResultExistingIstagsOverAllClusters{}
	// MergoNestedMaps(&filteredIsTags, allIsTags)
	// result.AllIstags = filteredIsTags

	LogMsg("Wait for filtered chanUsedIsTags")
	result = <-chanCompleteResults

	// Filter results for output
	FilterAllIstags(&result)

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
	// if CmdParams.Options.Profiler {
	// 	wg.Add(1)
	// 	wg.Wait()
	// }
}
