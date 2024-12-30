package main

import (
	"log"
	"net/http"
	"url-file-save/routes"
)

func main() {
	mux := http.NewServeMux()
	routes.FileRouter(mux)
	log.Println("Server is running on 8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
