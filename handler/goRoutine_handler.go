package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"url-file-save/constant"
	"url-file-save/controller"
)

func GoRoutineHandler(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	var urls []string
	err := json.NewDecoder(r.Body).Decode(&urls)

	if err != nil {
		http.Error(w, "Error in parsing body "+err.Error(), http.StatusInternalServerError)
		return
	}

	errorChannel := make(chan error, len(urls))
	results := make(chan string, len(urls))

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			if _, err := controller.DownloadAndSaveFile(url, constant.FILE_DOWNLOAD_PATH); err != nil {
				errorChannel <- fmt.Errorf("fail to process on %s: %v", url, err)
				return
			}
			results <- fmt.Sprintf("successfully processed %s", url)
		}(url)
	}

	go func() {
		wg.Wait()
		close(errorChannel)
		close(results)
	}()

	var errorFiles []string
	var successFiles []string

	for err := range errorChannel {
		errorFiles = append(errorFiles, err.Error())
	}

	for success := range results {
		successFiles = append(successFiles, success)
	}
	response := map[string]interface{}{
		"sucess": successFiles,
		"errors": errorFiles,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
