package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"

	"github.com/davidpalves06/GodeoEffects/internal/videoeffects"
)

const MAX_UPLOAD_FILE_SIZE int64 = 1 << 30 // 1GB
var channelMapping = make(map[string](chan uint8))

func handleDownloads(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, filePath)
}

func handleProgressUpdates(w http.ResponseWriter, r *http.Request) {
	processID := r.URL.Query().Get("processID")
	if processID == "" {
		http.Error(w, "ProcessID is missing", http.StatusBadRequest)
		return
	}
	log.Println("CHECKING PROGRESS FOR", processID)

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	processChannel := channelMapping[processID]

	for percentage := range processChannel {
		fmt.Fprintf(w, "event: progress\n")
		fmt.Fprintf(w, "data: %d%%\n\n", percentage)
		flusher.Flush()
		log.Println("Flushed:", percentage)
		if percentage == 100 {
			delete(channelMapping, processID)
			break
		}
	}

}

func handleIndexPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	log.Println("Serving HTML to index route")
	html, err := template.ParseFiles("static/index.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal("HTML template failed to load!")
	}
	html.Execute(w, nil)
}

func handleFileUploads(w http.ResponseWriter, r *http.Request) {
	log.Println("Upload file request")
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var processId = generateRandomID(10)
	channelMapping[processId] = make(chan uint8, 1)

	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_FILE_SIZE)

	err := r.ParseMultipartForm(MAX_UPLOAD_FILE_SIZE)
	if err != nil {
		w.WriteHeader(413)
		http.Error(w, "Unable to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var outputFile string = "FastMotion-" + handler.Filename
	var tmpOutputFile = "uploads/" + outputFile
	var fileBytes []byte = videoeffects.GetFileBytes(file)

	go videoeffects.VideoConversion(fileBytes, tmpOutputFile, channelMapping[processId])

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"processID": processId})

}

func generateRandomID(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func main() {
	http.Handle("/static/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/download/", handleDownloads)
	http.HandleFunc("/upload", handleFileUploads)
	http.HandleFunc("/progress", handleProgressUpdates)
	http.HandleFunc("/", handleIndexPage)

	log.Println("Web Server starting to listen at port 8080")
	err := http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		log.Printf("Error starting server: %v", err)
	} else {
		log.Println("Web Server started at port 8080")
	}
}
