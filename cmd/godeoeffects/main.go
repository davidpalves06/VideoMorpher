package main

import (
	"net/http"

	"github.com/davidpalves06/GodeoEffects/internal/handlers"
	"github.com/davidpalves06/GodeoEffects/internal/logger"
)

func main() {
	http.Handle("/static/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/download/", handlers.HandleDownloads)
	http.HandleFunc("/upload", handlers.HandleFileUploads)
	http.HandleFunc("/progress", handlers.HandleProgressUpdates)
	http.HandleFunc("/", handlers.HandleIndexPage)

	logger.Info().Println("Web server starting at port 8080")
	err := http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		logger.Error().Printf("Error starting server : %v\n", err)
	}
}
