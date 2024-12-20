package main

import (
	"log"
	"net/http"

	"github.com/davidpalves06/GodeoEffects/internal/handlers"
)

func main() {
	http.Handle("/static/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/download/", handlers.HandleDownloads)
	http.HandleFunc("/upload", handlers.HandleFileUploads)
	http.HandleFunc("/progress", handlers.HandleProgressUpdates)
	http.HandleFunc("/", handlers.HandleIndexPage)

	log.Println("Web Server starting to listen at port 8080")
	err := http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		log.Printf("Error starting server: %v", err)
	} else {
		log.Println("Web Server started at port 8080")
	}
}
