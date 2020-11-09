package main

import (
	"clean-istags/ocrequest"
	. "fmt"
	_ "net/http/pprof"
)

func init() {
	ocrequest.Init()
}

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	ocrequest.EvalFlags()
	ocrequest.GetIsNamesForFamily(ocrequest.CmdParams.Family)
	result := ocrequest.T_completeResults{}
	if ocrequest.CmdParams.Output.All ||
		ocrequest.CmdParams.Output.Is ||
		ocrequest.CmdParams.Output.Istag ||
		ocrequest.CmdParams.Output.Sha {
		ocrequest.OcGetAllImagesOfCluster(ocrequest.CmdParams.Cluster)
		// ocrequest.GetIsNamesForFamily(ocrequest.CmdParams.Family)
		allIsTags := (ocrequest.GetAllIstagsForFamilyInCluster())
		filteredIsTags := ocrequest.FilterAllIstags(allIsTags)
		result.AllIstags = filteredIsTags
	}
	if ocrequest.CmdParams.Output.All || ocrequest.CmdParams.Output.Used {
		usedIsTags := (ocrequest.GetUsedIstagsForFamily())
		result.UsedIstags = usedIsTags
	}

	switch {
	case ocrequest.CmdParams.Json:
		Println(ocrequest.GetJsonFromMap(result))
	case ocrequest.CmdParams.Yaml:
		Println(ocrequest.GetYamlFromMap(result))
	case ocrequest.CmdParams.Csv:
		ocrequest.GetCsvFromMap(result)
	}

	// ocrequest.Test_MergeNestedMaps()
}
