package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"url-file-save/constant"
	"url-file-save/controller"
)

func worker(wg *sync.WaitGroup, id int, files <-chan string, results chan<- string) {
	defer wg.Done()
	for file := range files {
		log.Printf("Worker ID: %d started processing file: %s", id, file)
		if _, err := controller.DownloadAndSaveFile(file, constant.FILE_DOWNLOAD_PATH); err != nil {
			results <- "Worker ID " + strconv.Itoa(id) + " failed to process file " + file + " with error: " + err.Error()
		} else {
			results <- "woker id" + strconv.Itoa(id) + "is sucess with url" + file
		}

	}
}

func WorkerPoolHandler(w http.ResponseWriter, r *http.Request) {
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
		go worker(&wg, i, files, results)
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
