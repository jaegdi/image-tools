package ocrequest

import (
	"fmt"
)

var chanAllIsTags = make(chan T_ResultExistingIstagsOverAllClusters, 1)
var chanUsedIsTags = make(chan T_usedIstagsResult, 1)
var chanCompleteResults = make(chan T_completeResults, 1)
var chanInitAllImages = make(chan T_ImagesMapAllClusters, 1)

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
	if CmdParams.Filter.Maxage != -1 {
		s = s + "_Maxage:_" + string(CmdParams.Filter.Maxage)
	}
	if CmdParams.Filter.Minage != -1 {
		s = s + "_Minage:_" + string(CmdParams.Filter.Minage)
	}
	return s
}

func CmdlineMode() {
	result := T_completeResults{}
	go InitAllImages(chanInitAllImages)
	go GetUsedIstagsForFamily(chanUsedIsTags)

	DebugLogger.Println("Wait for chanInitAllImages")
	AllImages := <-chanInitAllImages
	DebugLogger.Println("Image clusters:", len(AllImages))

	go GetAllIstagsForFamily(chanAllIsTags)

	DebugLogger.Println("Wait for chanAllIsTags")
	result.AllIstags = <-chanAllIsTags

	DebugLogger.Println("Wait for chanUsedIsTags")
	result.UsedIstags = <-chanUsedIsTags

	go PutShaIntoUsedIstags(chanCompleteResults, result)

	DebugLogger.Println("Wait for filtered chanCompleteResults")
	result = <-chanCompleteResults
	DebugLogger.Println("result: ", result)

	if CmdParams.Output.UnUsed {
		FilterUnusedIstags(&result)
	}
	// Filter results for output
	FilterAllIstags(&result)
	DebugLogger.Println("result after filter: ", result)

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
		} else if CmdParams.DeleteOpts.NonBuild {
			FilterNonbuildIstagsToDelete(
				resultFamilies,
				CmdParams.Family,
				CmdParams.Cluster,
				CmdParams.DeleteOpts.MinAge)
		} else {
			DebugLogger.Println(
				"\n--main--::\n",
				"filter minAge: '"+fmt.Sprint(CmdParams.DeleteOpts.MinAge)+"'\n",
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

		if CmdParams.DeleteOpts.Confirm {
			DebugLogger.Println("execute oc adm prune")
		} else {
			DebugLogger.Println("run in dry run mode")
		}
	}
}
