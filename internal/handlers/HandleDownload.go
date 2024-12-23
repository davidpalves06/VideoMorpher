package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func HandleDownloads(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "File name is missing", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join("./uploads", fileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Error opening file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var format = filepath.Ext(filePath)[1:]

	shouldStream := r.URL.Query().Get("stream") == "enabled"
	if shouldStream {
		if format != "mp4" && format != "webm" {
			http.Error(w, "format not accepted to stream", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", fmt.Sprintf("video/%s", format))
		w.Header().Set("Accept-Ranges", "bytes")
	} else {
		w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	stat, _ := file.Stat()
	http.ServeContent(w, r, filePath, stat.ModTime(), file)
}
