package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/davidpalves06/GodeoEffects/internal/cleaner"
	"github.com/davidpalves06/GodeoEffects/internal/config"
	"github.com/davidpalves06/GodeoEffects/internal/handlers"
	"github.com/davidpalves06/GodeoEffects/internal/logger"
)

func init() {

	config.LoadConfigurations()

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
	serverConfig := config.ApplicationConfig.ServerConfig
	if serverConfig.Host == "" || serverConfig.Port <= 0 {
		logger.Error().Fatalln("Server configs are not valid")
	}
	http.Handle("/static/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/download/", handlers.HandleDownloads)
	http.HandleFunc("/upload", handlers.HandleFileUploads)
	http.HandleFunc("/progress", handlers.HandleProgressUpdates)
	http.HandleFunc("/", handlers.HandleIndexPage)

	logger.Info().Printf("Web server starting at port %d\n", serverConfig.Port)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port), nil)
	if err != nil {
		logger.Error().Printf("Error starting server : %v\n", err)
	}
}
