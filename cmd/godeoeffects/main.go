package main

import (
	"net/http"
	"os"

	"github.com/davidpalves06/GodeoEffects/internal/cleaner"
	"github.com/davidpalves06/GodeoEffects/internal/handlers"
	"github.com/davidpalves06/GodeoEffects/internal/logger"
)

func init() {
	if _, err := os.Stat(handlers.UPLOAD_DIRECTORY); os.IsNotExist(err) {
		err := os.MkdirAll(handlers.UPLOAD_DIRECTORY, os.ModePerm)
		if err != nil {
			logger.Debug().Printf("Error creating upload directory: %v\n", err)
			return
		}
		logger.Debug().Println("Upload Directory created successfully")
	} else {
		logger.Debug().Println("Upload Directory already exists")
	}
	cleaner.StartRoutineToCleanOldUploadFiles()
}

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
