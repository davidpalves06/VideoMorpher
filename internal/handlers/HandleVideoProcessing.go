package handlers

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"github.com/davidpalves06/GodeoEffects/internal/logger"
	"github.com/davidpalves06/GodeoEffects/internal/videoeffects"
)

func generateRandomID(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func HandleFileUploads(w http.ResponseWriter, r *http.Request) {
	logger.Info().Println("File upload request received")
	if r.Method != "POST" {
		logger.Error().Println("Method not allowed. Request failed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var processId = generateRandomID(10)
	channelMapping[processId] = make(chan uint8, 1)

	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_FILE_SIZE)

	err := r.ParseMultipartForm(MAX_UPLOAD_FILE_SIZE)
	if err != nil {
		logger.Error().Println("Failed parsing the form. Request failed")
		http.Error(w, "Unable to parse form: "+err.Error(), http.StatusRequestEntityTooLarge)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		logger.Error().Println("Error retrieving file. Request failed")
		http.Error(w, "Error retrieving the file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var outputFile string

	tmpFile, _ := os.CreateTemp("", "upload-*.tmp")
	logger.Debug().Printf("Temporary file %s created\n", tmpFile.Name())
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, file)

	if err != nil {
		logger.Error().Println("Error creating temporary file. Request failed")
		http.Error(w, "error writing to tmp file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	operation := r.FormValue("operation")
	logger.Debug().Printf("Operation: %s\n", operation)
	if operation == "motion" {
		motionSpeed, err := strconv.ParseFloat(r.FormValue("motionSpeed"), 32)
		if err != nil {
			logger.Error().Println("Motion speed is not a number. Request failed")
			http.Error(w, "Motion speed is not a number!", http.StatusBadRequest)
			return
		}
		logger.Debug().Printf("Changing video motion speed by %f\n", motionSpeed)
		outputFile, err = videoeffects.ChangeVideoMotionSpeed(tmpFile.Name(), handler.Filename, UPLOAD_DIRECTORY, processId, channelMapping[processId], float32(motionSpeed))

		if err != nil {
			logger.Error().Println("Error processing video. Request failed")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	} else if operation == "conversion" {
		outputFormat := r.FormValue("conversionFormat")

		logger.Debug().Printf("Changing video format to %s\n", outputFormat)
		outputFile, err = videoeffects.ChangeVideoFormat(tmpFile.Name(), handler.Filename, UPLOAD_DIRECTORY, processId, channelMapping[processId], outputFormat)

		if err != nil {
			logger.Error().Println("Error converting video. Request failed")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		logger.Error().Println("Operation not recognized. Request failed")
		http.Error(w, "Operation not recognized!", http.StatusBadRequest)
	}

	var encodedFileName = url.QueryEscape(filepath.Base(outputFile))
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"processID": processId, "generatedFile": encodedFileName})
	logger.Info().Println("File upload request handled")
}
