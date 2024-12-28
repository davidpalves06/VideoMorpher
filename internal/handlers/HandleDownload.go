package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/davidpalves06/GodeoEffects/internal/logger"
)

func HandleDownloads(w http.ResponseWriter, r *http.Request) {
	logger.Info().Println("Download Request received")
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		logger.Error().Println("File name is missing from request. Request failed")
		http.Error(w, "File name is missing", http.StatusBadRequest)
		return
	}

	logger.Debug().Printf("Download request for file %s\n", fileName)

	filePath := filepath.Join("./uploads", fileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		logger.Error().Println("File not found on uploads directory. Request failed")
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		logger.Error().Println("Error while opening file to download. Request failed")
		http.Error(w, "Error opening file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	shouldStream := r.URL.Query().Get("stream") == "enabled"
	if shouldStream {
		logger.Debug().Printf("Streaming file %s\n", fileName)
		w.Header().Set("Content-Type", "video")
		w.Header().Set("Accept-Ranges", "bytes")
	} else {
		logger.Debug().Printf("Responding with file %s\n", fileName)
		w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	stat, _ := file.Stat()
	http.ServeContent(w, r, filePath, stat.ModTime(), file)
	logger.Info().Printf("Download Request handled\n")
}
