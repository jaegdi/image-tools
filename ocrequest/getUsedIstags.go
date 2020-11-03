package ocrequest

import (
	"encoding/json"
	// "fmt"
	// "github.com/imdario/mergo"
	"log"
)

var cluster string = "cid-apc0"
var namespace string = "pkp-unicorn-drei"

func ocGetAllUsedIstagsOfNamespace(cluster string, token string, namespace string) map[string]map[string]interface{} {
	istagsDcJson := ocAPiCall(cluster, token, namespace, "deploymentconfigs", "")
	istagsJobJson := ocAPiCall(cluster, token, namespace, "jobs", "")
	istagsCronjobJson := ocAPiCall(cluster, token, namespace, "cronjobs", "")
	istagsPodJson := ocAPiCall(cluster, token, namespace, "pods", "")

	var istagsDcResult map[string]interface{}
	var istagsJobResult map[string]interface{}
	var istagsCronjobResult map[string]interface{}
	var istagsPodResult map[string]interface{}
	result := map[string]map[string]interface{}{}

	var err error
	err = json.Unmarshal([]byte(istagsDcJson), &istagsDcResult)
	if err != nil {
		log.Println("ERROR: Query dc", err.Error())
	}
	err = json.Unmarshal([]byte(istagsJobJson), &istagsJobResult)
	if err != nil {
		log.Println("ERROR: Query job", err.Error())
	}
	err = json.Unmarshal([]byte(istagsCronjobJson), &istagsCronjobResult)
	if err != nil {
		log.Println("ERROR: Query cronjob", err.Error())
	}
	err = json.Unmarshal([]byte(istagsPodJson), &istagsPodResult)
	if err != nil {
		log.Println("ERROR: Query pod", err.Error())
	}

	// if err := mergo.Merge(&result, istagsDcResult); err != nil {
	// 	log.Println("ERROR: " + "merge mySha to resultSha" + ": failed: " + err.Error())
	// }
	// if err := mergo.Merge(&result, istagsJobResult); err != nil {
	// 	log.Println("ERROR: " + "merge mySha to resultSha" + ": failed: " + err.Error())
	// }
	// if err := mergo.Merge(&result, istagsCronjobResult); err != nil {
	// 	log.Println("ERROR: " + "merge mySha to resultSha" + ": failed: " + err.Error())
	// }
	// if err := mergo.Merge(&result, istagsPodResult); err != nil {
	// 	log.Println("ERROR: " + "merge mySha to resultSha" + ": failed: " + err.Error())
	// }
	result["dc"] = istagsDcResult
	result["job"] = istagsJobResult
	result["cronjob"] = istagsCronjobResult
	result["pod"] = istagsPodResult
	return result
}
