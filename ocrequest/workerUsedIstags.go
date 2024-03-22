package ocrequest

import (
	// "fmt"
	"sync"
	// "time"
	// "log"
)

type JobUsedIstags struct {
	id        int
	cluster   T_clName
	namespace T_nsName
}

type ResultUsedIstags struct {
	job    JobUsedIstags
	istags T_usedIstagsResult
}

var jobsUsedIstags chan JobUsedIstags
var jobResultsUsedIstags chan ResultUsedIstags

var channelsizeUsedIstags = 500
var noOfWorkersUsedIstags = 100

func workerUsedIstags(wg *sync.WaitGroup) {
	for job := range jobsUsedIstags {
		output := ResultUsedIstags{job, ocGetAllUsedIstagsOfNamespace(job.cluster, job.namespace)}
		jobResultsUsedIstags <- output
	}
	wg.Done()
}

func createWorkerPoolUsedIstags(noOfWorkersUsedIstags int) {
	var wg sync.WaitGroup
	for i := 0; i < noOfWorkersUsedIstags; i++ {
		wg.Add(1)
		go workerUsedIstags(&wg)
	}
	wg.Wait()
	close(jobResultsUsedIstags)
}

func allocateUsedIstags(clusters []T_clName, clusterAppNamsepaces T_FamilyAppNamespaces) {
	jobNr := 0
	for _, cluster := range clusters {
		InfoLogger.Println("Start JobUsedIstags for cluster" + cluster)
		namespaces := clusterAppNamsepaces[cluster]
		for _, namespace := range namespaces {
			InfoLogger.Println("Start job for cluster " + string(cluster) + " in namespace " + string(namespace))
			job := JobUsedIstags{jobNr, cluster, namespace}
			jobsUsedIstags <- job
			jobNr++
		}
	}
	InfoLogger.Println("close jobsUsedIstags")
	close(jobsUsedIstags)
}

func goGetUsedIstagsForFamilyInAllClusters(family T_familyName) T_usedIstagsResult {

	istagResult := T_usedIstagsResult{}
	allClusterFamilyNamespaces := goGetAppNamespacesForFamily(family)

	if CmdParams.Filter.Namespace == "" {

		jobsUsedIstags = make(chan JobUsedIstags, channelsizeUsedIstags)
		jobResultsUsedIstags = make(chan ResultUsedIstags, channelsizeUsedIstags)

		InfoLogger.Println("Allocate and start JobsUsedIstags")
		go allocateUsedIstags(FamilyNamespaces[family].Stages, allClusterFamilyNamespaces)

		InfoLogger.Println("Create Worker Pool for Used IsTags")
		createWorkerPoolUsedIstags(noOfWorkersUsedIstags)

		InfoLogger.Println("Collect results for Used IsTags")
		for result := range jobResultsUsedIstags {
			// for is := range result.istags {
			// 	for tag := range result.istags[is] {
			// 		a := []T_usedIstag{}
			// 		if istagResult[is] == nil {
			// 			istagResult[is] = map[T_tagName][]T_usedIstag{}
			// 		}
			// 		if istagResult[is][tag] == nil {
			// 			istagResult[is][tag] = []T_usedIstag{}
			// 		}
			// 		if istagResult[is][tag] != nil {
			// 			a = istagResult[is][tag]
			// 		}
			// 		for i := range result.istags[is][tag] {
			// 			a = append(a, result.istags[is][tag][i])
			// 		}
			// 		istagResult[is][tag] = a
			// 	}

			// }
			// TODO: Mergen funktioniert nicht richtig, merged immer nur zwei namespaces
			// r := result.istags
			MergoNestedMaps(&istagResult, result.istags)
		}

	} else {
		for _, cluster := range FamilyNamespaces[family].Stages {
			namespaces := GetAppNamespacesForFamily(cluster, family)
			for _, namespace := range namespaces {
				if namespace == CmdParams.Filter.Namespace {
					r := ocGetAllUsedIstagsOfNamespace(cluster, namespace)
					MergoNestedMaps(&istagResult, r)
				}
			}
		}
	}
	InfoLogger.Println(istagResult)
	return istagResult
}
