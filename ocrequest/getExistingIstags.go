// Package ocrequest provides primitives to query an oc-cluster for istags
// calculate and filter them. And provides a json output primitive.
package ocrequest

import (
	"encoding/json"
	"fmt"
	"strings"
)

var IsNamesForFamily T_IsNamesForFamily

// joinShaStreams join keys of a map to an array.
func joinShaStreams(mymap map[string]bool) []string {
	keys := []string{}
	for k := range mymap {
		keys = append(keys, k)
	}
	return keys
}

// appendJoinedNamesToImagestreams get as params the imageStream-map and then joinedNames of istags and
// put them on the map under imagestreamMap.imagestream.image .
func appendJoinedNamesToImagestreams(istream T_resIs, imagestreamName string, image string, joinedNames []string) T_resIs {
	if istream[imagestreamName] == nil {
		istream[imagestreamName] = T_is{}
	}
	if istream[imagestreamName][image] == nil {
		istream[imagestreamName][image] = T_isShaTagnames{}
	}
	for _, v := range joinedNames {
		istream[imagestreamName][image][v] = true
	}
	return istream
}

// InitIsNamesForFamily initializes the package var IsNamesForFamily with all imagestreams from
// the build namespaces of the family.
func InitIsNamesForFamily(family string) {
	cluster := Clusters.Buildstage
	isResult := map[string]interface{}{}
	result := make(T_IsNamesForFamily)
	result[family] = make(map[string]bool)
	for _, ns := range FamilyNamespaces[family][cluster] {
		if ns == "openshift" {
			continue
		}
		isJson := ocGetCall(cluster, ns, "imagestreams", "")
		if err := json.Unmarshal([]byte(isJson), &isResult); err != nil {
			LogError("Unmarshal imagestreams. " + err.Error())
		}
		if isResult["items"] != nil {
			for _, v := range isResult["items"].([]interface{}) {
				val := v.(map[string]interface{})["metadata"].(map[string]interface{})["name"].(string)
				result[family][val] = true
			}
		}
	}
	IsNamesForFamily = result
}

// func setBuildLabels(buildLabelsMap map[string]interface{}) T_istagBuildLabels {
// 	buildLabels := T_istagBuildLabels{}
// 	buildLabelsJSON := []byte(GetJsonFromMap(buildLabelsMap))
// 	if err := json.Unmarshal(buildLabelsJSON, &buildLabels); err != nil {
// 		LogError("Unmarshal unescaped String", err)
// 	}
// 	return buildLabels
// }

// OcGetAllIstagsOfNamespace generates a map of all istags
// selected by (cluster, namespace) and append it to result map
// and return the result map
func OcGetAllIstagsOfNamespace(result T_result, cluster string, namespace string) T_result {
	istagsJson := ocGetCall(cluster, namespace, "imagestreamtags", "")
	var istagsMap map[string]interface{}
	if err := json.Unmarshal([]byte(istagsJson), &istagsMap); err != nil {
		LogError("unmarshal imagestreamtags.\n" + istagsJson + "\n" + err.Error())
		return T_result{}
	}
	resultIstag := make(T_resIstag)
	resultSha := make(T_resSha)
	resultIstream := make(T_resIs)

	itemsMap := istagsMap["items"].([]interface{})
	shaNames := make(T_shaNames)
	shaStreams := make(T_shaStreams)
	var metadata map[string]interface{}
	var imageMetadata map[string]interface{}

	for _, content := range itemsMap {
		metadata = content.(map[string]interface{})["metadata"].(map[string]interface{})
		imageMetadata = content.(map[string]interface{})["image"].(map[string]interface{})["metadata"].(map[string]interface{})
		istagname := metadata["name"].(string)
		if CmdParams.Filter.Istagname != "" && istagname != CmdParams.Filter.Istagname {
			continue
		}
		isNamespace := metadata["namespace"].(string)
		isLink := metadata["selfLink"].(string)
		isDate := metadata["creationTimestamp"].(string)
		sha := imageMetadata["name"].(string)
		if CmdParams.Filter.Imagename != "" && sha != CmdParams.Filter.Imagename {
			continue
		}

		buildLabelsMap := T_istagBuildLabels{}
		if ImagesMap[cluster][sha].(map[string]interface{})["dockerImageMetadata"].(map[string]interface{})["Config"].(map[string]interface{})["Labels"] != nil {
			buildLabelsMap.Set(ImagesMap[cluster][sha].(map[string]interface{})["dockerImageMetadata"].(map[string]interface{})["Config"].(map[string]interface{})["Labels"].(map[string]interface{}))
		}
		imagestreamfields := strings.Split(istagname, `:`)
		imagestreamName := imagestreamfields[0]
		if !IsNamesForFamily[CmdParams.Family][imagestreamName] {
			continue
		}
		tagName := imagestreamfields[1]
		isAge := fmt.Sprintf("%v", ageInDays(isDate))
		if CmdParams.Filter.Isname != "" && imagestreamName != CmdParams.Filter.Isname {
			continue
		}
		if CmdParams.Filter.Tagname != "" && tagName != CmdParams.Filter.Tagname {
			continue
		}

		shaNames.Add(sha, isNamespace+"/"+istagname)

		myIstag := T_istag{
			Imagestream: imagestreamName,
			Tagname:     tagName,
			Namespace:   isNamespace,
			Link:        isLink,
			Date:        isDate,
			AgeInDays:   isAge,
			Image:       sha,
			Build:       buildLabelsMap,
		}

		shaStreams.Add(imagestreamName, sha, myIstag)

		mySha := map[string]T_sha{
			istagname: {
				Istags:      shaNames[sha],
				Imagestream: imagestreamName,
				Namespace:   isNamespace,
				Link:        isLink,
				Date:        isDate,
				AgeInDays:   isAge,
			},
		}

		joinedNames := joinShaStreams(shaNames[sha])
		resultIstream = appendJoinedNamesToImagestreams(resultIstream, imagestreamName, sha, joinedNames)
		resultIstag[istagname] = myIstag

		if resultSha[sha] == nil {
			resultSha[sha] = make(map[string]T_sha)
		}
		t := map[string]T_sha{}
		MergoNestedMaps(&t, resultSha[sha], mySha)
		resultSha[sha] = t
	}
	tmp_result := T_result{
		Istag: resultIstag,
		Image: resultSha,
		Is:    resultIstream,
	}

	t := T_result{}
	MergoNestedMaps(&t, result, tmp_result)
	result = t
	n_istags := len(result.Istag)
	n_shas := len(result.Image)
	n_is := len(result.Is)
	result.Report = T_resReport{
		Anz_ImageStreamTags: n_istags,
		Anz_Images:          n_shas,
		Anz_ImageStreams:    n_is,
	}
	return result
}

// GetAllIstagsForFamilyInCluster generates a map of all istags
// selected by (family, cluster, namespace) in CmdParams
// and return a map with the results
func GetAllIstagsForFamily(c chan T_ResultExistingIstagsOverAllClusters) {
	family := CmdParams.Family
	namespace := CmdParams.Filter.Namespace
	var result = T_ResultExistingIstagsOverAllClusters{}
	if Multiproc {
		result = goGetExistingIstagsForFamilyInAllClusters(family)
	} else {
		for _, cluster := range Clusters.Stages {
			if namespace == "" {
				for _, ns := range FamilyNamespaces[family][cluster] {
					r := T_ResultExistingIstagsOverAllClusters{cluster: OcGetAllIstagsOfNamespace(result[cluster], cluster, ns)}
					t := T_ResultExistingIstagsOverAllClusters{}
					MergoNestedMaps(&t, result, r)
					result = t

				}
			} else {
				result = T_ResultExistingIstagsOverAllClusters{cluster: OcGetAllIstagsOfNamespace(result[cluster], cluster, namespace)}
			}
		}
	}
	c <- result
}
