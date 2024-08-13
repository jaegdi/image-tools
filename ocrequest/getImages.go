package ocrequest

import (
	"encoding/json"
)

type T_ImagesMap map[string]interface{}
type T_ImagesMapAllClusters map[T_clName]T_ImagesMap

var ImagesMap T_ImagesMapAllClusters

// InitAllImagesOfCluster reads images from cluster and converts the items array to
// a map image=>item and set the package var ImagesMap to the result.
func GetAllImagesOfCluster(cluster T_clName) T_ImagesMap {
	imagesJson := ocGetCall(cluster, "", "images", "")
	var imagesMap map[string]interface{}
	if err := json.Unmarshal([]byte(imagesJson), &imagesMap); err != nil {
		ErrorLogger.Println("unmarschal images." + err.Error())
	}
	var metadata map[string]interface{}
	result := T_ImagesMap{}

	if imagesMap["items"] != nil {
		for _, content := range imagesMap["items"].([]interface{}) {
			metadata = content.(map[string]interface{})["metadata"].(map[string]interface{})
			image := metadata["name"].(string)
			if result[image] == nil {
				result[image] = map[string]interface{}{}
			}
			result[image] = content
		}
	}
	return result
}

// InitAllImages switches between multithreading and sequencial execution depending on the variable Multiproc
// func InitAllImages(c chan T_ImagesMapAllClusters) {
func InitAllImages(c chan T_ImagesMapAllClusters) {
	if Multiproc {
		ImagesMap = goGetExistingImagesInAllClusters()
	} else {
		for _, cluster := range FamilyNamespaces[CmdParams.Family].Stages {
			r := T_ImagesMapAllClusters{}
			r[cluster] = GetAllImagesOfCluster(cluster)
			MergoNestedMaps(&ImagesMap, r)
		}
	}
	for _, cluster := range FamilyNamespaces[CmdParams.Family].Stages {
		InfoMsg("Number of Images found in", cluster, ":", len(ImagesMap[cluster]))
	}
	c <- ImagesMap
	// c <- "InitAllImages Done!"
}
