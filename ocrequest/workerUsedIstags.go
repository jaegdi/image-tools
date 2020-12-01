package ocrequest

import (
	// "fmt"
	"sync"
	// "time"
	// "log"
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
		LogMsg("Start JobUsedIstags for cluster" + clusters[cl])
		namespaces := clusterAppNamsepaces[clusters[cl]]
		for i := 0; i < len(namespaces); i++ {
			LogMsg("Start job for cluster " + clusters[cl] + " in namespace " + namespaces[i])
			job := JobUsedIstags{jobNr, clusters[cl], namespaces[i]}
			jobsUsedIstags <- job
			jobNr++
		}
	}
	LogMsg("close jobsUsedIstags")
	close(jobsUsedIstags)
}

func goGetUsedIstagsForFamilyInAllClusters(family string) T_usedIstagsResult {

	istagResult := T_usedIstagsResult{}
	allClusterFamilyNamespaces := goGetAppNamespacesForFamily(family)

	if CmdParams.Filter.Namespace == "" {

		jobsUsedIstags = make(chan JobUsedIstags, channelsizeUsedIstags)
		jobResultsUsedIstags = make(chan ResultUsedIstags, channelsizeUsedIstags)

		LogMsg("Allocate and start JobsUsedIstags")
		go allocateUsedIstags(Clusters.Stages, allClusterFamilyNamespaces)

		LogMsg("Create Worker Pool")
		createWorkerPool(noOfWorkersUsedIstags)

		LogMsg("Collect results")
		for result := range jobResultsUsedIstags {
			for is := range result.istags {
				for tag := range result.istags[is] {
					a := []T_usedIstag{}
					if istagResult[is] == nil {
						istagResult[is] = map[string][]T_usedIstag{}
					}
					if istagResult[is][tag] == nil {
						istagResult[is][tag] = []T_usedIstag{}
					}
					if istagResult[is][tag] != nil {
						a = istagResult[is][tag]
					}
					for i := range result.istags[is][tag] {
						a = append(a, result.istags[is][tag][i])
					}
					istagResult[is][tag] = a
				}

			}
			// TODO: Mergen funktioniert nicht richtig, merged immer nur zwei namespaces
			// r := result.istags
			// MergoNestedMaps(&istagResult, r)
		}

	} else {
		for _, cluster := range Clusters.Stages {
			namespaces := GetAppNamespacesForFamily(cluster, family)
			for _, namespace := range namespaces {
				if namespace == CmdParams.Filter.Namespace {
					r := ocGetAllUsedIstagsOfNamespace(cluster, namespace)
					MergoNestedMaps(&istagResult, r)
				}
			}
		}
	}
	return istagResult
}
