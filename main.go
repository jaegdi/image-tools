package main

import (
	"fmt"
	// "net/http"
	// _ "net/http/pprof"
	. "image-tools/ocrequest"
	// "sync"
)

func init() {
	Init()
}

var chanAllIsTags = make(chan T_ResultExistingIstagsOverAllClusters, 1)
var chanUsedIsTags = make(chan T_usedIstagsResult, 1)
var chanCompleteResults = make(chan T_completeResults, 1)
var chanInitAllImages = make(chan T_ImagesMapAllClusters, 1)
var LogfileName string

func getCause() string {
	s := "filter"
	if CmdParams.DeleteOpts.Pattern != "" {
		s = s + "_pattern:_" + CmdParams.DeleteOpts.Pattern
	}
	if CmdParams.Filter.Isname != "" {
		s = s + "_Isname:_" + string(CmdParams.Filter.Isname)
	}
	if CmdParams.Filter.Tagname != "" {
		s = s + "_Tagname:_" + string(CmdParams.Filter.Tagname)
	}
	if CmdParams.Filter.Istagname != "" {
		s = s + "_Istagname:_" + string(CmdParams.Filter.Istagname)
	}
	if CmdParams.Filter.Namespace != "" {
		s = s + "_Namespace:_" + string(CmdParams.Filter.Namespace)
	}
	return s
}

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

	LogDebug("Wait for chanInitAllImages")
	AllImages := <-chanInitAllImages
	LogDebug("Image clusters:", len(AllImages))

	go GetAllIstagsForFamily(chanAllIsTags)

	LogDebug("Wait for chanAllIsTags")
	result.AllIstags = <-chanAllIsTags

	LogDebug("Wait for chanUsedIsTags")
	result.UsedIstags = <-chanUsedIsTags

	go PutShaIntoUsedIstags(chanCompleteResults, result)

	LogDebug("Wait for filtered chanCompleteResults")
	result = <-chanCompleteResults

	if CmdParams.Output.UnUsed {
		FilterUnusedIstags(&result)
	}
	// Filter results for output
	FilterAllIstags(&result)

	resultFamilies := T_completeResultsFamilies{}
	resultFamilies[CmdParams.Family] = result
	if !CmdParams.Delete {
		switch {
		case CmdParams.Json:
			fmt.Println(GetJsonFromMap(resultFamilies))
		case CmdParams.Yaml:
			fmt.Println(GetYamlFromMap(resultFamilies))
		case CmdParams.Csv:
			GetCsvFromMap(result, CmdParams.Family)
		case (CmdParams.Table || CmdParams.TabGroup):
			GetTableFromMap(result, CmdParams.Family)
		}
	} else {
		if CmdParams.DeleteOpts.Snapshots {
			FilterIstagsToDelete(
				resultFamilies,
				CmdParams.Family,
				CmdParams.Cluster,
				`snapshot|SNAPSHOT|\:PR-|\:[[:digit:]](\.[[:digit:]]+){2}\-202[[:digit:]]{5}\.[[:digit:]]{6}\-[[:digit:]]`,
				CmdParams.DeleteOpts.MinAge,
				"snapshots and pull requests")
		}
		if CmdParams.DeleteOpts.Pattern != "" ||
			CmdParams.Filter.Isname != "" ||
			CmdParams.Filter.Tagname != "" ||
			CmdParams.Filter.Istagname != "" ||
			CmdParams.Filter.Namespace != "" {
			LogDebug(
				"main::",
				"filter pattern: '"+CmdParams.DeleteOpts.Pattern+"'\n",
				"filter Isname: '"+CmdParams.Filter.Isname+"'\n",
				"filter Tagname: '"+CmdParams.Filter.Tagname+"'\n",
				"filter Istagname: '"+CmdParams.Filter.Istagname+"'\n",
				"filter Namespace: '"+CmdParams.Filter.Namespace+"'")
			FilterIstagsToDelete(
				resultFamilies,
				CmdParams.Family,
				CmdParams.Cluster,
				CmdParams.DeleteOpts.Pattern,
				CmdParams.DeleteOpts.MinAge,
				getCause())
		}
		if CmdParams.DeleteOpts.NonBuild {
			FilterNonbuildIstagsToDelete(
				resultFamilies,
				CmdParams.Family,
				CmdParams.Cluster,
				CmdParams.DeleteOpts.MinAge)
		}
		if CmdParams.DeleteOpts.Confirm {
			LogDebug("execute oc adm prune")
		} else {
			LogDebug("run in dry run mode")
		}
		// if CmdParams.Options.Profiler {
		// 	wg.Add(1)
		// 	wg.Wait()
		// }
	}
}
