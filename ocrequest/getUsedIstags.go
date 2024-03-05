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
func GetAppNamespacesForFamily(cluster T_clName, family T_familyName) []T_nsName {
	namespacesJson := ocGetCall(cluster, "", "namespaces", "")
	namespacesMap := map[string]interface{}{}
	namespaceList := []T_nsName{}
	err := json.Unmarshal([]byte(namespacesJson), &namespacesMap)
	if err != nil {
		ErrorLogger.Println("generate Map for AppNamespaces." + err.Error())
	} else {
		// InfoLogger.Println("CHECK: cluster:"+cluster+" family:"+family+" => map:", namespacesMap)
		if len(namespacesMap["metadata"].(map[string]interface{})) > 0 && len(namespacesMap["items"].([]interface{})) > 0 {
			for _, v := range namespacesMap["items"].([]interface{}) {
				if v.(map[string]interface{})["metadata"].(map[string]interface{})["name"] != nil {
					ns := v.(map[string]interface{})["metadata"].(map[string]interface{})["name"].(string)
					if ns != "" && regexValidNamespace.MatchString(ns) {
						namespaceList = append(namespaceList, T_nsName(ns))
					}
				}
			}
		}
	}
	return namespaceList
}

// GetIsAndTag takes a iamge-string and cut it into the single parts T_shaName, T_istagName, T_isName, T_tagName, T_nsName
func GetIsAndTag(imagestr string) (T_shaName, T_istagName, T_isName, T_tagName, T_nsName) {
	var is T_isName
	var tag T_tagName
	var image T_shaName
	imageParts := strings.Split(imagestr, "/")
	var fromNamespace T_nsName
	if len(imageParts) < 2 {
		// it seems, that the image comes from sattelite
		imageParts = strings.Split(imagestr, "_")
		if len(imageParts) < 2 {
			fromNamespace = ""
		} else {
			fromNamespace = T_nsName("Sattelite_" + T_nsName(imageParts[len(imageParts)-2]))
		}
	} else {
		fromNamespace = T_nsName(imageParts[len(imageParts)-2])
	}
	istag := T_istagName(imageParts[len(imageParts)-1])
	istagParts := strings.Split(istag.str(), "@")
	if len(istagParts) < 2 {
		istagParts = strings.Split(istag.str(), ":")
		image = ""
	} else {
		image = T_shaName(istagParts[len(istagParts)-1])
	}
	if len(istagParts) < 2 {
		is = T_isName(istagParts[len(istagParts)-1])
		tag = T_tagName("")
	} else {
		is = T_isName(istagParts[len(istagParts)-2])
		tag = T_tagName(istagParts[len(istagParts)-1])
	}
	return image, istag, is, tag, fromNamespace
}

// getIstagFromContainer get the image url from each container and extract the imagestream and the istag or image, whichever is defined.
func getIstagFromContainer(cluster T_clName, namespace T_nsName, containers []interface{}, results T_usedIstagsResult) T_usedIstagsResult {
	for _, container := range containers {
		image := T_shaName(container.(map[string]interface{})["image"].(string))
		image, istag, is, tag, fromNamespace := GetIsAndTag(string(image))
		if CmdParams.Filter.Isname != "" && is != CmdParams.Filter.Isname && !CmdParams.FilterReg.Isname.MatchString(string(is)) {
			continue
		}
		if CmdParams.Filter.Tagname != "" && tag != CmdParams.Filter.Tagname && !CmdParams.FilterReg.Tagname.MatchString(string(tag)) {
			continue
		}
		if CmdParams.Filter.Istagname != "" && istag != CmdParams.Filter.Istagname && !CmdParams.FilterReg.Istagname.MatchString(string(istag)) {
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
			results[is] = map[T_tagName][]T_usedIstag{}
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
func FilterIstagsFromRunningObjects(cluster T_clName, namespace T_nsName, data T_runningObjects) T_usedIstagsResult {
	results := T_usedIstagsResult{}
	if !(data.Dc == nil || data.Dc["items"] == nil || data.Dc["items"].([]interface{}) == nil) {
		for _, content := range data.Dc["items"].([]interface{}) {
			containers := content.(map[string]interface{})["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"]
			results = getIstagFromContainer(cluster, namespace, containers.([]interface{}), results)
		}
	}
	if !(data.Deploy == nil || data.Deploy["items"] == nil || data.Deploy["items"].([]interface{}) == nil) {
		for _, content := range data.Deploy["items"].([]interface{}) {
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
func ocGetAllUsedIstagsOfNamespace(cluster T_clName, namespace T_nsName) T_usedIstagsResult {
	istagsDcJson := ocGetCall(cluster, namespace, "deploymentconfigs", "")
	istagsDeployJson := ocGetCall(cluster, namespace, "deployments", "")
	istagsJobJson := ocGetCall(cluster, namespace, "jobs", "")
	istagsCronjobJson := ocGetCall(cluster, namespace, "cronjobs", "")
	istagsPodJson := ocGetCall(cluster, namespace, "pods", "")

	var istagsDcResult T_DcResults
	var istagsDeployResult T_DcResults
	var istagsJobResult T_JobResults
	var istagsCronjobResult T_CronjobResults
	var istagsPodResult T_Results
	result := T_runningObjects{}

	var err error
	err = json.Unmarshal([]byte(istagsDcJson), &istagsDcResult)
	if err != nil {
		ErrorLogger.Println("Query dc" + err.Error())
	}
	err = json.Unmarshal([]byte(istagsDeployJson), &istagsDeployResult)
	if err != nil {
		ErrorLogger.Println("Query deployment" + err.Error())
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
	result.Deploy = istagsDeployResult
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
func GetUsedIstagsForFamilyInCluster(family T_familyName, cluster T_clName) T_usedIstagsResult {
	namespace := CmdParams.Filter.Namespace

	var result T_usedIstagsResult
	if namespace == "" {
		for _, ns := range GetAppNamespacesForFamily(cluster, family) {
			InfoLogger.Println("Get used istags of cluster: ", cluster, "in namespace:", ns)
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
			istag := T_istagName(is.str() + ":" + tag.str())
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
		clusters := FamilyNamespaces[CmdParams.Family].Stages
		for _, cluster := range clusters {
			InfoLogger.Println("Get used istags in cluster:", cluster)
			resultCluster := GetUsedIstagsForFamilyInCluster(CmdParams.Family, cluster)
			MergoNestedMaps(&result, resultCluster)
		}
	}
	c <- result
}
