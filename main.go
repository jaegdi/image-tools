package main

import (
	"clean-istags/ocrequest"
	. "fmt"
)

var familyNamespaces = ocrequest.T_famNs{
	"pkp": []string{"ms-jenkins", "openshift", "images-pkp"},
	"ssp": []string{"ssp-jenkins", "images-ssp"},
	"aps": []string{"aps-jenkins", "images-aps"},
	"fpc": []string{"fpc-jenkins", "images-fpc"},
}

func main() {
	flags := ocrequest.EvalFlags(familyNamespaces)

	result := map[string]interface{}{}
	allIsTags := (ocrequest.GetAllIstagsForFamilyInCluster(flags, familyNamespaces))
	allIsTags = ocrequest.FilterAllIstags(allIsTags, flags)

	usedIsTags := (ocrequest.GetUsedIstagsForFamilyInCluster(flags, familyNamespaces))

	result["all-istags"] = allIsTags
	result["used-istags"] = usedIsTags

	Println(ocrequest.GetJsonFromMap(result, flags))
}
