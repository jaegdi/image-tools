package ocrequest

import (
	// "fmt"
	"sync"
	// "time"
	// "log"
)

type JobExistingIstags struct {
	id        int
	cluster   T_clName
	namespace T_nsName
}

// cluster
type T_ResultExistingIstagsOverAllClusters map[T_clName]T_result

type ResultExistingIstags struct {
	job    JobExistingIstags
	Istags T_ResultExistingIstagsOverAllClusters
}

var jobsExistingIstags chan JobExistingIstags
var jobResultsExistingIstags chan ResultExistingIstags

var channelsizeExistingIstags = 60
var noOfWorkersExistingIstags = 50

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

func allocateExistingIstags(family T_familyName, clusters []T_clName) {
	jobNr := 0
	for cl := 0; cl < len(clusters); cl++ {
		InfoMsg("Start JobExistingIstags for cluster" + clusters[cl])
		clusterFamNamespaces := FamilyNamespaces[family].ImageNamespaces[clusters[cl]]
		for i := 0; i < len(clusterFamNamespaces); i++ {
			InfoMsg("Start job for cluster " + string(clusters[cl]) + " in namespace " + string(clusterFamNamespaces[i]))
			job := JobExistingIstags{jobNr, clusters[cl], clusterFamNamespaces[i]}
			jobsExistingIstags <- job
			jobNr++
		}
	}
	InfoMsg("close jobsExistingIstags")
	close(jobsExistingIstags)
}

func goGetExistingIstagsForFamilyInAllClusters(family T_familyName) T_ResultExistingIstagsOverAllClusters {

	istagResult := T_ResultExistingIstagsOverAllClusters{}

	if CmdParams.Filter.Namespace == "" {

		jobsExistingIstags = make(chan JobExistingIstags, channelsizeExistingIstags)
		jobResultsExistingIstags = make(chan ResultExistingIstags, channelsizeExistingIstags)

		InfoMsg("Allocate and start JobsExistingIstags")
		go allocateExistingIstags(family, FamilyNamespaces[family].Stages)

		InfoMsg("Create Worker Pool")
		createWorkerPoolForExistingIstags(noOfWorkersExistingIstags)

		InfoMsg("Collect results")
		for result := range jobResultsExistingIstags {
			r := result.Istags
			MergoNestedMaps(&istagResult, r)
		}

	} else {
		for _, cluster := range FamilyNamespaces[family].Stages {
			namespaces := FamilyNamespaces[family].ImageNamespaces[cluster]
			for _, namespace := range namespaces {
				r := T_ResultExistingIstagsOverAllClusters{
					cluster: OcGetAllIstagsOfNamespace(T_result{}, cluster, namespace)}
				MergoNestedMaps(&istagResult, r)
			}
		}
	}
	return istagResult
}
