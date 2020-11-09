package ocrequest

import (
	"encoding/json"
	"log"
	"strings"
)

func GetAppNamespacesForFamily(cluster string, family string) []string {
	namespacesJson := ocGetCall(cluster, "", "namespaces", "")
	namespacesMap := map[string]interface{}{}
	namespaceList := []string{}
	err := json.Unmarshal([]byte(namespacesJson), &namespacesMap)
	if err != nil {
		ErrorLogger.Println("generate Map for AppNamespaces." + err.Error())
	}
	if len(namespacesMap["items"].([]interface{})) > 0 {
		for _, v := range namespacesMap["items"].([]interface{}) {
			if v.(map[string]interface{})["metadata"].(map[string]interface{})["name"] != nil {
				ns := v.(map[string]interface{})["metadata"].(map[string]interface{})["name"].(string)
				if strings.Contains(ns, "pkp-") {
					namespaceList = append(namespaceList, ns)
				}
			}
		}
	}
	return namespaceList
}

func FilterDcResults(cluster string, ns string, data T_runningObjects) T_usedIstagsResult {
	results := T_usedIstagsResult{}
	// ToDo Implementation
	dc := data.Dc["items"]
	for _, content := range dc.([]interface{}) {
		for _, containers := range content.(map[string]interface{})["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"].([]interface{}) {
			image := containers.(map[string]interface{})["image"].(string)
			imageParts := strings.Split(image, "/")
			istag := imageParts[len(imageParts)-1]
			var is string
			var tag string
			var sha string
			istagParts := strings.Split(istag, "@")
			if len(istagParts) < 2 {
				istagParts = strings.Split(istag, ":")
				sha = ""
			} else {
				sha = istagParts[len(istagParts)-1]
			}
			if len(istagParts) < 2 {
				is = istagParts[len(istagParts)-1]
				tag = ""
			} else {
				is = istagParts[len(istagParts)-2]
				tag = istagParts[len(istagParts)-1]
			}
			usedIstag := T_usedIstag{
				UsedInNamespace: ns,
				Sha:             sha,
				Cluster:         cluster,
			}
			if results[is] == nil {
				results[is] = map[string][]T_usedIstag{}
			}
			if results[is][tag] == nil {
				results[is][tag] = []T_usedIstag{}
			}
			results[is][tag] = append(results[is][tag], usedIstag)
		}
	}
	return results
}

func ocGetAllUsedIstagsOfNamespace(cluster string, namespace string) T_usedIstagsResult {
	istagsDcJson := ocGetCall(cluster, namespace, "deploymentconfigs", "")
	istagsJobJson := ocGetCall(cluster, namespace, "jobs", "")
	istagsCronjobJson := ocGetCall(cluster, namespace, "cronjobs", "")
	istagsPodJson := ocGetCall(cluster, namespace, "pods", "")

	var istagsDcResult T_DcResults
	var istagsJobResult T_JobResults
	var istagsCronjobResult T_CronjobResults
	var istagsPodResult T_Results
	result := T_runningObjects{}

	var err error
	err = json.Unmarshal([]byte(istagsDcJson), &istagsDcResult)
	if err != nil {
		ErrorLogger.Println("Query dc" + err.Error())
	}
	err = json.Unmarshal([]byte(istagsJobJson), &istagsJobResult)
	if err != nil {
		ErrorLogger.Println("Query job" + err.Error())
	}
	err = json.Unmarshal([]byte(istagsCronjobJson), &istagsCronjobResult)
	if err != nil {
		ErrorLogger.Println("Query cronjob" + err.Error())
	}
	err = json.Unmarshal([]byte(istagsPodJson), &istagsPodResult)
	if err != nil {
		ErrorLogger.Println("Query pod" + err.Error())
	}

	result.Dc = istagsDcResult
	result.Job = istagsJobResult
	result.Cronjob = istagsCronjobResult
	result.Pod = istagsPodResult
	filteredUsedIsTags := FilterDcResults(cluster, namespace, result)
	return filteredUsedIsTags
}

func GetUsedIstagsForFamilyInCluster(cluster string) T_usedIstagsResult {
	family := CmdParams.Family
	namespace := CmdParams.Filter.Namespace

	var result T_usedIstagsResult
	if namespace == "" {
		for _, ns := range GetAppNamespacesForFamily(cluster, family) {
			InfoLogger.Println("Get used istags of cluster: " + cluster + " in namespace: " + ns)
			log.Println("Get used istags of cluster: " + cluster + " in namespace: " + ns)
			r := ocGetAllUsedIstagsOfNamespace(cluster, ns)
			// if err := mergo.Merge(&result, r); err != nil {
			// 	ErrorLogger.Println("merge mySha to resultSha" + ": failed: " + err.Error())
			// }
			t := T_usedIstagsResult{}
			MergoNestedMaps(&t, result, r)
			result = t
		}
	} else {
		result = ocGetAllUsedIstagsOfNamespace(cluster, namespace)
	}
	return result
}

func GetUsedIstagsForFamily() T_usedIstagsResult {
	var result T_usedIstagsResult
	for _, cluster := range Clusters["stages"].([]string) {
		InfoLogger.Println("Get used istags in cluster: " + cluster)
		log.Println("Get used istags in cluster: " + cluster)
		resultCluster := GetUsedIstagsForFamilyInCluster(cluster)
		t := T_usedIstagsResult{}
		MergoNestedMaps(&t, result, resultCluster)
		result = t
	}
	return result
}
