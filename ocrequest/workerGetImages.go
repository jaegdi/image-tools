package ocrequest

import (
	// "fmt"
	"sync"
	// "time"
	"log"
)

type JobExistingImages struct {
	id      int
	cluster string
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

func allocateExistingImages(clusters []string) {
	jobNr := 0
	for cl := 0; cl < len(clusters); cl++ {
		log.Println("Start JobExistingImages for cluster" + clusters[cl])
		job := JobExistingImages{jobNr, clusters[cl]}
		jobsExistingImages <- job
		jobNr++
	}
	log.Println("close jobsExistingImages")
	close(jobsExistingImages)
}

// func getResult(istagResult T_ImagesMap) { // done chan bool,
// 	for result := range jobResultsExistingImages {
// 		t := T_ImagesMap{}
// 		MergoNestedMaps(&t, istagResult, result.images)
// 		istagResult = t
// 	}
// 	// done <- true
// }

func goGetExistingImagesInAllClusters() T_ImagesMapAllClusters {

	istagResult := T_ImagesMapAllClusters{}

	jobsExistingImages = make(chan JobExistingImages, channelsizeExistingImages)
	jobResultsExistingImages = make(chan ResultExistingImages, channelsizeExistingImages)

	log.Println("Allocate and start JobsExistingImages")
	go allocateExistingImages(Clusters.Stages)

	log.Println("Create Worker Pool for Existing Images")
	createWorkerPoolExistingImages(noOfWorkersExistingImages)

	log.Println("Collect results for existing images")
	for result := range jobResultsExistingImages {
		t := T_ImagesMapAllClusters{}
		MergoNestedMaps(&t, istagResult, result.images)
		istagResult = t
	}
	return istagResult
}
