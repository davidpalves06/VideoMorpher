package handlers

import (
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

	shouldStream := r.URL.Query().Get("stream") == "enabled"
	if shouldStream {
		//TODO: CHECK FILE FORMAT
		w.Header().Set("Content-Type", "video/mp4")
		w.Header().Set("Accept-Ranges", "bytes")
	} else {
		w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	stat, _ := file.Stat()
	http.ServeContent(w, r, filePath, stat.ModTime(), file)
}
