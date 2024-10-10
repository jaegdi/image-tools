package ocrequest

import (
	"fmt"
)

var chanAllIsTags = make(chan T_ResultExistingIstagsOverAllClusters, 1)
var chanUsedIsTags = make(chan T_usedIstagsResult, 1)
var chanCompleteResults = make(chan T_completeResults, 1)
var chanInitAllImages = make(chan T_ImagesMapAllClusters, 1)

// getCause returns the cause string based on the command line parameters.
// It constructs the cause string by concatenating different filter options.
// If the delete pattern is specified, it appends "_pattern:_" followed by the pattern value.
// If the Isname filter is specified, it appends "_Isname:_" followed by the Isname value.
// If the Tagname filter is specified, it appends "_Tagname:_" followed by the Tagname value.
// If the Istagname filter is specified, it appends "_Istagname:_" followed by the Istagname value.
// If the Namespace filter is specified, it appends "_Namespace:_" followed by the Namespace value.
// If the Maxage filter is specified, it appends "_Maxage:_" followed by the Maxage value.
// If the Minage filter is specified, it appends "_Minage:_" followed by the Minage value.
// The constructed cause string is then returned.
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
		s = s + "_Maxage:_" + fmt.Sprintf("%d", CmdParams.Filter.Maxage)
	}
	if CmdParams.Filter.Minage != -1 {
		s = s + "_Minage:_" + fmt.Sprintf("%d", CmdParams.Filter.Minage)
	}
	return s
}

// CmdlineMode is a function that executes the command line mode of the ocrequest package.
// It performs various operations based on the command line parameters and returns the complete results.
// The function initializes the necessary channels and goroutines to retrieve and process the required data.
// It filters the results based on the specified filters and options, and performs additional operations if the delete option is enabled.
// Finally, it outputs the results in the desired format specified by the command line parameters.
// The function returns the complete results.
func CmdlineMode() T_completeResults {
	result := T_completeResults{}
	go InitAllImages(chanInitAllImages)
	go GetUsedIstagsForFamily(chanUsedIsTags)

	DebugMsg("Wait for chanInitAllImages")
	AllImages := <-chanInitAllImages
	DebugMsg("Image clusters:", len(AllImages))

	go GetAllIstagsForFamily(chanAllIsTags)

	DebugMsg("Wait for chanAllIsTags")
	result.AllIstags = <-chanAllIsTags

	DebugMsg("Wait for chanUsedIsTags")
	result.UsedIstags = <-chanUsedIsTags

	go PutShaIntoUsedIstags(chanCompleteResults, result)

	DebugMsg("Wait for filtered chanCompleteResults")

	result = <-chanCompleteResults
	DebugMsg("result: ", result)

	if CmdParams.Output.UnUsed {
		FilterUnusedIstags(&result)
	}
	// Filter results for output
	FilterAllIstags(&result)
	DebugMsg("result after filter: ", result)

	resultFamilies := T_completeResultsFamilies{}
	resultFamilies[CmdParams.Family] = result
	if CmdParams.Delete {
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
			DebugMsg("execute oc adm prune")
		} else {
			DebugMsg("run in dry run mode")
		}
	} else {
		switch {
		case CmdParams.Html:
			GetTextTableFromMap(result, CmdParams.Family)
		case CmdParams.Json:
			fmt.Println(GetJsonFromMap(resultFamilies))
		case CmdParams.Yaml:
			fmt.Println(GetYamlFromMap(resultFamilies))
		case CmdParams.Csv:
			GetCsvFromMap(result, CmdParams.Family)
		case (CmdParams.Table || CmdParams.TabGroup):
			GetTextTableFromMap(result, CmdParams.Family)
		case CmdParams.Options.ServerMode:
			return result
		}
	}

	return result
}
