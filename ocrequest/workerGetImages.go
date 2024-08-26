package ocrequest

import (
	"sync"
)

type JobExistingImages struct {
	id      int
	cluster T_clName
}

type ResultExistingImages struct {
	job    JobExistingImages
	images T_ImagesMapAllClusters
}

var jobsExistingImages chan JobExistingImages
var jobResultsExistingImages chan ResultExistingImages

var channelsizeExistingImages = 10
var noOfWorkersExistingImages = 5

// workerExistingImages processes jobs from the jobsExistingImages channel.
// It retrieves all images for the specified cluster and sends the results to the jobResultsExistingImages channel.
//
// Parameters:
// - wg: A pointer to a sync.WaitGroup to signal when the worker is done.
func workerExistingImages(wg *sync.WaitGroup) {
	for job := range jobsExistingImages {
		// Retrieve all images for the specified cluster
		output := ResultExistingImages{
			job:    job,
			images: T_ImagesMapAllClusters{job.cluster: GetAllImagesOfCluster(job.cluster)},
		}
		// Send the results to the jobResultsExistingImages channel
		jobResultsExistingImages <- output
	}
	// Signal that the worker is done
	wg.Done()
}

// createWorkerPoolExistingImages creates a pool of workers to process jobs from the jobsExistingImages channel.
// It waits for all workers to complete and then closes the jobResultsExistingImages channel.
//
// Parameters:
// - noOfWorkersExistingImages: The number of workers to create.
func createWorkerPoolExistingImages(noOfWorkersExistingImages int) {
	var wg sync.WaitGroup
	// Create the specified number of workers
	for i := 0; i < noOfWorkersExistingImages; i++ {
		wg.Add(1)
		go workerExistingImages(&wg)
	}
	// Wait for all workers to complete
	wg.Wait()
	// Close the jobResultsExistingImages channel
	close(jobResultsExistingImages)
}

// allocateExistingImages creates jobs for each cluster in the specified list.
// It sends the jobs to the jobsExistingImages channel and then closes the channel.
//
// Parameters:
// - clusters: A slice of cluster names (T_clName) to process.
func allocateExistingImages(clusters []T_clName) {
	jobNr := 0
	// Iterate through each cluster in the list
	for cl := 0; cl < len(clusters); cl++ {
		VerifyMsg("Start JobExistingImages for cluster" + clusters[cl])
		// Create a new JobExistingImages struct with the current job number and cluster
		job := JobExistingImages{jobNr, clusters[cl]}
		// Send the job to the jobsExistingImages channel
		jobsExistingImages <- job
		jobNr++
	}
	VerifyMsg("close jobsExistingImages")
	// Close the jobsExistingImages channel
	close(jobsExistingImages)
}

// goGetExistingImagesInAllClusters retrieves all images for all clusters in the specified family.
// It processes the clusters concurrently using workers.
//
// Returns:
// - A T_ImagesMapAllClusters containing the retrieved images.
func goGetExistingImagesInAllClusters() T_ImagesMapAllClusters {
	istagResult := T_ImagesMapAllClusters{}

	// Initialize channels for jobs and job results
	jobsExistingImages = make(chan JobExistingImages, channelsizeExistingImages)
	jobResultsExistingImages = make(chan ResultExistingImages, channelsizeExistingImages)

	VerifyMsg("Allocate and start JobsExistingImages")
	// Start allocating jobs for the specified clusters
	go allocateExistingImages(FamilyNamespaces[CmdParams.Family].Stages)

	VerifyMsg("Create Worker Pool for Existing Images")
	// Create a pool of workers to process the jobs
	createWorkerPoolExistingImages(noOfWorkersExistingImages)

	VerifyMsg("Collect results for existing images")
	// Collect results from the job results channel
	for result := range jobResultsExistingImages {
		// Merge the result into the final istagResult
		MergoNestedMaps(&istagResult, result.images)
	}
	return istagResult
}
