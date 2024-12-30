package handler

import (
	"encoding/json"
	"net/http"
	"url-file-save/constant"
	"url-file-save/controller"
)

func FileHandler(w http.ResponseWriter, r *http.Request) {
	var urls []string

	err := json.NewDecoder(r.Body).Decode(&urls)
	if err != nil {
		http.Error(w, "failed to parse body"+err.Error(), http.StatusInternalServerError)
		return
	}
	file, err := controller.DownloadAndSaveFile(urls[0], constant.FILE_DOWNLOAD_PATH)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(file)

}
