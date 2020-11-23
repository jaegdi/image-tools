package ocrequest

import (
	"encoding/json"
)

type T_ImagesMap map[string]interface{}
type T_ImagesMapAllClusters map[string]T_ImagesMap

var ImagesMap T_ImagesMapAllClusters

// InitAllImagesOfCluster reads images from cluster and converts the items array to
// a map image=>item and set the package var ImagesMap to the result.
func GetAllImagesOfCluster(cluster string) T_ImagesMap {
	imagesJson := ocGetCall(cluster, "", "images", "")
	var imagesMap map[string]interface{}
	if err := json.Unmarshal([]byte(imagesJson), &imagesMap); err != nil {
		LogError("unmarschal images." + err.Error())
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

func InitAllImages(c chan string) {
	if Multiproc {
		ImagesMap = goGetExistingImagesInAllClusters()
	} else {
		for _, cluster := range Clusters.Stages {
			r := T_ImagesMapAllClusters{}
			r[cluster] = GetAllImagesOfCluster(cluster)
			MergoNestedMaps(&ImagesMap, r)
		}
	}
	c <- "InitAllImages Done!"
}
