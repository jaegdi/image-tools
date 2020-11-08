package ocrequest

import (
	"encoding/json"
	)

type T_ImagesMap map[string]interface{}

var ImagesMap T_ImagesMap

// OcGetAllImagesOfCluster reads images from cluster and converts the items array to
// a map sha=>item
func OcGetAllImagesOfCluster(cluster string) {
	imagesJson := ocGetCall(cluster, "", "images", "")
	var imagesMap map[string]interface{}
	if err := json.Unmarshal([]byte(imagesJson), &imagesMap); err != nil {
		ErrorLogger.Println("unmarschal images." + err.Error())
	}
	var metadata map[string]interface{}
	result := T_ImagesMap{}

	for _, content := range imagesMap["items"].([]interface{}) {
		metadata = content.(map[string]interface{})["metadata"].(map[string]interface{})
		sha := metadata["name"].(string)
		if result[sha] == nil {
			result[sha] = map[string]interface{}{}
		}
		result[sha] = content
	}
	ImagesMap = result
}
