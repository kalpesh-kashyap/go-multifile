package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"url-file-save/constant"
	"url-file-save/controller"
)

func downloadProcess(wg sync.WaitGroup, url string) {
	defer wg.Done()
	_, err := controller.DownloadAndSaveFile(url, constant.FILE_DOWNLOAD_PATH)
	if err != nil {
		fmt.Println("error in file process" + err.Error())
	}

}

func GoRoutineHandler(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	var urls []string
	err := json.NewDecoder(r.Body).Decode(&urls)

	if err != nil {
		http.Error(w, "Error in parsing body "+err.Error(), http.StatusInternalServerError)
		return
	}

	for _, url := range urls {
		wg.Add(1)
		go downloadProcess(wg, url)
	}

	wg.Wait()

	w.Header().Set("Conent-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode("All files processed successfully")
}
