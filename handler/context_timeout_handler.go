package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
	"url-file-save/constant"
	"url-file-save/controller"
)

func contextWorker(wg *sync.WaitGroup, id int, files <-chan string, results chan<- string, ctx context.Context) {
	defer wg.Done()
	for file := range files {
		taskCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		log.Printf("Worker ID: %d started processing file: %s", id, file)
		err := controller.DownloadAndSaveFileWithContext(taskCtx, file, constant.FILE_DOWNLOAD_PATH)
		cancel()
		if err != nil {
			results <- "Worker ID " + strconv.Itoa(id) + " failed to process file " + file + " with error: " + err.Error()
		} else {
			results <- "woker id" + strconv.Itoa(id) + "is sucess with url" + file
		}
	}
}

func ContextWorkerPoolHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
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
		go contextWorker(&wg, i, files, results, ctx)
	}

	// Send tasks to the files channel
	go func() {
		for _, file := range urls {
			select {
			case <-ctx.Done():
				log.Println("Task submission canceled:", ctx.Err())
				return
			case files <- file:
			}
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
