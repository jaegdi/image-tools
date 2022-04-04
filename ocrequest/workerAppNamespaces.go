package ocrequest

import (
	// "fmt"
	"sync"
	// "time"
	// "log"
)

type T_CLusterAppNamespaces []T_nsName
type T_FamilyAppNamespaces map[T_clName]T_CLusterAppNamespaces

type JobAppNamespaces struct {
	id      int
	cluster T_clName
	family  T_family
}

type ResultAppNamespaces struct {
	job JobAppNamespaces
	// family     string
	cluster    T_clName
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

func allocateAppNamespaces(family T_family) {
	jobNr := 0
	clusters := FamilyNamespaces[family].Stages
	for cl := 0; cl < len(clusters); cl++ {
		InfoLogger.Println("Start JobAppNamespaces for cluster" + clusters[cl])
		job := JobAppNamespaces{jobNr, clusters[cl], family}
		jobsAppNamespaces <- job
		jobNr++
	}
	InfoLogger.Println("close jobsAppNamespaces")
	close(jobsAppNamespaces)
}

func goGetAppNamespacesForFamily(family T_family) T_FamilyAppNamespaces {

	appNameSpaces := T_FamilyAppNamespaces{}

	jobsAppNamespaces = make(chan JobAppNamespaces, channelsizeAppNamespaces)
	jobResultsAppNamespaces = make(chan ResultAppNamespaces, channelsizeAppNamespaces)
	InfoLogger.Println("Allocate and start JobsAppNamespaces")
	go allocateAppNamespaces(family)

	InfoLogger.Println("Create Worker Pool AppNamespaces")
	createWorkerPoolAppNamespaces(noOfWorkersAppNamespaces)

	InfoLogger.Println("Collect results")
	for jobResult := range jobResultsAppNamespaces {
		resmap := T_FamilyAppNamespaces{jobResult.cluster: jobResult.namespaces}
		MergoNestedMaps(&appNameSpaces, resmap)
	}

	return appNameSpaces
}
