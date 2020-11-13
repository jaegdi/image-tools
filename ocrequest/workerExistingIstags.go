package ocrequest

import (
	// "fmt"
	"sync"
	// "time"
	"log"
)

type T_clname string
type T_nsname string

type JobExistingIstags struct {
	id        int
	cluster   string
	namespace string
}

//                                             cluster
type T_ResultExistingIstagsOverAllClusters map[string]T_result

type ResultExistingIstags struct {
	job    JobExistingIstags
	istags T_ResultExistingIstagsOverAllClusters
}

var jobsExistingIstags chan JobExistingIstags
var jobResultsExistingIstags chan ResultExistingIstags

var channelsizeExistingIstags = 30
var noOfWorkersExistingIstags = 15

func workerExistingIstags(wg *sync.WaitGroup) {
	for job := range jobsExistingIstags {
		output := ResultExistingIstags{job,
			T_ResultExistingIstagsOverAllClusters{
				job.cluster: OcGetAllIstagsOfNamespace(T_result{}, job.cluster, job.namespace)}}
		jobResultsExistingIstags <- output
	}
	wg.Done()
}

func createWorkerPoolForExistingIstags(noOfWorkersExistingIstags int) {
	var wg sync.WaitGroup
	for i := 0; i < noOfWorkersExistingIstags; i++ {
		wg.Add(1)
		go workerExistingIstags(&wg)
	}
	wg.Wait()
	close(jobResultsExistingIstags)
}

func allocateExistingIstags(clusters []string, clusterFamNamespaces []string) {
	jobNr := 0
	for cl := 0; cl < len(clusters); cl++ {
		log.Println("Start JobExistingIstags for cluster" + clusters[cl])
		for i := 0; i < len(clusterFamNamespaces); i++ {
			log.Println("Start job for cluster " + clusters[cl] + " in namespace " + clusterFamNamespaces[i])
			job := JobExistingIstags{jobNr, clusters[cl], clusterFamNamespaces[i]}
			jobsExistingIstags <- job
			jobNr++
		}
	}
	log.Println("close jobsExistingIstags")
	close(jobsExistingIstags)
}

// func getResult(istagResult T_usedIstagsResult) { // done chan bool,
// 	for result := range jobResultsExistingIstags {
// 		t := T_usedIstagsResult{}
// 		MergoNestedMaps(&t, istagResult, result.istags)
// 		istagResult = t
// 	}
// 	// done <- true
// }

func goGetExistingIstagsForFamilyInAllClusters(family string) T_ResultExistingIstagsOverAllClusters {

	istagResult := T_ResultExistingIstagsOverAllClusters{}

	if CmdParams.Filter.Namespace == "" {

		jobsExistingIstags = make(chan JobExistingIstags, channelsizeExistingIstags)
		jobResultsExistingIstags = make(chan ResultExistingIstags, channelsizeExistingIstags)

		log.Println("Allocate and start JobsExistingIstags")
		go allocateExistingIstags(Clusters.Stages, FamilyNamespaces[family])

		log.Println("Create Worker Pool")
		createWorkerPoolForExistingIstags(noOfWorkersExistingIstags)

		log.Println("Collect results")
		for result := range jobResultsExistingIstags {
			t := T_ResultExistingIstagsOverAllClusters{}
			MergoNestedMaps(&t, istagResult, result.istags)
			istagResult = t
		}

	} else {
		for _, cluster := range Clusters.Stages {
			namespaces := FamilyNamespaces[family]
			for _, namespace := range namespaces {
				r := T_ResultExistingIstagsOverAllClusters{
					cluster: OcGetAllIstagsOfNamespace(T_result{}, cluster, namespace)}
				t := T_ResultExistingIstagsOverAllClusters{}
				MergoNestedMaps(&t, istagResult, r)
				istagResult = t
			}
		}
	}
	return istagResult
}
