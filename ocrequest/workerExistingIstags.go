package ocrequest

import (
	"sync"
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

// workerExistingIstags processes jobs from the jobsExistingIstags channel.
// It retrieves image stream tags for the specified namespace and cluster, and sends the results to the jobResultsExistingIstags channel.
//
// Parameters:
// - wg: A pointer to a sync.WaitGroup to signal when the worker is done.
func workerExistingIstags(wg *sync.WaitGroup) {
	for job := range jobsExistingIstags {
		// Retrieve image stream tags for the specified namespace and cluster
		output := ResultExistingIstags{
			job: job,
			Istags: T_ResultExistingIstagsOverAllClusters{
				job.cluster: OcGetAllIstagsOfNamespace(T_result{}, job.cluster, job.namespace),
			},
		}
		// Send the results to the jobResultsExistingIstags channel
		jobResultsExistingIstags <- output
	}
	// Signal that the worker is done
	wg.Done()
}

// createWorkerPoolForExistingIstags creates a pool of workers to process jobs from the jobsExistingIstags channel.
// It waits for all workers to complete and then closes the jobResultsExistingIstags channel.
//
// Parameters:
// - noOfWorkersExistingIstags: The number of workers to create.
func createWorkerPoolForExistingIstags(noOfWorkersExistingIstags int) {
	var wg sync.WaitGroup
	// Create the specified number of workers
	for i := 0; i < noOfWorkersExistingIstags; i++ {
		wg.Add(1)
		go workerExistingIstags(&wg)
	}
	// Wait for all workers to complete
	wg.Wait()
	// Close the jobResultsExistingIstags channel
	close(jobResultsExistingIstags)
}

// allocateExistingIstags creates jobs for each namespace in the specified family and clusters.
// It sends the jobs to the jobsExistingIstags channel and then closes the channel.
//
// Parameters:
// - family: The family name (T_familyName) to filter the namespaces.
// - clusters: A slice of cluster names (T_clName) to process.
func allocateExistingIstags(family T_familyName, clusters []T_clName) {
	jobNr := 0
	// Iterate through each cluster in the list
	for cl := 0; cl < len(clusters); cl++ {
		VerifyMsg("Start JobExistingIstags for cluster" + clusters[cl])
		clusterFamNamespaces := FamilyNamespaces[family].ImageNamespaces[clusters[cl]]
		// Iterate through each namespace in the current cluster
		for i := 0; i < len(clusterFamNamespaces); i++ {
			VerifyMsg("Start job for cluster " + string(clusters[cl]) + " in namespace " + string(clusterFamNamespaces[i]))
			// Create a new JobExistingIstags struct with the current job number, cluster, and namespace
			job := JobExistingIstags{jobNr, clusters[cl], clusterFamNamespaces[i]}
			// Send the job to the jobsExistingIstags channel
			jobsExistingIstags <- job
			jobNr++
		}
	}
	VerifyMsg("close jobsExistingIstags")
	// Close the jobsExistingIstags channel
	close(jobsExistingIstags)
}

// goGetExistingIstagsForFamilyInAllClusters retrieves image stream tags for all namespaces in the specified family and clusters.
// It either processes the namespaces concurrently using workers or sequentially if a specific namespace filter is set.
//
// Parameters:
// - family: The family name (T_familyName) to filter the namespaces.
//
// Returns:
// - A T_ResultExistingIstagsOverAllClusters containing the retrieved image stream tags.
func goGetExistingIstagsForFamilyInAllClusters(family T_familyName) T_ResultExistingIstagsOverAllClusters {
	istagResult := T_ResultExistingIstagsOverAllClusters{}

	if CmdParams.Filter.Namespace == "" {
		// Initialize channels for jobs and job results
		jobsExistingIstags = make(chan JobExistingIstags, channelsizeExistingIstags)
		jobResultsExistingIstags = make(chan ResultExistingIstags, channelsizeExistingIstags)

		VerifyMsg("Allocate and start JobsExistingIstags")
		// Start allocating jobs for the specified family and clusters
		go allocateExistingIstags(family, FamilyNamespaces[family].Stages)

		VerifyMsg("Create Worker Pool")
		// Create a pool of workers to process the jobs
		createWorkerPoolForExistingIstags(noOfWorkersExistingIstags)

		VerifyMsg("Collect results")
		// Collect results from the job results channel
		for result := range jobResultsExistingIstags {
			r := result.Istags
			// Merge the result into the final istagResult
			MergoNestedMaps(&istagResult, r)
		}
	} else {
		// If a specific namespace filter is set, process the namespaces sequentially
		for _, cluster := range FamilyNamespaces[family].Stages {
			namespaces := FamilyNamespaces[family].ImageNamespaces[cluster]
			for _, namespace := range namespaces {
				r := T_ResultExistingIstagsOverAllClusters{
					cluster: OcGetAllIstagsOfNamespace(T_result{}, cluster, namespace),
				}
				// Merge the result into the final istagResult
				MergoNestedMaps(&istagResult, r)
			}
		}
	}
	return istagResult
}
