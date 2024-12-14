package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/davidpalves06/GodeoEffects/internal/videoeffects"
)

const MAX_UPLOAD_FILE_SIZE int64 = 1 << 30 // 1GB

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
		w.WriteHeader(405)
		fmt.Fprint(w, "Path not found.")
		return
	}

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

	videoeffects.VideoConversion(fileBytes, tmpOutputFile)

	log.Println("Length of file bytes ", len(fileBytes))
	fmt.Fprintf(w, "<p>File processed successfully</p><a href='/download?file=%v'>Download file</a>", outputFile)

}

func main() {

	http.Handle("/static/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/download/", handleDownloads)
	http.HandleFunc("/upload", handleFileUploads)
	http.HandleFunc("/", handleIndexPage)

	log.Println("Web Server starting to listen at port 8080")
	err := http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		log.Printf("Error starting server: %v", err)
	} else {
		log.Println("Web Server started at port 8080")
	}
}
