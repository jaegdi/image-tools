package ocrequest

import (
	"encoding/json"
	// "fmt"
	// "github.com/imdario/mergo"
	"log"
)

// var cluster string = "cid-apc0"
// var namespace string = "pkp-unicorn-drei"

func filterDcResults(dc T_DcResults) T_usedIstagsResult {
	results := T_usedIstagsResult{}
	// ToDo Implementation
	return results
}

func ocGetAllUsedIstagsOfNamespace(cluster string, namespace string) T_runningObjects {
	istagsDcJson := ocAPiCall(cluster, namespace, "deploymentconfigs", "")
	istagsJobJson := ocAPiCall(cluster, namespace, "jobs", "")
	istagsCronjobJson := ocAPiCall(cluster, namespace, "cronjobs", "")
	istagsPodJson := ocAPiCall(cluster, namespace, "pods", "")

	var istagsDcResult T_DcResults
	var istagsJobResult T_JobResults
	var istagsCronjobResult T_CronjobResults
	var istagsPodResult T_Results
	result := T_runningObjects{}

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
	result.Dc = istagsDcResult
	result.Job = istagsJobResult
	result.Cronjob = istagsCronjobResult
	result.Pod = istagsPodResult
	return result
}
