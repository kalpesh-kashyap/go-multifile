package routes

import (
	"net/http"
	"url-file-save/handler"
)

func FileRouter(mux *http.ServeMux) {
	// mux.HandleFunc("/download", handler.FileHandler)
	mux.HandleFunc("/download", handler.GoRoutineHandler)
}
