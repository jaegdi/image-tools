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

var channelsizeUsedIstags = 800
var noOfWorkersUsedIstags = 200

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

func allocateUsedIstags(clusters []T_clName, clusterAppNamespaces T_FamilyAppNamespaces) {
	jobNr := 0
	for _, cluster := range clusters {
		VerifyMsg("Start JobUsedIstags for cluster"+cluster, "Namespaces:", clusterAppNamespaces)
		namespaces := clusterAppNamespaces[cluster]
		for _, namespace := range namespaces {
			VerifyMsg("Start job for cluster " + string(cluster) + " in namespace " + string(namespace))
			job := JobUsedIstags{jobNr, cluster, namespace}
			jobsUsedIstags <- job
			jobNr++
		}
	}
	VerifyMsg("close jobsUsedIstags")
	close(jobsUsedIstags)
}

func goGetUsedIstagsForFamilyInAllClusters(family T_familyName) T_usedIstagsResult {

	istagResult := T_usedIstagsResult{}
	allClusterFamilyNamespaces := goGetAppNamespacesForFamily(family)

	if CmdParams.Filter.Namespace == "" {
		VerifyMsg("Start without Namespace Filter")
		jobsUsedIstags = make(chan JobUsedIstags, channelsizeUsedIstags)
		jobResultsUsedIstags = make(chan ResultUsedIstags, channelsizeUsedIstags)

		VerifyMsg("Allocate and start JobsUsedIstags")
		go allocateUsedIstags(FamilyNamespaces[family].Stages, allClusterFamilyNamespaces)

		VerifyMsg("Create Worker Pool for Used IsTags")
		createWorkerPoolUsedIstags(noOfWorkersUsedIstags)

		VerifyMsg("Collect results for Used IsTags")
		VerifyMsg("jobResultsUsedIstags:", jobResultsUsedIstags)
		for result := range jobResultsUsedIstags {
			VerifyMsg("Result for job ", result.job.id, "content", result)
			MergoNestedMaps(&istagResult, result.istags)
		}

	} else {
		VerifyMsg("Start with Namespace Filter", CmdParams.Filter.Namespace)
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
	VerifyMsg(istagResult)
	return istagResult
}
