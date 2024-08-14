// Package ocrequest provides primitives to query an oc-cluster for istags
// calculate and filter them. And provides a json output primitive.
package ocrequest

import (
	"encoding/json"
	"strings"
)

var IsNamesForFamily T_IsNamesForFamily

// joinShaStreams join keys of a map to an array.
func joinShaStreams(mymap T_Istags_List) []T_istagName {
	keys := []T_istagName{}
	for k := range mymap {
		keys = append(keys, k)
	}
	return keys
}

// appendJoinedNamesToImagestreams get as params the imageStream-map and then joinedNames of istags and
// put them on the map under imagestreamMap.imagestream.image .
func appendJoinedNamesToImagestreams(istream T_resIs, imagestreamName T_isName, image T_shaName, joinedNames []T_istagName) T_resIs {
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
func InitIsNamesForFamily(family T_familyName) {
	// cluster := FamilyNamespaces[CmdParams.Family].Buildstage
	isResult := map[string]interface{}{}
	result := make(T_IsNamesForFamily)
	result[family] = make(map[T_isName]bool)
	for _, cluster := range FamilyNamespaces[CmdParams.Family].Buildstages {
		for _, ns := range FamilyNamespaces[family].ImageNamespaces[cluster] {
			if ns == "openshift" {
				continue
			}
			isJson := ocGetCall(cluster, ns, "imagestreams", "")
			if err := json.Unmarshal([]byte(isJson), &isResult); err != nil {
				ErrorLogger.Println("Unmarshal imagestreams. " + err.Error())
			}
			if isResult["items"] != nil {
				for _, v := range isResult["items"].([]interface{}) {
					val := T_isName(v.(map[string]interface{})["metadata"].(map[string]interface{})["name"].(string))
					result[family][val] = true
				}
			}
		}
	}
	IsNamesForFamily = result
}

func setBuildLabels(buildLabelsMap map[string]interface{}) T_istagBuildLabels {
	buildLabels := T_istagBuildLabels{}
	buildLabelsJSON := []byte(GetJsonFromMap(buildLabelsMap))
	if err := json.Unmarshal(buildLabelsJSON, &buildLabels); err != nil {
		ErrorLogger.Println("Unmarshal unescaped String", err)
	}
	return buildLabels
}

// OcGetAllIstagsOfNamespace generates a map of all istags
// selected by (cluster, namespace) and append it to result map
// and return the result map
func OcGetAllIstagsOfNamespace(result T_result, cluster T_clName, namespace T_nsName) T_result {
	istagsJson := ocGetCall(cluster, namespace, "imagestreamtags", "")
	DebugMsg("istagJson:", istagsJson)
	var istagsMap map[string]interface{}
	if err := json.Unmarshal([]byte(istagsJson), &istagsMap); err != nil {
		// logfix
		ErrorLogger.Println("Error: unmarshal imagestreamtags." + istagsJson)
		ErrorLogger.Println("Error: unmarshal imagestreamtags." + err.Error())
		return T_result{}
	}
	resultIstag := make(T_resIstag)
	resultSha := make(T_resSha)
	resultIstream := make(T_resIs)
	// var isLink string

	var itemsMap []interface{}
	if istagsMap["items"] != nil {
		itemsMap = istagsMap["items"].([]interface{})
	}
	shaNames := make(T_shaNames)
	shaStreams := make(T_shaStreams)
	var metadata map[string]interface{}
	var imageMetadata map[string]interface{}

	for _, content := range itemsMap {
		metadata = content.(map[string]interface{})["metadata"].(map[string]interface{})
		imageMetadata = content.(map[string]interface{})["image"].(map[string]interface{})["metadata"].(map[string]interface{})
		istagname := T_istagName(metadata["name"].(string))
		isNamespace := T_nsName(metadata["namespace"].(string))
		// if metadata["selfLink"] != nil {
		// 	isLink = metadata["selfLink"].(string)
		// } else {
		// 	isLink = ""
		// }
		isDate := metadata["creationTimestamp"].(string)
		sha := T_shaName(imageMetadata["name"].(string))
		if CmdParams.Filter.Imagename != "" && sha != CmdParams.Filter.Imagename {
			continue
		}
		if CmdParams.Filter.Istagname != "" && istagname != CmdParams.Filter.Istagname && !CmdParams.FilterReg.Istagname.MatchString(string(istagname)) {
			continue
		}

		buildLabelsMap := T_istagBuildLabels{}
		if CmdParams.Options.Debug {
			DebugMsg("IsTag: "+istagname, "ImagesMap: ", ImagesMap)
		}
		if len(ImagesMap[cluster]) > 0 && ImagesMap[cluster][sha.str()].(map[string]interface{})["dockerImageMetadata"].(map[string]interface{})["Config"].(map[string]interface{})["Labels"] != nil {
			buildLabelsMap.Set(ImagesMap[cluster][sha.str()].(map[string]interface{})["dockerImageMetadata"].(map[string]interface{})["Config"].(map[string]interface{})["Labels"].(map[string]interface{}))
		}
		imagestreamfields := strings.Split(istagname.str(), `:`)
		imagestreamName := T_isName(imagestreamfields[0])
		if !IsNamesForFamily[CmdParams.Family][imagestreamName] {
			continue
		}
		tagName := T_tagName(imagestreamfields[1])
		isAge := ageInDays(isDate)

		if !matchIsIstagToFilterParams(imagestreamName, tagName, istagname, isNamespace, isAge) {
			continue
		}

		shaNames.Add(sha, T_istagName(isNamespace.str()+"/"+istagname.str()))

		myIstag := T_istag{
			Imagestream: imagestreamName,
			Tagname:     tagName,
			Namespace:   isNamespace,
			// Link:        isLink,
			Date:      isDate,
			AgeInDays: isAge,
			Image:     sha,
			Build:     buildLabelsMap,
		}

		shaStreams.Add(imagestreamName, sha, myIstag)

		mySha := map[T_istagName]T_sha{
			istagname: {
				Istags:      shaNames[sha],
				Imagestream: imagestreamName,
				Namespace:   isNamespace,
				// Link:        isLink,
				Date:      isDate,
				AgeInDays: isAge,
			},
		}

		joinedNames := joinShaStreams(shaNames[sha])
		resultIstream = appendJoinedNamesToImagestreams(resultIstream, imagestreamName, sha, joinedNames)
		if resultIstag[istagname] == nil {
			resultIstag[istagname] = map[T_nsName]T_istag{}
		}
		// if resultIstag[istagname][isNamespace] == T_istag{} {
		// 	resultIstag[istagname][isNamespace] = T_istag{}
		// }
		resultIstag[istagname][isNamespace] = myIstag

		if resultSha[sha] == nil {
			resultSha[sha] = make(map[T_istagName]T_sha)
		}
		// t := map[string]T_sha{}
		t := resultSha[sha]
		MergoNestedMaps(&t, mySha)
		resultSha[sha] = t
	}
	tmp_result := T_result{
		Istag: resultIstag,
		Image: resultSha,
		Is:    resultIstream,
	}

	MergoNestedMaps(&result, tmp_result)
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
		for _, cluster := range FamilyNamespaces[CmdParams.Family].Stages {
			if namespace == "" {
				for _, ns := range FamilyNamespaces[family].ImageNamespaces[cluster] {
					r := T_ResultExistingIstagsOverAllClusters{cluster: OcGetAllIstagsOfNamespace(result[cluster], cluster, ns)}
					MergoNestedMaps(&result, r)

				}
			} else {
				result = T_ResultExistingIstagsOverAllClusters{cluster: OcGetAllIstagsOfNamespace(result[cluster], cluster, namespace)}
			}
		}
	}
	c <- result
}
