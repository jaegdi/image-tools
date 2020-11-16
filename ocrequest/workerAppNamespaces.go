package ocrequest

import (
	// "fmt"
	"sync"
	// "time"
	// "log"
)

type T_CLusterAppNamespaces []string
type T_FamilyAppNamespaces map[string]T_CLusterAppNamespaces

type JobAppNamespaces struct {
	id      int
	cluster string
	family  string
}

type ResultAppNamespaces struct {
	job JobAppNamespaces
	// family     string
	cluster    string
	namespaces T_CLusterAppNamespaces
}

var jobsAppNamespaces chan JobAppNamespaces
var jobResultsAppNamespaces chan ResultAppNamespaces

var channelsizeAppNamespaces = 30
var noOfWorkersAppNamespaces = 20

func workerGetAppNamespaces(wg *sync.WaitGroup) {
	for job := range jobsAppNamespaces {
		output := ResultAppNamespaces{job, job.cluster, GetAppNamespacesForFamily(job.cluster, job.family)}
		jobResultsAppNamespaces <- output
	}
	wg.Done()
}

func createWorkerPoolAppNamespaces(noOfWorkersAppNamespaces int) {
	var wg sync.WaitGroup
	for i := 0; i < noOfWorkersAppNamespaces; i++ {
		wg.Add(1)
		go workerGetAppNamespaces(&wg)
	}
	wg.Wait()
	close(jobResultsAppNamespaces)
}

func allocateAppNamespaces(family string) {
	jobNr := 0
	clusters := Clusters.Stages
	for cl := 0; cl < len(clusters); cl++ {
		LogMsg("Start JobAppNamespaces for cluster" + clusters[cl])
		job := JobAppNamespaces{jobNr, clusters[cl], family}
		jobsAppNamespaces <- job
		jobNr++
	}
	LogMsg("close jobsAppNamespaces")
	close(jobsAppNamespaces)
}

// func getResult(istagResult T_usedIstagsResult) { // done chan bool,
// 	for result := range jobResultsAppNamespaces {
// 		t := T_usedIstagsResult{}
// 		MergoNestedMaps(&t, istagResult, result.istags)
// 		istagResult = t
// 	}
// 	// done <- true
// }

func goGetAppNamespacesForFamily(family string) T_FamilyAppNamespaces {

	appNameSpaces := T_FamilyAppNamespaces{}

	jobsAppNamespaces = make(chan JobAppNamespaces, channelsizeAppNamespaces)
	jobResultsAppNamespaces = make(chan ResultAppNamespaces, channelsizeAppNamespaces)
	LogMsg("Allocate and start JobsAppNamespaces")
	go allocateAppNamespaces(family)

	LogMsg("Create Worker Pool AppNamespaces")
	createWorkerPoolAppNamespaces(noOfWorkersAppNamespaces)

	LogMsg("Collect results")
	for jobResult := range jobResultsAppNamespaces {
		t := T_FamilyAppNamespaces{}
		resmap := T_FamilyAppNamespaces{jobResult.cluster: jobResult.namespaces}
		MergoNestedMaps(&t, appNameSpaces, resmap)
		appNameSpaces = t
	}

	return appNameSpaces
}
