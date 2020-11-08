package main

import (
	"clean-istags/ocrequest"
	. "fmt"
)

func init() {
	ocrequest.Init()
}

func main() {
	ocrequest.EvalFlags()
	ocrequest.GetIsNamesForFamily(ocrequest.CmdParams.Family)
	ocrequest.OcGetAllImagesOfCluster(ocrequest.CmdParams.Cluster)
	// ocrequest.GetIsNamesForFamily(ocrequest.CmdParams.Family)
	result := ocrequest.T_completeResults{}
	allIsTags := (ocrequest.GetAllIstagsForFamilyInCluster())
	filteredIsTags := ocrequest.FilterAllIstags(allIsTags)
	result.AllIstags = filteredIsTags

	if ocrequest.CmdParams.Output.All || ocrequest.CmdParams.Output.Used {
		usedIsTags := (ocrequest.GetUsedIstagsForFamilyInCluster())
		result.UsedIstags = usedIsTags
	}

	Println(ocrequest.GetJsonFromMap(result))

	// ocrequest.Test_MergeNestedMaps()
}
