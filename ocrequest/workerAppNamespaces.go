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
	family  T_familyName
}

type ResultAppNamespaces struct {
	job        JobAppNamespaces
	cluster    T_clName
	namespaces T_CLusterAppNamespaces
}

var jobsAppNamespaces chan JobAppNamespaces
var jobResultsAppNamespaces chan ResultAppNamespaces

var channelsizeAppNamespaces = 30
var noOfWorkersAppNamespaces = 10

// workerProcessJobsOfGetAppNamespaces processes jobs from the jobsAppNamespaces channel.
// For each job, it retrieves the namespaces for the specified family and cluster,
// and sends the result to the jobResultsAppNamespaces channel.
//
// Parameters:
// - wg: A WaitGroup to synchronize the completion of goroutines.
func workerProcessJobsOfGetAppNamespaces(wg *sync.WaitGroup) {
	for job := range jobsAppNamespaces {
		VerifyMsg("get output GetAppNamespacesForFamily", job)
		output := ResultAppNamespaces{job, job.cluster, GetAppNamespacesForFamily(job.cluster, job.family)}
		jobResultsAppNamespaces <- output
	}
	wg.Done()
}

// createWorkerPoolAppNamespaces creates a pool of worker goroutines to process jobs.
// It starts the specified number of worker goroutines and waits for them to complete.
//
// Parameters:
// - noOfWorkersAppNamespaces: The number of worker goroutines to start.
func createWorkerPoolAppNamespaces(noOfWorkersAppNamespaces int) {
	var wg sync.WaitGroup
	for i := 0; i < noOfWorkersAppNamespaces; i++ {
		wg.Add(1)
		VerifyMsg("Try go workerProcessJobsOfGetAppNamespaces")
		go workerProcessJobsOfGetAppNamespaces(&wg)
	}
	wg.Wait()
	close(jobResultsAppNamespaces)
}

// createJobsToAllocateAppNamespaces creates jobs for each cluster in the specified family.
// It sends the jobs to the jobsAppNamespaces channel and then closes the channel.
//
// Parameters:
// - family: The family name (T_familyName) to filter the namespaces.
func createJobsToAllocateAppNamespaces(family T_familyName) {
	// Initialize the job number counter
	jobNr := 0
	// Get the list of clusters for the specified family
	clusters := FamilyNamespaces[family].Stages
	// Iterate through each cluster in the list
	for cl := 0; cl < len(clusters); cl++ {
		// Log a verification message indicating the start of a job for the current cluster
		VerifyMsg("Start JobAppNamespaces for cluster" + clusters[cl])
		// Create a new JobAppNamespaces struct with the current job number, cluster, and family
		job := JobAppNamespaces{jobNr, clusters[cl], family}
		// Send the job to the jobsAppNamespaces channel
		jobsAppNamespaces <- job
		// Increment the job number counter
		jobNr++
	}
	// Log a verification message indicating the closure of the jobsAppNamespaces channel
	VerifyMsg("close jobsAppNamespaces")
	// Close the jobsAppNamespaces channel
	close(jobsAppNamespaces)
}

// goGetAppNamespacesForFamily retrieves namespaces for all clusters in a specified family.
// It initializes the job and result channels, allocates jobs, creates a worker pool,
// and collects the results into a T_FamilyAppNamespaces map.
//
// Parameters:
// - family: The family name (T_familyName) to filter the namespaces.
//
// Returns:
// - T_FamilyAppNamespaces: A map containing the namespaces for each cluster in the family.
func goGetAppNamespacesForFamily(family T_familyName) T_FamilyAppNamespaces {

	appNameSpaces := T_FamilyAppNamespaces{}

	// Initialize the job and result channels with the specified buffer size
	jobsAppNamespaces = make(chan JobAppNamespaces, channelsizeAppNamespaces)
	jobResultsAppNamespaces = make(chan ResultAppNamespaces, channelsizeAppNamespaces)

	VerifyMsg("Allocate and start JobsAppNamespaces")
	// Start a goroutine to allocate jobs for each cluster in the family
	go createJobsToAllocateAppNamespaces(family)

	VerifyMsg("Try to Create Worker Pool AppNamespaces")
	// Create a pool of worker goroutines to process the jobs
	createWorkerPoolAppNamespaces(noOfWorkersAppNamespaces)

	VerifyMsg("Collect results:", jobResultsAppNamespaces)
	// Collect the results from the jobResultsAppNamespaces channel
	for jobResult := range jobResultsAppNamespaces {
		resmap := T_FamilyAppNamespaces{jobResult.cluster: jobResult.namespaces}
		// Merge the results into the final appNameSpaces map
		MergoNestedMaps(&appNameSpaces, resmap)
	}
	VerifyMsg("\nAppNAmespaces:", appNameSpaces)
	return appNameSpaces
}
