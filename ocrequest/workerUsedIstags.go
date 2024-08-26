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

// workerUsedIstags processes jobs from the jobsUsedIstags channel and sends the results to the jobResultsUsedIstags channel.
// It runs in a loop, processing each job until the jobsUsedIstags channel is closed.
// After processing all jobs, it signals that it is done by calling wg.Done().
//
// Parameters:
// - wg: A WaitGroup used to signal when the worker has finished processing all jobs.
func workerUsedIstags(wg *sync.WaitGroup) {
	for job := range jobsUsedIstags {
		output := ResultUsedIstags{job, ocGetAllUsedIstagsOfNamespace(job.cluster, job.namespace)}
		jobResultsUsedIstags <- output
	}
	wg.Done()
}

// createWorkerPoolUsedIstags initializes a pool of workers to process jobs from the jobsUsedIstags channel.
// It creates a specified number of worker goroutines, each of which processes jobs and sends results to the jobResultsUsedIstags channel.
// The function waits for all workers to finish processing before closing the jobResultsUsedIstags channel.
//
// Parameters:
// - noOfWorkersUsedIstags: The number of worker goroutines to create.
func createWorkerPoolUsedIstags(noOfWorkersUsedIstags int) {
	var wg sync.WaitGroup
	for i := 0; i < noOfWorkersUsedIstags; i++ {
		wg.Add(1)
		go workerUsedIstags(&wg)
	}
	wg.Wait()
	close(jobResultsUsedIstags)
}

// allocateUsedIstags creates jobs for each cluster and its associated namespaces, and sends these jobs to the jobsUsedIstags channel.
// It iterates over the provided clusters and their namespaces, creating a job for each namespace.
// Each job is then sent to the jobsUsedIstags channel for processing by worker goroutines.
// After all jobs have been created and sent, the jobsUsedIstags channel is closed.
//
// Parameters:
// - clusters: A slice of T_clName representing the clusters to process.
// - clusterAppNamespaces: A map of T_FamilyAppNamespaces where each key is a cluster and the value is a slice of namespaces associated with that cluster.
func allocateUsedIstags(clusters []T_clName, clusterAppNamespaces T_FamilyAppNamespaces) {
	jobNr := 0
	for _, cluster := range clusters {
		// Log the start of job creation for the current cluster
		VerifyMsg("Start JobUsedIstags for cluster:", cluster, ", Namespacecount:", len(clusterAppNamespaces[cluster]), ", Namespaces:", clusterAppNamespaces)

		// Get the namespaces associated with the current cluster
		namespaces := clusterAppNamespaces[cluster]

		for _, namespace := range namespaces {
			// Log the start of a job for the current namespace
			VerifyMsg("JobNr:", jobNr, "Start job for cluster "+string(cluster)+" in namespace "+string(namespace))

			// Create a new job and send it to the jobsUsedIstags channel
			job := JobUsedIstags{jobNr, cluster, namespace}
			jobsUsedIstags <- job
			jobNr++
		}
	}

	// Log the closure of the jobsUsedIstags channel
	VerifyMsg("close jobsUsedIstags")
	close(jobsUsedIstags)
}

// goGetUsedIstagsForFamilyInAllClusters retrieves the used image stream tags for a given family across all clusters.
// It either processes all namespaces or filters by a specific namespace based on the CmdParams.Filter.Namespace value.
// The results are merged into a single T_usedIstagsResult and returned.
//
// Parameters:
// - family: The family name for which to retrieve the used image stream tags.
//
// Returns:
// - A T_usedIstagsResult containing the merged results of used image stream tags.
func goGetUsedIstagsForFamilyInAllClusters(family T_familyName) T_usedIstagsResult {

	// Initialize the result map for used image stream tags
	istagResult := T_usedIstagsResult{}
	// Retrieve all namespaces for the given family across all clusters
	allClusterFamilyNamespaces := goGetAppNamespacesForFamily(family)

	// Check if a namespace filter is applied
	if CmdParams.Filter.Namespace == "" {
		VerifyMsg("Start without Namespace Filter")

		// Initialize channels for jobs and job results
		jobsUsedIstags = make(chan JobUsedIstags, channelsizeUsedIstags)
		jobResultsUsedIstags = make(chan ResultUsedIstags, channelsizeUsedIstags)

		VerifyMsg("Allocate and start JobsUsedIstags")
		// Start job allocation in a separate goroutine
		go allocateUsedIstags(FamilyNamespaces[family].Stages, allClusterFamilyNamespaces)

		VerifyMsg("Create Worker Pool for Used IsTags")
		// Create a pool of worker goroutines to process the jobs
		createWorkerPoolUsedIstags(noOfWorkersUsedIstags)

		VerifyMsg("Collect results for Used IsTags")
		VerifyMsg("jobResultsUsedIstags:", jobResultsUsedIstags)
		// Collect and merge results from the job results channel
		count := 0
		for result := range jobResultsUsedIstags {

			for is, ismap := range result.istags {
				for tag, tagmap := range ismap {
					VerifyMsg("Count jobResultsUsedIstags", count,
						"Result for job ", result.job.id,
						"tagCount:", len(result.istags))
					for _, usedIsTag := range tagmap {
						// Add the used ImageStreamTag to the results map
						istagResult = AddUsedIstag(istagResult, is, tag, usedIsTag)
					}
				}
			}
			// Increment the count variable by 1 for each result
			count++
		}

	} else {
		VerifyMsg("Start with Namespace Filter", CmdParams.Filter.Namespace)
		// Process only the namespaces that match the filter
		count := 0
		for _, cluster := range FamilyNamespaces[family].Stages {
			namespaces := GetAppNamespacesForFamily(cluster, family)
			for _, namespace := range namespaces {
				if namespace == CmdParams.Filter.Namespace {
					// Retrieve used image stream tags for the specific namespace
					result := ocGetAllUsedIstagsOfNamespace(cluster, namespace)
					for is, ismap := range result {
						for tag, tagmap := range ismap {
							VerifyMsg("Count jobResultsUsedIstags", count,
								"Result for job namespace", namespace, "is", is, "tag", tag,
								"tagCount:", len(tagmap))
							count++
							for _, usedIsTag := range tagmap {
								// Add the used ImageStreamTag to the results map
								istagResult = AddUsedIstag(istagResult, is, tag, usedIsTag)
							}
						}
					}
				}
			}
		}
	}
	// Log the final result
	VerifyMsg("Result of goGetUsedIstagsForFamilyInAllClusters, Anzahl:", len(istagResult), "Map", GetJsonFromMap(istagResult))
	return istagResult
}
