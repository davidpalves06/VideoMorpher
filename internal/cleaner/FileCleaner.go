package cleaner

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/davidpalves06/GodeoEffects/internal/handlers"
	"github.com/davidpalves06/GodeoEffects/internal/logger"
)

func cleanUploadDirectory(currentTime time.Time) {
	logger.Debug().Println("Starting to clean UPLOAD DIRECTORY")
	if _, err := os.Stat(handlers.UPLOAD_DIRECTORY); errors.Is(err, os.ErrNotExist) {
		logger.Warn().Println("UPLOAD DIRECTORY does not exist in file system")
		return
	}
	var deadlineTime = currentTime.Add(-7 * time.Minute)
	dirEntries, err := os.ReadDir(handlers.UPLOAD_DIRECTORY)
	if err != nil {
		logger.Warn().Printf("Error reading UPLOAD DIRECTORY. %s\n", err.Error())
		return
	}
	for _, dirEntry := range dirEntries {
		file, err := dirEntry.Info()
		if err != nil {
			logger.Warn().Printf("Error getting info for file %s: %s", dirEntry.Name(), err.Error())
			return
		}
		var lastFileModTime = file.ModTime()
		filePath := filepath.Join(handlers.UPLOAD_DIRECTORY, file.Name())
		if lastFileModTime.Before(deadlineTime) {
			logger.Debug().Printf("Removing file %s for being too old\n", filePath)
			err = os.Remove(filePath)
			if err != nil {
				logger.Warn().Printf("Could not delete file %s : %s\n", filePath, err.Error())
			}
		}
	}
}

func StartRoutineToCleanOldUploadFiles() {
	logger.Debug().Println("Starting routine to clear old files")
	var routineTicker = time.NewTicker(1 * time.Minute)
	go func() {
		cleanUploadDirectory(time.Now())
		for currentTime := range routineTicker.C {
			logger.Debug().Println("Ticker fired")
			cleanUploadDirectory(currentTime)
		}
	}()
}
