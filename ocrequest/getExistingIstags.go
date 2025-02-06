// Package ocrequest provides primitives to query an oc-cluster for istags
// calculate and filter them. And provides a json output primitive.
package ocrequest

import (
	"encoding/json"
	"strings"
)

var IsNamesForFamily T_IsNamesForFamily

// joinShaStreams extracts and returns the keys from a map of ImageStream tags.
// This function takes a map where the keys are ImageStream tag names and the values
// are lists of tags. It iterates over the map, collects all the keys, and returns them
// as a slice.
//
// Parameters:
// - mymap: A map where the keys are ImageStream tag names and the values are lists of tags.
//
// Returns:
// - A slice of ImageStream tag names extracted from the map.
//
// Example:
//
//	mymap := T_Istags_List{
//	    "tag1": {"sha1", "sha2"},
//	    "tag2": {"sha3", "sha4"},
//	}
//	result := joinShaStreams(mymap)
//	>// result will be []T_istagName{"tag1", "tag2"}
func joinShaStreams(mymap T_Istags_List) []T_istagName {
	keys := []T_istagName{}
	for k := range mymap {
		keys = append(keys, k)
	}
	return keys
}

// appendJoinedNamesToImagestreams appends joined names to the specified ImageStream.
// This function takes an ImageStream map, an ImageStream name, an image name, and a list of joined names.
// It ensures that the ImageStream and image entries exist in the map, and then appends each joined name
// to the image entry.
//
// Parameters:
// - istream: The ImageStream map to be updated.
// - imagestreamName: The name of the ImageStream to be updated.
// - image: The name of the image to which the joined names will be appended.
// - joinedNames: A list of joined names to be appended to the image entry.
//
// Returns:
// - The updated ImageStream map.
//
// Example:
//
//	istream := T_resIs{}
//	imagestreamName := T_isName("example-stream")
//	image := T_shaName("example-image")
//	joinedNames := []T_istagName{"tag1", "tag2", "tag3"}
//	result := appendJoinedNamesToImagestreams(istream, imagestreamName, image, joinedNames)
//	>// result will be updated to include the joined names for the specified image in the ImageStream.
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

// InitIsNamesForFamily initializes the names of ImageStreams for a given family.
// This function takes a family name and a list of ImageStream names, and initializes
// the ImageStream names for that family. It returns a map where the keys are the family
// names and the values are the corresponding ImageStream names.
//
// Parameters:
// - family: The name of the family for which the ImageStream names are to be initialized.
// - isNames: A list of ImageStream names to be associated with the family.
//
// Returns:
// - A map where the keys are the family names and the values are the corresponding ImageStream names.
//
// Example:
//
//	family := "example-family"
//	isNames := []string{"is1", "is2", "is3"}
//	result := InitIsNamesForFamily(family, isNames)
//	>// result will be map[string][]string{"example-family": {"is1", "is2", "is3"}}
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
				ErrorMsg("Unmarshal imagestreams. " + err.Error())
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

// setBuildLabels adds or updates labels for a build object.
// This function takes a build object and a map of labels, and sets the labels
// in the build object accordingly. If a label already exists, it is updated;
// otherwise, it is added.
//
// Parameters:
// - build: The build object to be updated.
// - labels: A map of labels to be added or updated in the build object.
//
// Example:
//
//	build := &Build{}
//	labels := map[string]string{"env": "production", "version": "1.0"}
//	setBuildLabels(build, labels)
func setBuildLabels(buildLabelsMap map[string]interface{}) T_istagBuildLabels {
	buildLabels := T_istagBuildLabels{}
	buildLabelsJSON := []byte(GetJsonFromMap(buildLabelsMap))
	if err := json.Unmarshal(buildLabelsJSON, &buildLabels); err != nil {
		ErrorMsg("Unmarshal unescaped String", err)
	}
	return buildLabels
}

// OcGetAllIstagsOfNamespace retrieves all ImageStreamTags of a specified namespace from a given cluster.
// This function makes an OpenShift API call to get the ImageStreamTags, parses the JSON response,
// and processes the tags to filter and organize them into the result structure.
//
// Parameters:
// - result: The initial result structure to be updated with the retrieved ImageStreamTags.
// - cluster: The name of the cluster from which to retrieve the ImageStreamTags.
// - namespace: The name of the namespace from which to retrieve the ImageStreamTags.
//
// Returns:
// - An updated result structure containing the filtered and organized ImageStreamTags.
//
// Example:
//
//	result := T_result{}
//	cluster := T_clName("example-cluster")
//	namespace := T_nsName("example-namespace")
//	updatedResult := OcGetAllIstagsOfNamespace(result, cluster, namespace)
//	// updatedResult will contain the ImageStreamTags from the specified namespace and cluster.
func OcGetAllIstagsOfNamespace(result T_result, cluster T_clName, namespace T_nsName) T_result {
	// Read the imageStreamTags by making an OpenShift API call
	istagsJson := ocGetCall(cluster, namespace, "imagestreamtags", "")
	DebugMsg("istagJson:", istagsJson)

	// Parse the JSON response into a map
	istagsMap, err := parseIstagsJson(istagsJson)
	if err != nil {
		// Log error messages if JSON parsing fails
		ErrorMsg("Error: unmarshal imagestreamtags.", istagsJson)
		ErrorMsg("Error: unmarshal imagestreamtags.", err.Error())
		return T_result{}
	}

	// Initialize result maps for ImageStreamTags, images, and ImageStreams
	resultIstag := make(T_resIstag)
	resultSha := make(T_resSha)
	resultIstream := make(T_resIs)

	var itemsMap []interface{}
	// Check if the "items" key exists in the parsed JSON map
	if istagsMap["items"] != nil {
		itemsMap = istagsMap["items"].([]interface{})
	}

	// Initialize maps for SHA names and SHA streams
	shaNames := make(T_shaNames)
	shaStreams := make(T_shaStreams)

	// Process each item in the itemsMap
	for _, content := range itemsMap {
		processItem(content, cluster, &resultIstag, &resultSha, &resultIstream, shaNames, shaStreams)
	}

	// Create a temporary result structure with the processed data
	tmp_result := T_result{
		Istag: resultIstag,
		Image: resultSha,
		Is:    resultIstream,
	}

	// Merge the temporary result into the initial result structure
	MergoNestedMaps(&result, tmp_result)

	// Update the report with the counts of ImageStreamTags, images, and ImageStreams
	n_istags := len(result.Istag)
	n_shas := len(result.Image)
	n_is := len(result.Is)
	result.Report = T_resReport{
		Anz_ImageStreamTags: n_istags,
		Anz_Images:          n_shas,
		Anz_ImageStreams:    n_is,
	}

	// Return the updated result structure
	return result
}

// parseIstagsJson parses the given JSON string into a map.
// This function takes a JSON string representing ImageStreamTags and unmarshals it into a map.
// It returns the parsed map and any error encountered during unmarshalling.
//
// Parameters:
// - istagsJson: A string containing the JSON representation of ImageStreamTags.
//
// Returns:
// - map[string]interface{}: A map containing the parsed JSON data.
// - error: An error object if any error occurred during JSON unmarshalling.
func parseIstagsJson(istagsJson string) (map[string]interface{}, error) {
	var istagsMap map[string]interface{}
	err := json.Unmarshal([]byte(istagsJson), &istagsMap)
	return istagsMap, err
}

// processItem processes a single item from the ImageStreamTags list.
// This function extracts metadata and image metadata from the given content,
// checks if the item should be processed based on filters, and updates the result maps accordingly.
//
// Parameters:
// - content: An interface{} representing a single item from the ImageStreamTags list.
// - cluster: The cluster name (T_clName) where the item is located.
// - resultIstag: A pointer to the result map (T_resIstag) for ImageStreamTags.
// - resultSha: A pointer to the result map (T_resSha) for image SHAs.
// - resultIstream: A pointer to the result map (T_resIs) for ImageStreams.
// - shaNames: A map (T_shaNames) to store SHA names.
// - shaStreams: A map (T_shaStreams) to store SHA streams.
func processItem(content interface{}, cluster T_clName, resultIstag *T_resIstag, resultSha *T_resSha, resultIstream *T_resIs, shaNames T_shaNames, shaStreams T_shaStreams) {
	metadata := content.(map[string]interface{})["metadata"].(map[string]interface{})
	imageMetadata := content.(map[string]interface{})["image"].(map[string]interface{})["metadata"].(map[string]interface{})
	istagname := T_istagName(metadata["name"].(string))
	isNamespace := T_nsName(metadata["namespace"].(string))
	isDate := metadata["creationTimestamp"].(string)
	sha := T_shaName(imageMetadata["name"].(string))

	if !shouldProcessItem(istagname, sha) {
		return
	}

	buildLabelsMap := getBuildLabelsMap(cluster, istagname, sha)
	imagestreamName, tagName := getImageStreamAndTagName(istagname)
	isAge := ageInDays(isDate)

	if !matchIsIstagToFilterParams(imagestreamName, tagName, istagname, isNamespace, isAge) {
		return
	}

	shaNames.Add(sha, T_istagName(isNamespace.str()+"/"+istagname.str()))

	myIstag := T_istag{
		Imagestream: imagestreamName,
		Tagname:     tagName,
		Namespace:   isNamespace,
		Date:        isDate,
		AgeInDays:   isAge,
		Image:       sha,
		Build:       buildLabelsMap,
	}

	shaStreams.Add(imagestreamName, sha, myIstag)

	updateResultMaps(resultIstag, resultSha, resultIstream, istagname, isNamespace, myIstag, sha, shaNames)
}

// shouldProcessItem determines if an item should be processed based on filter criteria.
// This function checks if the given ImageStreamTag name and SHA match the filter parameters
// specified in CmdParams. If the filters are not empty and the item does not match the filters,
// the function returns false, indicating that the item should not be processed.
//
// Parameters:
// - istagname: The name of the ImageStreamTag (T_istagName).
// - sha: The SHA of the image (T_shaName).
//
// Returns:
// - bool: True if the item should be processed, false otherwise.
func shouldProcessItem(istagname T_istagName, sha T_shaName) bool {
	if CmdParams.Filter.Imagename != "" && sha != CmdParams.Filter.Imagename {
		return false
	}
	if CmdParams.Filter.Istagname != "" && istagname != CmdParams.Filter.Istagname && !CmdParams.FilterReg.Istagname.MatchString(string(istagname)) {
		return false
	}
	return true
}

// getBuildLabelsMap retrieves the build labels map for a given image SHA in a specified cluster.
// This function checks if the ImagesMap contains metadata for the given cluster and SHA.
// If the metadata includes build labels, it sets these labels in the buildLabelsMap structure.
//
// Parameters:
// - cluster: The cluster name (T_clName) where the image is located.
// - sha: The SHA of the image (T_shaName).
//
// Returns:
// - T_istagBuildLabels: A structure containing the build labels for the specified image.
func getBuildLabelsMap(cluster T_clName, istagname T_istagName, sha T_shaName) T_istagBuildLabels {
	buildLabelsMap := T_istagBuildLabels{}
	if CmdParams.Options.Debug {
		DebugMsg("IsTag: "+istagname, "ImagesMap: ", ImagesMap)
	}
	imageMetadata, ok := ImagesMap[cluster][sha.str()].(map[string]interface{})["dockerImageMetadata"]
	if !ok || imageMetadata == nil {
		return buildLabelsMap
	}
	config, ok := imageMetadata.(map[string]interface{})["Config"]
	if !ok || config == nil {
		return buildLabelsMap
	}
	labels, ok := config.(map[string]interface{})["Labels"]
	if !ok || labels == nil {
		return buildLabelsMap
	}
	buildLabelsMap.Set(labels.(map[string]interface{}))
	return buildLabelsMap
}

// getImageStreamAndTagName extracts the ImageStream name and tag name from a given ImageStreamTag name.
// This function splits the ImageStreamTag name (istagname) into two parts: the ImageStream name and the tag name.
// It assumes that the ImageStreamTag name is in the format "ImageStream:Tag".
//
// Parameters:
// - istagname: The name of the ImageStreamTag (T_istagName).
//
// Returns:
// - T_isName: The name of the ImageStream.
// - T_tagName: The name of the tag.
func getImageStreamAndTagName(istagname T_istagName) (T_isName, T_tagName) {
	imagestreamfields := strings.Split(istagname.str(), `:`)
	imagestreamName := T_isName(imagestreamfields[0])
	tagName := T_tagName(imagestreamfields[1])
	return imagestreamName, tagName
}

// updateResultMaps updates the result maps with the processed ImageStreamTag data.
// This function adds the given ImageStreamTag (myIstag) to the result maps for ImageStreamTags (resultIstag),
// image SHAs (resultSha), and ImageStreams (resultIstream). It ensures that the maps are properly initialized
// and merges the new data with existing entries.
//
// Parameters:
// - resultIstag: A pointer to the result map (T_resIstag) for ImageStreamTags.
// - resultSha: A pointer to the result map (T_resSha) for image SHAs.
// - resultIstream: A pointer to the result map (T_resIs) for ImageStreams.
// - istagname: The name of the ImageStreamTag (T_istagName).
// - isNamespace: The namespace of the ImageStreamTag (T_nsName).
// - myIstag: The ImageStreamTag structure (T_istag) containing the processed data.
// - sha: The SHA of the image (T_shaName).
// - shaNames: A map (T_shaNames) to store SHA names.
func updateResultMaps(resultIstag *T_resIstag, resultSha *T_resSha, resultIstream *T_resIs, istagname T_istagName, isNamespace T_nsName, myIstag T_istag, sha T_shaName, shaNames T_shaNames) {
	mySha := map[T_istagName]T_sha{
		istagname: {
			Istags:      shaNames[sha],
			Imagestream: myIstag.Imagestream,
			Namespace:   myIstag.Namespace,
			Date:        myIstag.Date,
			AgeInDays:   myIstag.AgeInDays,
		},
	}

	joinedNames := joinShaStreams(shaNames[sha])
	*resultIstream = appendJoinedNamesToImagestreams(*resultIstream, myIstag.Imagestream, sha, joinedNames)
	if (*resultIstag)[istagname] == nil {
		(*resultIstag)[istagname] = map[T_nsName]T_istag{}
	}
	(*resultIstag)[istagname][isNamespace] = myIstag

	if (*resultSha)[sha] == nil {
		(*resultSha)[sha] = make(map[T_istagName]T_sha)
	}
	t := (*resultSha)[sha]
	MergoNestedMaps(&t, mySha)
	(*resultSha)[sha] = t
}

// GetAllIstagsForFamily retrieves all ImageStreamTags for a specified family across all clusters.
// This function checks if multiprocessing is enabled and either calls a function to get the tags
// concurrently or iterates through each cluster and namespace to retrieve the tags sequentially.
// The results are merged into a single result structure.
//
// Parameters:
// - c: A channel to send the result structure containing all ImageStreamTags.
//
// Example:
//
//	resultChannel := make(chan T_ResultExistingIstagsOverAllClusters)
//	go GetAllIstagsForFamily(resultChannel)
//	result := <-resultChannel
func GetAllIstagsForFamily(c chan T_ResultExistingIstagsOverAllClusters) {
	// Retrieve the family and namespace filter from command parameters
	family := CmdParams.Family
	namespace := CmdParams.Filter.Namespace
	// Initialize the result map for existing image stream tags
	var result = T_ResultExistingIstagsOverAllClusters{}

	// Check if multiprocessing is enabled
	if Multiproc {
		// Retrieve existing image stream tags for the family across all clusters using multiprocessing
		result = goGetExistingIstagsForFamilyInAllClusters(family)
	} else {
		// Iterate over each cluster stage for the specified family
		for _, cluster := range FamilyNamespaces[CmdParams.Family].Stages {
			// Check if a namespace filter is applied
			if namespace == "" {
				// Process all image namespaces of the family for the cluster
				for _, ns := range FamilyNamespaces[family].ImageNamespaces[cluster] {
					// Retrieve existing image stream tags for the namespace and merge into the result map
					r := T_ResultExistingIstagsOverAllClusters{cluster: OcGetAllIstagsOfNamespace(result[cluster], cluster, ns)}
					MergoNestedMaps(&result, r)
				}
			} else {
				// Retrieve existing image stream tags for the specific namespace and merge into the result map
				result = T_ResultExistingIstagsOverAllClusters{cluster: OcGetAllIstagsOfNamespace(result[cluster], cluster, namespace)}
			}
		}
	}
	// Send the result through the provided channel
	c <- result
}
