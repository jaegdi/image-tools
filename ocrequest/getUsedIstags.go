package ocrequest

import (
	"encoding/json"

	"github.com/stretchr/stew/slice"

	// "log"
	"strings"
)

// GetAppNamespacesForFamily gets all namespaces of a cluster and
// filters them with the pattern '^<family>-' to find all namespaces, which
// names starting with the family-name followed by a dash. It returns a slice
// list with the application namespaces for the family.
func GetAppNamespacesForFamily(cluster string, family string) []string {
	namespacesJson := ocGetCall(cluster, "", "namespaces", "")
	namespacesMap := map[string]interface{}{}
	namespaceList := []string{}
	err := json.Unmarshal([]byte(namespacesJson), &namespacesMap)
	if err != nil {
		LogError("generate Map for AppNamespaces." + err.Error())
	} else {
		// LogMsg("CHECK: cluster:"+cluster+" family:"+family+" => map:", namespacesMap)
		if len(namespacesMap["metadata"].(map[string]interface{})) > 0 && len(namespacesMap["items"].([]interface{})) > 0 {
			for _, v := range namespacesMap["items"].([]interface{}) {
				if v.(map[string]interface{})["metadata"].(map[string]interface{})["name"] != nil {
					ns := v.(map[string]interface{})["metadata"].(map[string]interface{})["name"].(string)
					if ns != "" && regexValidNamespace.MatchString(ns) {
						namespaceList = append(namespaceList, ns)
					}
				}
			}
		}
	}
	return namespaceList
}

// getIstagFromContainer get the image url from each container and extract the imagestream and the istag or image, whichever is defined.
func getIstagFromContainer(cluster string, namespace string, containers []interface{}, results T_usedIstagsResult) T_usedIstagsResult {
	var is string
	var tag string
	// var image string
	for _, container := range containers {
		image := container.(map[string]interface{})["image"].(string)
		imageParts := strings.Split(image, "/")
		fromNamespace := imageParts[len(imageParts)-2]
		istag := imageParts[len(imageParts)-1]
		istagParts := strings.Split(istag, "@")
		if len(istagParts) < 2 {
			istagParts = strings.Split(istag, ":")
			image = ""
		} else {
			image = istagParts[len(istagParts)-1]
		}
		if len(istagParts) < 2 {
			is = istagParts[len(istagParts)-1]
			tag = ""
		} else {
			is = istagParts[len(istagParts)-2]
			tag = istagParts[len(istagParts)-1]
		}
		if CmdParams.Filter.Isname != "" && is != CmdParams.Filter.Isname {
			continue
		}
		if CmdParams.Filter.Tagname != "" && tag != CmdParams.Filter.Tagname {
			continue
		}
		if CmdParams.Filter.Istagname != "" && istag != CmdParams.Filter.Istagname {
			continue
		}
		if CmdParams.Filter.Imagename != "" && image != CmdParams.Filter.Imagename {
			continue
		}
		usedIstag := T_usedIstag{
			UsedInNamespace: namespace,
			FromNamespace:   fromNamespace,
			Image:           image,
			Cluster:         cluster,
		}
		if results[is] == nil {
			results[is] = map[string][]T_usedIstag{}
		}
		if results[is][tag] == nil {
			results[is][tag] = []T_usedIstag{}
		}
		if !slice.Contains(results[is][tag], usedIstag) {
			results[is][tag] = append(results[is][tag], usedIstag)
		}
	}
	return results
}

// FilterIstagsFromRunningObjects get the right position in the data map, where the containers are defined
// and calls getIsTagFromContainer with this data node to get the istag of this container.
func FilterIstagsFromRunningObjects(cluster string, namespace string, data T_runningObjects) T_usedIstagsResult {
	results := T_usedIstagsResult{}
	if !(data.Dc == nil || data.Dc["items"] == nil || data.Dc["items"].([]interface{}) == nil) {
		for _, content := range data.Dc["items"].([]interface{}) {
			containers := content.(map[string]interface{})["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"]
			results = getIstagFromContainer(cluster, namespace, containers.([]interface{}), results)
		}
	}
	if !(data.Job == nil || data.Job["items"] == nil || data.Job["items"].([]interface{}) == nil) {
		for _, content := range data.Job["items"].([]interface{}) {
			containers := content.(map[string]interface{})["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"]
			results = getIstagFromContainer(cluster, namespace, containers.([]interface{}), results)
		}
	}
	if !(data.Cronjob == nil || data.Cronjob["items"] == nil || data.Cronjob["items"].([]interface{}) == nil) {
		for _, content := range data.Cronjob["items"].([]interface{}) {
			jobtemplate := content.(map[string]interface{})["spec"].(map[string]interface{})["jobTemplate"].(map[string]interface{})
			containers := jobtemplate["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"]
			results = getIstagFromContainer(cluster, namespace, containers.([]interface{}), results)
		}
	}
	if !(data.Pod == nil || data.Pod["items"] == nil || data.Pod["items"].([]interface{}) == nil) {
		for _, content := range data.Pod["items"].([]interface{}) {
			containers := content.(map[string]interface{})["spec"].(map[string]interface{})["containers"]
			results = getIstagFromContainer(cluster, namespace, containers.([]interface{}), results)
		}
	}
	return results
}

// ocGetAllUsedIstagsOfNamespace get __used__ istags from __a namespace__
// of __one__ cluster.
// It looks for deploymentconfigs, jobs, cronjobs and pods in all
// namespaces, that belongs to the family,
// registers the images of these objects and generates
// a map with all these istags and return this map as result.
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
		LogError("Query dc" + err.Error())
	}
	err = json.Unmarshal([]byte(istagsJobJson), &istagsJobResult)
	if err != nil {
		LogError("Query job" + err.Error())
	}
	err = json.Unmarshal([]byte(istagsCronjobJson), &istagsCronjobResult)
	if err != nil {
		LogError("Query cronjob" + err.Error())
	}
	err = json.Unmarshal([]byte(istagsPodJson), &istagsPodResult)
	if err != nil {
		LogError("Query pod" + err.Error())
	}

	result.Dc = istagsDcResult
	result.Job = istagsJobResult
	result.Cronjob = istagsCronjobResult
	result.Pod = istagsPodResult
	filteredUsedIsTags := FilterIstagsFromRunningObjects(cluster, namespace, result)
	return filteredUsedIsTags
}

// GetUsedIstagsForFamilyInCluster get __used__ istags from __one__ cluster.
// It looks for deploymentconfigs, jobs, cronjobs and pods in all
// namespaces, that belongs to the family,
// registers the images of these objects and generates
// a map with all these istags and return this map as result.
func GetUsedIstagsForFamilyInCluster(family string, cluster string) T_usedIstagsResult {
	namespace := CmdParams.Filter.Namespace

	var result T_usedIstagsResult
	if namespace == "" {
		for _, ns := range GetAppNamespacesForFamily(cluster, family) {
			LogMsg("Get used istags of cluster: ", cluster, "in namespace:", ns)
			r := ocGetAllUsedIstagsOfNamespace(cluster, ns)
			MergoNestedMaps(&result, r)
		}
	} else {
		result = ocGetAllUsedIstagsOfNamespace(cluster, namespace)
	}
	return result
}

// PutShaIntoUsedIstags writes sha from allTags into usedIsTags if there the sha is empty
func PutShaIntoUsedIstags(c chan T_completeResults, results T_completeResults) {
	for is, isMap := range results.UsedIstags {
		isMap := isMap
		for tag, tagArray := range isMap {
			tagArray := tagArray
			istag := is + ":" + tag
			for i, tagMap := range tagArray {
				tagMap := tagMap
				if tagMap.Image == "" && tagMap.Cluster != "" {
					if results.AllIstags[tagMap.Cluster].Istag[istag][tagMap.FromNamespace].Image != "" {
						results.UsedIstags[is][tag][i].Image = results.AllIstags[tagMap.Cluster].Istag[istag][tagMap.FromNamespace].Image
						results.UsedIstags[is][tag][i].AgeInDays = results.AllIstags[tagMap.Cluster].Istag[istag][tagMap.FromNamespace].AgeInDays
					}
				}
			}
		}
	}
	c <- results
}

// GetUsedIstagsForFamily get __used__ istags from __all__ clusters.
// It looks for deploymentconfigs, jobs, cronjobs and pods in all
// namespaces, that belongs to the family,
// registers the images of these objects and generates
// a map with all these istags and return this map as result.
func GetUsedIstagsForFamily(c chan T_usedIstagsResult) {
	var result T_usedIstagsResult
	if Multiproc {
		result = goGetUsedIstagsForFamilyInAllClusters(CmdParams.Family)
	} else {
		clusters := Clusters.Stages
		for _, cluster := range clusters {
			LogMsg("Get used istags in cluster:", cluster)
			resultCluster := GetUsedIstagsForFamilyInCluster(CmdParams.Family, cluster)
			MergoNestedMaps(&result, resultCluster)
		}
	}
	c <- result
}
