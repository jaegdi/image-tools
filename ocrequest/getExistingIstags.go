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
func joinShaStreams(mymap map[string]interface{}) []string {
	keys := []string{}
	for k := range mymap {
		keys = append(keys, k)
	}
	return keys
}

// addNamesToShaNames initialize the sha-names map if it is nil and add new names to this map.
func addNamesToShaNames(a T_shaNames, key string, b string) T_shaNames {
	if a == nil {
		a = T_shaNames{}
	}
	if a[key] == nil {
		a[key] = map[string]interface{}{}
	}
	a[key][b] = true
	return a
}

// addNamesToShaStreams initialize the sha-streams map if it is nil and add new names to this map.
func addNamesToShaStreams(a T_shaStreams, is string, sha string, istag T_istag) T_shaStreams {
	if a == nil {
		a = T_shaStreams{}
	}
	if a[is] == nil {
		a[is] = T_resIstag{}
	}
	if a[is][sha] == (T_istag{}) {
		a[is][sha] = T_istag{}
	}
	// for k, v := range istag {
	a[is][sha] = istag
	// }
	return a
}

// appendJoinedNamesToImagestreams get as params the imageStream-map and then joinedNames of istags and
// put them on the map under imagestreamMap.imagestream.sha .
func appendJoinedNamesToImagestreams(istream T_resIs, imagestreamName string, sha string, joinedNames []string) T_resIs {
	if istream[imagestreamName] == nil {
		istream[imagestreamName] = T_is{}
	}
	if istream[imagestreamName][sha] == nil {
		istream[imagestreamName][sha] = T_isShaTagnames{}
	}
	for _, v := range joinedNames {
		istream[imagestreamName][sha][v] = true
	}
	return istream
}

func GetIsNamesForFamily(family string) {
	cluster := "cid"
	isResult := map[string]interface{}{}
	result := make(T_IsNamesForFamily)
	result[family] = make(map[string]bool)
	for _, ns := range FamilyNamespaces[family] {
		if ns == "openshift" {
			continue
		}
		isJson := ocGetCall(cluster, ns, "imagestreams", "")
		if err := json.Unmarshal([]byte(isJson), &isResult); err != nil {
			ErrorLogger.Println("Unmarshal imagestreams. " + err.Error())
		}
		for _, v := range isResult["items"].([]interface{}) {
			val := v.(map[string]interface{})["metadata"].(map[string]interface{})["name"].(string)
			result[family][val] = true
		}
	}
	IsNamesForFamily = result
}

func setBuildLabels(buildLabelsMap map[string]interface{}) T_istagBuildLabels {
	buildLabels := T_istagBuildLabels{}
	// buildLabelsJSON, err := json.MarshalIndent(buildLabelsMap, "", "    ")
	// if err != nil {
	// 	ErrorLogger.Println("Marshal buildLabels to jSON failed in setBuildLabels", err)
	// }
	buildLabelsJSON := []byte(GetJsonFromMap(buildLabelsMap))

	// fmt.Println(string(buildLabelsUnescJSON))
	if err := json.Unmarshal(buildLabelsJSON, &buildLabels); err != nil {
		ErrorLogger.Println("Unmarshal unescaped String", err)
	}
	return buildLabels
}

func OcGetAllIstagsOfNamespace(result T_result, cluster string, namespace string) T_result {
	istagsJson := ocGetCall(cluster, namespace, "imagestreamtags", "")

	var istagsMap map[string]interface{}
	if err := json.Unmarshal([]byte(istagsJson), &istagsMap); err != nil {
		ErrorLogger.Println("unmarshal imagestreamtags." + err.Error())
	}

	// result := make(map[string]interface{})s
	resultIstag := make(T_resIstag)
	resultSha := make(T_resSha)
	resultIstream := make(T_resIs)

	itemsMap := istagsMap["items"].([]interface{})
	shaNames := make(T_shaNames)
	shaStreams := make(T_shaStreams)
	var metadata map[string]interface{}
	var imageMetadata map[string]interface{}
	var buildLabelsMap map[string]interface{}

	for _, content := range itemsMap {
		metadata = content.(map[string]interface{})["metadata"].(map[string]interface{})
		imageMetadata = content.(map[string]interface{})["image"].(map[string]interface{})["metadata"].(map[string]interface{})
		istagname := metadata["name"].(string)
		isNamespace := metadata["namespace"].(string)
		isLink := metadata["selfLink"].(string)
		isDate := metadata["creationTimestamp"].(string)
		sha := imageMetadata["name"].(string)
		// InfoLogger.Println("image test:" + ImagesMap[sha])
		// os.Exit(1)
		if ImagesMap[sha].(map[string]interface{})["dockerImageMetadata"].(map[string]interface{})["Config"].(map[string]interface{})["Labels"] != nil {
			buildLabelsMap = ImagesMap[sha].(map[string]interface{})["dockerImageMetadata"].(map[string]interface{})["Config"].(map[string]interface{})["Labels"].(map[string]interface{})
		}
		imagestreamfields := strings.Split(istagname, `:`)
		imagestreamName := imagestreamfields[0]
		if !IsNamesForFamily[CmdParams.Family][imagestreamName] {
			continue
		}
		tagName := imagestreamfields[1]
		isAge := fmt.Sprintf("%v", ageInDays(isDate))

		shaNames = addNamesToShaNames(shaNames, sha, isNamespace+"/"+istagname)

		myIstag := T_istag{
			Imagestream: imagestreamName,
			Tagname:     tagName,
			Namespace:   isNamespace,
			Link:        isLink,
			Date:        isDate,
			AgeInDays:   isAge,
			Sha:         sha,
			Build:       setBuildLabels(buildLabelsMap),
		}

		shaStreams = addNamesToShaStreams(shaStreams, imagestreamName, sha, myIstag)

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
		// tmp := resultSha[sha]
		// if err := mergo.Merge(&tmp, mySha); err != nil {
		// 	ErrorLogger.Println("merge mySha to resultSha" + ": failed: " + err.Error())
		// }
		// resultSha[sha] = tmp
		if resultSha[sha] == nil {
			resultSha[sha] = make(map[string]T_sha)
		}
		t := map[string]T_sha{}
		MergoNestedMaps(&t, resultSha[sha], mySha)
		resultSha[sha] = t
		// resultSha[sha] = mySha
	}
	tmp_result := T_result{
		Istag: resultIstag,
		Sha:   resultSha,
		Is:    resultIstream,
	}
	// if err := mergo.Merge(&result, tmp_result); err != nil {
	// 	exitWithError("ERROR: " + "merge result of namespace to result" + ": failed: " + err.Error())
	// }
	t := T_result{}
	MergoNestedMaps(&t, result, tmp_result)
	result = t
	n_istags := len(result.Istag)
	n_shas := len(result.Sha)
	n_is := len(result.Is)
	result.Report = T_resReport{
		AnzNames:    n_istags,
		AnzShas:     n_shas,
		AnzIstreams: n_is,
	}
	// testIfShaHasMultiIstags(shaNames, namespace)
	return result
}

// Generate map of all istags and return json string with the results
func GetAllIstagsForFamilyInCluster() T_result {
	family := CmdParams.Family
	cluster := CmdParams.Cluster
	namespace := CmdParams.Filter.Namespace

	var result = T_result{}
	if namespace == "" {
		for _, ns := range FamilyNamespaces[family] {
			r := OcGetAllIstagsOfNamespace(result, cluster, ns)
			// if err := mergo.Merge(&result, r); err != nil {
			// 	ErrorLogger.Println("merge mySha to resultSha" + ": failed: " + err.Error())
			// }
			t := T_result{}
			MergoNestedMaps(&t, result, r)
			result = t

		}
	} else {
		result = OcGetAllIstagsOfNamespace(result, cluster, namespace)
	}
	return result
}
