package ocrequest

import (
	"encoding/json"

	"github.com/stretchr/stew/slice"

	// "log"
	"strings"
)

// GetAppNamespacesForFamily retrieves the namespaces for a given family in a specified cluster.
// This function makes an OpenShift API call to get the list of namespaces in the cluster,
// parses the JSON response, and filters the namespaces based on a regex pattern.
// It returns a list of namespaces that match the criteria.
//
// Parameters:
// - cluster: The cluster name (T_clName) where the namespaces are located.
// - family: The family name (T_familyName) to filter the namespaces.
//
// Returns:
// - []T_nsName: A slice of namespace names (T_nsName) that match the criteria.
func GetAppNamespacesForFamily(cluster T_clName, family T_familyName) []T_nsName {
	// Make an OpenShift API call to get the list of namespaces
	namespacesJson := ocGetCall(cluster, "", "namespaces", "")
	namespacesMap := map[string]interface{}{}
	namespaceList := []T_nsName{}

	// Parse the JSON response into a map
	err := json.Unmarshal([]byte(namespacesJson), &namespacesMap)
	VerifyMsg("namespacesMap:", GetJsonFromMap(namespacesMap))

	if err != nil {
		// Log an error message if JSON unmarshalling fails
		ErrorMsg("generate Map for AppNamespaces." + err.Error())
	} else {
		DebugMsg("CHECK: cluster:", cluster, "family:", family, " => map:", namespacesMap)
		if items, ok := namespacesMap["items"]; ok {
			VerifyMsg("Length of namespacesMap in cluster", cluster, "is:", len(items.([]interface{})))

			// Check if the namespacesMap contains metadata and items
			if len(namespacesMap["metadata"].(map[string]interface{})) > 0 && len(items.([]interface{})) > 0 {
				// Iterate over the items in the namespacesMap
				for _, v := range namespacesMap["items"].([]interface{}) {
					if v.(map[string]interface{})["metadata"].(map[string]interface{})["name"] != nil {
						ns := v.(map[string]interface{})["metadata"].(map[string]interface{})["name"].(string)
						DebugMsg("Iterating over namespaceMap:", ns)

						// add namespace to namespacelist if the namespace name matches the regex pattern
						if len(ns) > 0 && regexValidNamespace.MatchString(ns) {
							DebugMsg("Add Namespace to list:", ns)
							namespaceList = append(namespaceList, T_nsName(ns))
						}
					}
				}
			}
		}
	}
	return namespaceList
}

// GetIsAndTag extracts the SHA, ImageStreamTag, ImageStream name, tag name, and namespace from a given image string.
// This function parses the image string to determine the namespace, ImageStream name, tag name, and SHA.
// It handles different formats of image strings, including those from satellite sources.
//
// Parameters:
// - imagestr: A string representing the image, which may include namespace, ImageStream, tag, and SHA.
//
// Returns:
// - T_shaName: The SHA of the image.
// - T_istagName: The full ImageStreamTag name.
// - T_isName: The name of the ImageStream.
// - T_tagName: The name of the tag.
// - T_nsName: The namespace of the image.
func GetIsAndTag(imagestr string) (T_shaName, T_istagName, T_isName, T_tagName, T_nsName, T_registryUrl) {
	var is T_isName
	var tag T_tagName
	var image T_shaName
	imageParts := strings.Split(imagestr, "/")
	var fromNamespace T_nsName
	var registryURL T_registryUrl

	// Check if the image string contains a
	registryURL = T_registryUrl(imagestr)

	// Check if the image string contains a namespace
	if len(imageParts) < 2 {
		// If not, assume the image comes from a satellite source
		imageParts = strings.Split(imagestr, "_")
		if len(imageParts) < 2 {
			fromNamespace = ""
		} else {
			fromNamespace = T_nsName("Sattelite_" + T_nsName(imageParts[len(imageParts)-2]))
		}
	} else {
		fromNamespace = T_nsName(imageParts[len(imageParts)-2])
	}

	// Extract the ImageStreamTag name
	istag := T_istagName(imageParts[len(imageParts)-1])
	istagParts := strings.Split(istag.str(), "@")

	// Check if the ImageStreamTag name contains a SHA
	if len(istagParts) < 2 {
		istagParts = strings.Split(istag.str(), ":")
		image = ""
	} else {
		image = T_shaName(istagParts[len(istagParts)-1])
	}

	// Extract the ImageStream name and tag name
	if len(istagParts) < 2 {
		is = T_isName(istagParts[len(istagParts)-1])
		tag = T_tagName("")
	} else {
		is = T_isName(istagParts[len(istagParts)-2])
		tag = T_tagName(istagParts[len(istagParts)-1])
	}

	return image, istag, is, tag, fromNamespace, registryURL
}

// getIstagFromContainer extracts ImageStreamTag information from a list of containers and updates the results map.
// This function iterates over the provided containers, extracts the image information, and applies various filters
// to determine if the ImageStreamTag should be included in the results. If the ImageStreamTag passes the filters,
// it is added to the results map.
//
// Parameters:
// - cluster: The cluster name (T_clName) where the containers are located.
// - namespace: The namespace (T_nsName) where the containers are located.
// - containers: A slice of interfaces representing the containers.
// - results: A map (T_usedIstagsResult) to store the used ImageStreamTags.
//
// Returns:
// - T_usedIstagsResult: The updated results map containing the used ImageStreamTags.
func getIstagFromContainer(cluster T_clName, namespace T_nsName, containers []interface{}, results T_usedIstagsResult) T_usedIstagsResult {
	for _, container := range containers {
		// Extract the image SHA from the container
		image := T_shaName(container.(map[string]interface{})["image"].(string))
		// Get the ImageStreamTag details from the image SHA
		image, istag, is, tag, fromNamespace, fromregistry := GetIsAndTag(string(image))

		// Apply filters to determine if the ImageStreamTag should be included in the results
		if CmdParams.Filter.Isname != "" && is != CmdParams.Filter.Isname && !CmdParams.FilterReg.Isname.MatchString(string(is)) {
			continue
		}
		if CmdParams.Filter.Tagname != "" && tag != CmdParams.Filter.Tagname && !CmdParams.FilterReg.Tagname.MatchString(string(tag)) {
			continue
		}
		if CmdParams.Filter.Istagname != "" && istag != CmdParams.Filter.Istagname && !CmdParams.FilterReg.Istagname.MatchString(string(istag)) {
			continue
		}
		if CmdParams.Filter.Istagname != "" && istag != T_istagName(CmdParams.Filter.Namespace) && !CmdParams.FilterReg.Istagname.MatchString(string(istag)) {
			continue
		}
		if CmdParams.Filter.Imagename != "" && image != CmdParams.Filter.Imagename {
			continue
		}

		// Create a T_usedIstag structure with the extracted data
		usedIstag := T_usedIstag{
			UsedInNamespace: namespace,
			FromNamespace:   fromNamespace,
			Image:           image,
			Cluster:         cluster,
			RegistryUrl:     fromregistry,
		}

		results = AddUsedIstag(results, is, tag, usedIstag)
	}
	return results
}

// AddUsedIstag adds a T_usedIstag to the results map if it is not already present.
//
// Parameters:
// - results: The map to store the used image stream tags.
// - is: The image stream name.
// - tag: The tag name.
// - usedIstag: The T_usedIstag structure to add.
func AddUsedIstag(results T_usedIstagsResult, is T_isName, tag T_tagName, usedIstag T_usedIstag) T_usedIstagsResult {
	// Initialize the results map for the ImageStream if necessary
	if results[is] == nil {
		results[is] = map[T_tagName][]T_usedIstag{}
	}
	// Initialize the results map for the tag if necessary
	if results[is][tag] == nil {
		results[is][tag] = []T_usedIstag{}
	}
	// Add the used ImageStreamTag to the results map if it is not already present
	if !slice.Contains(results[is][tag], usedIstag) {
		results[is][tag] = append(results[is][tag], usedIstag)
	}
	return results
}

// FilterIstagsFromRunningObjects filters ImageStreamTags from running objects in a specified cluster and namespace.
// This function processes various types of running objects (DeploymentConfigs, Deployments, Jobs, CronJobs, and Pods)
// to extract ImageStreamTag information from their containers. It iterates over the items in each type of object,
// retrieves the container specifications, and calls getIstagFromContainer to update the results map with the used ImageStreamTags.
//
// Parameters:
// - cluster: The cluster name (T_clName) where the running objects are located.
// - namespace: The namespace (T_nsName) where the running objects are located.
// - data: A structure (T_runningObjects) containing the running objects data.
//
// Returns:
// - T_usedIstagsResult: A map containing the used ImageStreamTags.
func FilterIstagsFromRunningObjects(cluster T_clName, namespace T_nsName, data T_runningObjects) T_usedIstagsResult {
	results := T_usedIstagsResult{}

	// Process DeploymentConfigs
	if !(data.Dc == nil || data.Dc["items"] == nil || data.Dc["items"].([]interface{}) == nil) {
		for _, content := range data.Dc["items"].([]interface{}) {
			containers := content.(map[string]interface{})["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"]
			results = getIstagFromContainer(cluster, namespace, containers.([]interface{}), results)
		}
	}

	// Process Deployments
	if !(data.Deploy == nil || data.Deploy["items"] == nil || data.Deploy["items"].([]interface{}) == nil) {
		for _, content := range data.Deploy["items"].([]interface{}) {
			containers := content.(map[string]interface{})["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"]
			results = getIstagFromContainer(cluster, namespace, containers.([]interface{}), results)
		}
	}

	// Process Jobs
	if !(data.Job == nil || data.Job["items"] == nil || data.Job["items"].([]interface{}) == nil) {
		for _, content := range data.Job["items"].([]interface{}) {
			containers := content.(map[string]interface{})["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"]
			results = getIstagFromContainer(cluster, namespace, containers.([]interface{}), results)
		}
	}

	// Process CronJobs
	if !(data.Cronjob == nil || data.Cronjob["items"] == nil || data.Cronjob["items"].([]interface{}) == nil) {
		for _, content := range data.Cronjob["items"].([]interface{}) {
			jobtemplate := content.(map[string]interface{})["spec"].(map[string]interface{})["jobTemplate"].(map[string]interface{})
			containers := jobtemplate["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"]
			results = getIstagFromContainer(cluster, namespace, containers.([]interface{}), results)
		}
	}

	// Process Pods
	if !(data.Pod == nil || data.Pod["items"] == nil || data.Pod["items"].([]interface{}) == nil) {
		for _, content := range data.Pod["items"].([]interface{}) {
			containers := content.(map[string]interface{})["spec"].(map[string]interface{})["containers"]
			results = getIstagFromContainer(cluster, namespace, containers.([]interface{}), results)
		}
	}

	return results
}

// ocGetAllUsedIstagsOfNamespace retrieves used ImageStreamTags from a specified namespace in a cluster.
// This function queries the OpenShift API to get deploymentconfigs, deployments, jobs, cronjobs, and pods
// in the given namespace. It unmarshals the JSON responses into respective result structures, aggregates
// these results into a T_runningObjects structure, and then filters the ImageStreamTags using the
// FilterIstagsFromRunningObjects function. The final result is a map of used ImageStreamTags.
//
// Parameters:
// - cluster: The cluster name (T_clName) where the namespace is located.
// - namespace: The namespace (T_nsName) from which to retrieve the used ImageStreamTags.
//
// Returns:
// - T_usedIstagsResult: A map containing the used ImageStreamTags.
func ocGetAllUsedIstagsOfNamespace(cluster T_clName, namespace T_nsName) T_usedIstagsResult {
	// Query the OpenShift API to get deploymentconfigs, deployments, jobs, cronjobs, and pods
	istagsDcJson := ocGetCall(cluster, namespace, "deploymentconfigs", "")
	istagsDeployJson := ocGetCall(cluster, namespace, "deployments", "")
	istagsJobJson := ocGetCall(cluster, namespace, "jobs", "")
	istagsCronjobJson := ocGetCall(cluster, namespace, "cronjobs", "")
	istagsPodJson := ocGetCall(cluster, namespace, "pods", "")

	// Initialize result structures for each type of object
	var istagsDcResult T_DcResults
	var istagsDeployResult T_DcResults
	var istagsJobResult T_JobResults
	var istagsCronjobResult T_CronjobResults
	var istagsPodResult T_Results
	result := T_runningObjects{}

	var err error
	// Unmarshal the JSON responses into the respective result structures
	err = json.Unmarshal([]byte(istagsDcJson), &istagsDcResult)
	if err != nil {
		ErrorMsg("Query dc" + err.Error())
	}
	err = json.Unmarshal([]byte(istagsDeployJson), &istagsDeployResult)
	if err != nil {
		ErrorMsg("Query deployment" + err.Error())
	}
	err = json.Unmarshal([]byte(istagsJobJson), &istagsJobResult)
	if err != nil {
		ErrorMsg("Query job" + err.Error())
	}
	err = json.Unmarshal([]byte(istagsCronjobJson), &istagsCronjobResult)
	if err != nil {
		ErrorMsg("Query cronjob" + err.Error())
	}
	err = json.Unmarshal([]byte(istagsPodJson), &istagsPodResult)
	if err != nil {
		ErrorMsg("Query pod" + err.Error())
	}

	// Aggregate the results into a T_runningObjects structure
	result.Dc = istagsDcResult
	result.Deploy = istagsDeployResult
	result.Job = istagsJobResult
	result.Cronjob = istagsCronjobResult
	result.Pod = istagsPodResult

	// Filter the ImageStreamTags from the running objects and return the result
	filteredUsedIsTags := FilterIstagsFromRunningObjects(cluster, namespace, result)
	return filteredUsedIsTags
}

// GetUsedIstagsForFamilyInCluster retrieves used ImageStreamTags from a specified family in a cluster.
// This function searches for deploymentconfigs, jobs, cronjobs, and pods in all namespaces that belong to the family.
// It registers the images of these objects and generates a map with all these ImageStreamTags, returning this map as the result.
//
// Parameters:
// - family: The family name (T_familyName) to filter the namespaces.
// - cluster: The cluster name (T_clName) where the namespaces are located.
//
// Returns:
// - T_usedIstagsResult: A map containing the used ImageStreamTags.
func GetUsedIstagsForFamilyInCluster(family T_familyName, cluster T_clName) T_usedIstagsResult {
	namespace := CmdParams.Filter.Namespace

	var result T_usedIstagsResult
	// If no specific namespace is provided, iterate over all namespaces for the family
	if namespace == "" {
		count := 0
		for _, ns := range GetAppNamespacesForFamily(cluster, family) {
			VerifyMsg("Count:", count, "Get used istags of cluster: ", cluster, "in namespace:", ns)
			count++
			// Retrieve used ImageStreamTags for the current namespace
			r := ocGetAllUsedIstagsOfNamespace(cluster, ns)
			// Merge the results into the final result map
			MergoNestedMaps(&result, r)
		}
	} else {
		// If a specific namespace is provided, retrieve used ImageStreamTags for that namespace
		result = ocGetAllUsedIstagsOfNamespace(cluster, namespace)
	}
	return result
}

// PutShaIntoUsedIstags writes the SHA from allTags into usedIsTags if the SHA is empty.
// This function iterates over the used ImageStreamTags in the results, checks if the SHA is empty,
// and if so, retrieves the SHA from the allTags map and updates the used ImageStreamTags accordingly.
//
// Parameters:
// - c: A channel (chan T_completeResults) to send the updated results.
// - results: A structure (T_completeResults) containing the used and all ImageStreamTags.
func PutShaIntoUsedIstags(c chan T_completeResults, results T_completeResults) {
	count1 := 0
	count2 := 0
	count3 := 0
	for is, isMap := range results.UsedIstags {
		VerifyMsg("Count:", count1, "iterate over IS")
		count1++
		isMap := isMap
		for tag, tagArray := range isMap {
			VerifyMsg("Count:", count2, "iterate over TAG")
			count2++
			tagArray := tagArray
			istag := T_istagName(is.str() + ":" + tag.str())
			for i, tagMap := range tagArray {
				VerifyMsg("Count:", count3, "iterate over TAGMAP")
				count3++
				tagMap := tagMap
				// Check if the SHA is empty and the cluster is not empty
				if tagMap.Image == "" && tagMap.Cluster != "" {
					// Retrieve the SHA from the allTags map and update the used ImageStreamTags
					if results.AllIstags[tagMap.Cluster].Istag[istag][tagMap.FromNamespace].Image != "" {
						results.UsedIstags[is][tag][i].Image = results.AllIstags[tagMap.Cluster].Istag[istag][tagMap.FromNamespace].Image
						results.UsedIstags[is][tag][i].AgeInDays = results.AllIstags[tagMap.Cluster].Istag[istag][tagMap.FromNamespace].AgeInDays
					}
				}
			}
		}
	}
	// Send the updated results through the channel
	c <- results
}

// GetUsedIstagsForFamily retrieves used ImageStreamTags from all clusters for a specified family.
// This function searches for deploymentconfigs, jobs, cronjobs, and pods in all namespaces that belong to the family.
// It registers the images of these objects and generates a map with all these ImageStreamTags, returning this map as the result.
// The function can operate in multiprocess mode, where it concurrently retrieves ImageStreamTags from all clusters,
// or in single-process mode, where it sequentially processes each cluster.
//
// Parameters:
// - c: A channel (chan T_usedIstagsResult) to send the final result map containing the used ImageStreamTags.
func GetUsedIstagsForFamily(c chan T_usedIstagsResult) {
	var result T_usedIstagsResult
	if Multiproc {
		// Multiprocess mode: concurrently retrieve used ImageStreamTags for the family from all clusters
		VerifyMsg("Get used istags for family: ", CmdParams.Family)
		result = goGetUsedIstagsForFamilyInAllClusters(CmdParams.Family)
		VerifyMsg(GetJsonFromMap(result))
	} else {
		// Single-process mode: sequentially retrieve used ImageStreamTags for the family from each cluster
		clusters := FamilyNamespaces[CmdParams.Family].Stages
		for _, cluster := range clusters {
			VerifyMsg("Get used istags in cluster:", cluster)
			resultCluster := GetUsedIstagsForFamilyInCluster(CmdParams.Family, cluster)
			MergoNestedMaps(&result, resultCluster)
		}
	}
	// Send the final result map through the channel
	c <- result
	DebugMsg("UsedImagetags scraped")
}
