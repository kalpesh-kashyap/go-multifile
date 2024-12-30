package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"url-file-save/constant"
	"url-file-save/controller"
)

func worker(id int, files <-chan string, results chan<- string) {
	for file := range files {
		log.Printf("worker id:%d", id)
		if _, err := controller.DownloadAndSaveFile(file, constant.FILE_DOWNLOAD_PATH); err != nil {
			results <- "worked with id" + strconv.Itoa(id) + "failed to process file with error" + err.Error()
			return
		}
		results <- "woker id" + strconv.Itoa(id) + "is sucess with url" + file
	}
}

func WorkerPoolHandler(w http.ResponseWriter, r *http.Request) {
	var urls []string
	err := json.NewDecoder(r.Body).Decode(&urls)

	if err != nil {
		http.Error(w, "Error in parsing body "+err.Error(), http.StatusInternalServerError)
		return
	}

	files := make(chan string, len(urls))
	results := make(chan string, len(urls))
	var numberOforkers int = 4

	for i := 1; i <= numberOforkers; i++ {
		go worker(i, files, results)
	}

	for _, file := range urls {
		files <- file
	}

	close(files)

	var allResults []string

	for result := range results {
		allResults = append(allResults, result)
	}

	close(results)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(allResults)
}
