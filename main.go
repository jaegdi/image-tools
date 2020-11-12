package main

import (
	"clean-istags/ocrequest"
	. "fmt"
	// _ "github.com/rakyll/gom/http"
	"net/http"
	_ "net/http/pprof"
)

func init() {
	ocrequest.Init()
}

func main() {
	go func() {
		ocrequest.InfoLogger.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	// var str string
	// Scanln(&str)
	// ocrequest.InfoLogger.Println(str)
	result := ocrequest.T_completeResults{}
	if !ocrequest.CmdParams.Output.Used {
		ocrequest.InitAllImagesOfCluster(ocrequest.CmdParams.Cluster)
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
	case ocrequest.CmdParams.Table || ocrequest.CmdParams.TabGroup:
		ocrequest.GetTableFromMap(result)
	}
	// Scanln(&str)
	// ocrequest.InfoLogger.Println(str)
	// ocrequest.Test_MergeNestedMaps()
}
