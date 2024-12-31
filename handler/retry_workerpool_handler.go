package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
	"url-file-save/constant"
	"url-file-save/controller"
)

func retryWorker(wg *sync.WaitGroup, id int, files <-chan string, results chan<- string) {
	defer wg.Done()
	maxAttempts := 3

	for file := range files {
		log.Printf("Worker ID: %d started processing file: %s", id, file)
		for i := 1; i <= maxAttempts; i++ {
			if _, err := controller.DownloadAndSaveFile(file, constant.FILE_DOWNLOAD_PATH); err != nil {
				if i == maxAttempts {
					results <- "Worker ID " + strconv.Itoa(id) + " failed to process file " + file + " with error: " + err.Error() + "on attempt" + strconv.Itoa(i)
				} else {
					time.Sleep(time.Second * time.Duration(i))
				}
			} else {
				results <- "woker id" + strconv.Itoa(id) + "is sucess with url" + file
			}
		}

	}
}

func RetryWorkerPoolHandler(w http.ResponseWriter, r *http.Request) {
	var urls []string
	err := json.NewDecoder(r.Body).Decode(&urls)

	if err != nil {
		http.Error(w, "Error in parsing body "+err.Error(), http.StatusInternalServerError)
		return
	}
	var wg sync.WaitGroup
	files := make(chan string, len(urls))
	results := make(chan string, len(urls))
	var numberOforkers int = 4

	for i := 1; i <= numberOforkers; i++ {
		wg.Add(1)
		go retryWorker(&wg, i, files, results)
	}

	// Send tasks to the files channel
	go func() {
		for _, file := range urls {
			files <- file
		}
		close(files)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	var allResults []string
	for result := range results {
		allResults = append(allResults, result)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(allResults)
}
