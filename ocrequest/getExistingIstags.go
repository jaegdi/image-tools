// Package ocrequest provides primitives to query an oc-cluster for istags
// calculate and filter them. And provides a json output primitive.
package ocrequest

import (
	"encoding/json"
	"fmt"
	"github.com/imdario/mergo"
	"log"
	"strings"
)

type T_shaStreams map[string]map[string]map[string]interface{}
type T_shaNames map[string]map[string]interface{}
type T_resIstag map[string]map[string]interface{}

// type T_resSha map[string]map[string]string
type T_resIs map[string]map[string][]string
type T_istag map[string]interface{}
type T_sha map[string]interface{}
type T_isShaTagnames []string
type T_is map[string][]string
type T_resSha map[string]map[string]map[string]interface{}
type T_result map[string]interface{}

// joinShaNames join keys of map to a comma sepearted string.
func joinShaNames(mymap map[string]bool) []string {
	keys := []string{}
	for k := range mymap {
		keys = append(keys, k)
	}
	return keys
}

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
	if a[is][sha] == nil {
		a[is][sha] = T_sha{}
	}
	for k, v := range istag {
		a[is][sha][k] = v
	}
	return a
}

// testIfShaHasMultiIstags test if a sha is tagged with more than one istag and prints them out.
//This function is for logging purposes.
func testIfShaHasMultiIstags(mymap T_shaNames, ns string) {
	for k := range mymap {
		s := "Ns " + ns + ": More than one tag on sha in "
		if len(mymap[k]) > 1 {
			s = s + k + " ==>\n"
			for n := range mymap[k] {
				s = s + "  " + n + "\n"
			}
			log.Println(s)
		}
	}
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
		istream[imagestreamName][sha] =
			append(istream[imagestreamName][sha], v)
	}
	return istream
}

func getIsNamesForFamily() {

}

func OcGetAllIstagsOfNamespace(result map[string]interface{}, cluster string, token string, namespace string) T_result {
	istagsJson := ocAPiCall(cluster, token, namespace, "imagestreamtags", "")

	var istagsresult map[string]interface{}
	err := json.Unmarshal([]byte(istagsJson), &istagsresult)
	if err != nil {
		log.Println("ERROR:", err.Error())
	}

	// result := make(map[string]interface{})
	resultIstag := make(T_resIstag)
	resultSha := make(T_resSha)
	resultIstream := make(T_resIs)

	items := istagsresult["items"].([]interface{})
	shaNames := make(T_shaNames)
	shaStreams := make(T_shaStreams)

	for _, content := range items {
		istagname := content.(map[string]interface{})["metadata"].(map[string]interface{})["name"].(string)
		imagestreamfields := strings.Split(istagname, `:`)
		imagestreamName := imagestreamfields[0]
		tagName := imagestreamfields[1]
		isNamespace := content.(map[string]interface{})["metadata"].(map[string]interface{})["namespace"].(string)
		isLink := content.(map[string]interface{})["metadata"].(map[string]interface{})["selfLink"].(string)
		isDate := content.(map[string]interface{})["metadata"].(map[string]interface{})["creationTimestamp"].(string)
		isAge := fmt.Sprintf("%v", ageInDays(isDate))
		sha := content.(map[string]interface{})["image"].(map[string]interface{})["metadata"].(map[string]interface{})["name"].(string)

		shaNames = addNamesToShaNames(shaNames, sha, isNamespace+"/"+istagname)

		myIstag := T_istag{
			"imagestream": imagestreamName,
			"tagname":     tagName,
			"namespace":   isNamespace,
			"link":        isLink,
			"date":        isDate,
			"ageInDays":   isAge,
			"sha":         sha,
		}

		shaStreams = addNamesToShaStreams(shaStreams, imagestreamName, sha, myIstag)

		mySha := map[string]map[string]interface{}{
			istagname: T_sha{
				"istags":      shaNames[sha],
				"imagestream": imagestreamName,
				"namespace":   isNamespace,
				"link":        isLink,
				"date":        isDate,
				"ageInDays":   isAge,
			},
		}

		joinedNames := joinShaStreams(shaNames[sha])
		resultIstream = appendJoinedNamesToImagestreams(resultIstream, imagestreamName, sha, joinedNames)
		resultIstag[istagname] = myIstag
		tmp := resultSha[sha]
		if err := mergo.Merge(&tmp, mySha); err != nil {
			log.Println("ERROR: " + "merge mySha to resultSha" + ": failed: " + err.Error())
		}
		resultSha[sha] = tmp
		// resultSha[sha] = mySha
	}
	tmp_result := map[string]interface{}{
		"istag": resultIstag,
		"sha":   resultSha,
		"is":    resultIstream,
	}
	if err := mergo.Merge(&result, tmp_result); err != nil {
		// log.Println("ERROR: " + "merge result of namespace to result" + ": failed: " + err.Error())
		exitWithError("ERROR: " + "merge result of namespace to result" + ": failed: " + err.Error())
	}
	n_istags := len(result["istag"].(T_resIstag))
	n_shas := len(result["sha"].(T_resSha))
	n_is := len(result["is"].(T_resIs))
	result["report"] = map[string]int{
		"anz-names":         n_istags,
		"anz-shas":          n_shas,
		"anz-resultIstream": n_is,
	}
	// testIfShaHasMultiIstags(shaNames, namespace)
	return result
}
