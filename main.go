package main

import (
	"clean-istags/ocrequest"
	. "fmt"
)

type completeResults struct {
	AllIstags  ocrequest.T_result
	UsedIstags ocrequest.T_runningObjects
}


func main() {
	ocrequest.EvalFlags()

	result := completeResults{}
	allIsTags := (ocrequest.GetAllIstagsForFamilyInCluster())
	filteredIsTags := ocrequest.FilterAllIstags(allIsTags)

	usedIsTags := (ocrequest.GetUsedIstagsForFamilyInCluster())

	result.AllIstags = filteredIsTags
	result.UsedIstags = usedIsTags

	Println(ocrequest.GetJsonFromMap(result))
}
