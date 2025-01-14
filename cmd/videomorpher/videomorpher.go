package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/davidpalves06/GodeoEffects/internal/cleaner"
	"github.com/davidpalves06/GodeoEffects/internal/config"
	"github.com/davidpalves06/GodeoEffects/internal/handlers"
	"github.com/davidpalves06/GodeoEffects/internal/logger"
	"github.com/davidpalves06/GodeoEffects/internal/videoeffects"
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

	var server = http.Server{
		Addr: fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port),
	}

	http.Handle("/static/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/download/", handlers.HandleDownloads)
	http.HandleFunc("/upload", handlers.HandleFileUploads)
	http.HandleFunc("/progress", handlers.HandleProgressUpdates)
	http.HandleFunc("/", handlers.HandleIndexPage)

	go func() {
		logger.Info().Printf("Web server starting at port %d\n", serverConfig.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error().Printf("Error with HTTP server : %s\n", err.Error())
			os.Exit(1)
		}
	}()

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)

	<-sigChannel

	logger.Info().Println("Starting Shutdown")
	shutdownContext, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()
	if err := server.Shutdown(shutdownContext); err != nil {
		fmt.Printf("Shutdown timeout is done. Finishing all open commands\n")
		for _, cmd := range videoeffects.ActiveCommands {
			cmd.Process.Kill()
		}
		videoeffects.VideoCommandWaitGroup.Wait()
	}

	logger.Info().Println("Shutdown complete")
}
