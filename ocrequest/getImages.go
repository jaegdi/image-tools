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
	// Führt einen oc-Befehl aus, um die Images im angegebenen Cluster abzurufen
	imagesJson := ocGetCall(cluster, "", "images", "")
	var imagesMap map[string]interface{}

	// Unmarshalt die JSON-Antwort in eine Map
	if err := json.Unmarshal([]byte(imagesJson), &imagesMap); err != nil {
		ErrorMsg("unmarschal images." + err.Error())
	}

	var metadata map[string]interface{}
	result := T_ImagesMap{}

	// Überprüft, ob die "items"-Schlüssel in der Map vorhanden ist
	if imagesMap["items"] != nil {
		// Iteriert über die Items und extrahiert die Metadaten
		for _, content := range imagesMap["items"].([]interface{}) {
			metadata = content.(map[string]interface{})["metadata"].(map[string]interface{})
			image := metadata["name"].(string)

			// Initialisiert die Map für das Image, falls noch nicht vorhanden
			if result[image] == nil {
				result[image] = map[string]interface{}{}
			}

			// Setzt das Image in das Ergebnis
			result[image] = content
		}
	}

	// Gibt die resultierende Map zurück
	return result
}

// InitAllImages switches between multithreading and sequential execution depending on the variable Multiproc.
// It initializes the ImagesMap with images from all clusters and sends the result through the provided channel.
func InitAllImages(c chan T_ImagesMapAllClusters) {
	if Multiproc {
		// Wenn Multiproc aktiviert ist, rufe die Funktion auf, die Images in allen Clustern parallel abruft
		ImagesMap = goGetExistingImagesInAllClusters()
	} else {
		// Andernfalls iteriere über alle Cluster in den Stages der angegebenen Familie
		for _, cluster := range FamilyNamespaces[CmdParams.Family].Stages {
			// Initialisiere eine neue Map für die Images aller Cluster
			r := T_ImagesMapAllClusters{}

			// Rufe die Images für den aktuellen Cluster ab und speichere sie in der Map
			r[cluster] = GetAllImagesOfCluster(cluster)

			// Füge die abgerufenen Images zur globalen ImagesMap hinzu
			MergoNestedMaps(&ImagesMap, r)
		}
	}

	// Iteriere über alle Cluster in den Stages der angegebenen Familie
	for _, cluster := range FamilyNamespaces[CmdParams.Family].Stages {
		// Protokolliere die Anzahl der in jedem Cluster gefundenen Images
		VerifyMsg("Number of Images found in", cluster, ":", len(ImagesMap[cluster]))
	}

	// Sende die resultierende ImagesMap durch den bereitgestellten Kanal
	c <- ImagesMap
}
