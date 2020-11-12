package ocrequest

import (
	"encoding/json"
)

type T_ImagesMap map[string]interface{}

var ImagesMap T_ImagesMap

// InitAllImagesOfCluster reads images from cluster and converts the items array to
// a map image=>item and set the package var ImagesMap to the result.
func InitAllImagesOfCluster(cluster string) {
	imagesJson := ocGetCall(cluster, "", "images", "")
	var imagesMap map[string]interface{}
	if err := json.Unmarshal([]byte(imagesJson), &imagesMap); err != nil {
		ErrorLogger.Println("unmarschal images." + err.Error())
	}
	var metadata map[string]interface{}
	result := T_ImagesMap{}

	for _, content := range imagesMap["items"].([]interface{}) {
		metadata = content.(map[string]interface{})["metadata"].(map[string]interface{})
		image := metadata["name"].(string)
		if result[image] == nil {
			result[image] = map[string]interface{}{}
		}
		result[image] = content
	}
	ImagesMap = result
}
