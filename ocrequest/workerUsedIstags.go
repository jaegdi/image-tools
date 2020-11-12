package ocrequest

import (
	// "fmt"
	"sync"
	// "time"
	"log"
)

type JobUsedIstags struct {
	id        int
	cluster   string
	namespace string
}

type ResultUsedIstags struct {
	job    JobUsedIstags
	istags T_usedIstagsResult
}

var jobsUsedIstags chan JobUsedIstags
var jobResultsUsedIstags chan ResultUsedIstags

var channelsizeUsedIstags = 300
var noOfWorkersUsedIstags = 100

func workerUsedIstags(wg *sync.WaitGroup) {
	for job := range jobsUsedIstags {
		output := ResultUsedIstags{job, ocGetAllUsedIstagsOfNamespace(job.cluster, job.namespace)}
		jobResultsUsedIstags <- output
	}
	wg.Done()
}

func createWorkerPool(noOfWorkersUsedIstags int) {
	var wg sync.WaitGroup
	for i := 0; i < noOfWorkersUsedIstags; i++ {
		wg.Add(1)
		go workerUsedIstags(&wg)
	}
	wg.Wait()
	close(jobResultsUsedIstags)
}

func allocateUsedIstags(clusters []string, clusterAppNamsepaces T_FamilyAppNamespaces) {
	jobNr := 0
	for cl := 0; cl < len(clusters); cl++ {
		log.Println("Start JobUsedIstags for cluster" + clusters[cl])
		namespaces := clusterAppNamsepaces[clusters[cl]]
		for i := 0; i < len(namespaces); i++ {
			log.Println("Start job for cluster " + clusters[cl] + " in namespace " + namespaces[i])
			job := JobUsedIstags{jobNr, clusters[cl], namespaces[i]}
			jobsUsedIstags <- job
			jobNr++
		}
	}
	log.Println("close jobsUsedIstags")
	close(jobsUsedIstags)
}

// func getResult(istagResult T_usedIstagsResult) { // done chan bool,
// 	for result := range jobResultsUsedIstags {
// 		t := T_usedIstagsResult{}
// 		MergoNestedMaps(&t, istagResult, result.istags)
// 		istagResult = t
// 	}
// 	// done <- true
// }

func goGetUsedIstagsForFamily(family string) T_usedIstagsResult {

	istagResult := T_usedIstagsResult{}
	allClusterFamilyNamespaces := goGetAppNamespacesForFamily(family)

	if CmdParams.Filter.Namespace == "" {

		jobsUsedIstags = make(chan JobUsedIstags, channelsizeUsedIstags)
		jobResultsUsedIstags = make(chan ResultUsedIstags, channelsizeUsedIstags)

		log.Println("Allocate and start JobsUsedIstags")
		go allocateUsedIstags(Clusters.Stages, allClusterFamilyNamespaces)

		log.Println("Create Worker Pool")
		createWorkerPool(noOfWorkersUsedIstags)

		log.Println("Collect results")
		for result := range jobResultsUsedIstags {
			t := T_usedIstagsResult{}
			MergoNestedMaps(&t, istagResult, result.istags)
			istagResult = t
		}

	} else {
		for _, cluster := range Clusters.Stages {
			namespaces := GetAppNamespacesForFamily(cluster, family)
			for _, namespace := range namespaces {
				r := ocGetAllUsedIstagsOfNamespace(cluster, namespace)
				t := T_usedIstagsResult{}
				MergoNestedMaps(&t, istagResult, r)
				istagResult = t
			}
		}
	}
	return istagResult
}
