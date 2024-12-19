package handlers

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"

	"github.com/davidpalves06/GodeoEffects/internal/videoeffects"
)

const UPLOAD_DIRECTORY = "uploads/"

func generateRandomID(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func HandleFileUploads(w http.ResponseWriter, r *http.Request) {
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

	var outputFile string
	var fileBytes []byte = videoeffects.GetFileBytes(file)

	operation := r.FormValue("operation")
	if operation == "motion" {
		motionSpeed, err := strconv.ParseFloat(r.FormValue("motionSpeed"), 32)
		if err != nil {
			log.Println("Motion speed is not a number!", err)
			http.Error(w, "Motion speed is not a number!", http.StatusBadRequest)
			return
		}
		outputFile, err = videoeffects.ChangeVideoMotionSpeed(fileBytes, handler.Filename, UPLOAD_DIRECTORY, channelMapping[processId], float32(motionSpeed))

		if err != nil {
			log.Println("Error processing file!", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	} else if operation == "reverse" {
		//TODO: CHECK OTHER EFFECTS
		log.Println("Reverse not accepted yet!")
		http.Error(w, "Reverse not implemented yet!", http.StatusInternalServerError)
		return
	} else if operation == "conversion" {
		outputFormat := r.FormValue("conversionFormat")
		outputFile, err = videoeffects.ChangeVideoFormat(fileBytes, handler.Filename, UPLOAD_DIRECTORY, channelMapping[processId], outputFormat)

		if err != nil {
			log.Println("Error converting video!")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		log.Println("Operation not recognized!")
		http.Error(w, "Operation not recognized!", http.StatusBadRequest)
	}

	var encodedFileName = url.QueryEscape(filepath.Base(outputFile))
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"processID": processId, "generatedFile": encodedFileName})
}
