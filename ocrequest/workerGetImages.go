package ocrequest

import (
	// "fmt"
	"sync"
	// "time"
	// "log"
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

func workerExistingImages(wg *sync.WaitGroup) {
	for job := range jobsExistingImages {
		output := ResultExistingImages{job, T_ImagesMapAllClusters{job.cluster: GetAllImagesOfCluster(job.cluster)}}
		jobResultsExistingImages <- output
	}
	wg.Done()
}

func createWorkerPoolExistingImages(noOfWorkersExistingImages int) {
	var wg sync.WaitGroup
	for i := 0; i < noOfWorkersExistingImages; i++ {
		wg.Add(1)
		go workerExistingImages(&wg)
	}
	wg.Wait()
	close(jobResultsExistingImages)
}

func allocateExistingImages(clusters []T_clName) {
	jobNr := 0
	for cl := 0; cl < len(clusters); cl++ {
		LogMsg("Start JobExistingImages for cluster" + clusters[cl])
		job := JobExistingImages{jobNr, clusters[cl]}
		jobsExistingImages <- job
		jobNr++
	}
	LogMsg("close jobsExistingImages")
	close(jobsExistingImages)
}

func goGetExistingImagesInAllClusters() T_ImagesMapAllClusters {

	istagResult := T_ImagesMapAllClusters{}

	jobsExistingImages = make(chan JobExistingImages, channelsizeExistingImages)
	jobResultsExistingImages = make(chan ResultExistingImages, channelsizeExistingImages)

	LogMsg("Allocate and start JobsExistingImages")
	go allocateExistingImages(Clusters.Stages)

	LogMsg("Create Worker Pool for Existing Images")
	createWorkerPoolExistingImages(noOfWorkersExistingImages)

	LogMsg("Collect results for existing images")
	for result := range jobResultsExistingImages {
		MergoNestedMaps(&istagResult, result.images)
	}
	return istagResult
}
