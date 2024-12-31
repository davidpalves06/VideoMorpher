package handlers

import (
	"fmt"
	"net/http"

	"github.com/davidpalves06/GodeoEffects/internal/logger"
)

func HandleProgressUpdates(w http.ResponseWriter, r *http.Request) {
	logger.Info().Println("Progress update request for a process received")
	if r.Method != "GET" {
		logger.Error().Println("Method not allowed. Request failed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	processID := r.URL.Query().Get("processID")
	if processID == "" || channelMapping[processID] == nil {
		logger.Error().Println("ProcessID is missing. Request failed")
		http.Error(w, "ProcessID is missing", http.StatusBadRequest)
		return
	}
	logger.Debug().Printf("Process to be tracked is %s\n", processID)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		logger.Error().Println("Streaming unsupported. Request failed")
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	processChannel := channelMapping[processID]

	for percentage := range processChannel {
		if percentage == 255 {
			fmt.Fprintf(w, "event: error\n\n")
			flusher.Flush()
			delete(channelMapping, processID)
			logger.Warn().Println("Sending Error event. Progress update finished before expected")
			return
		} else {
			fmt.Fprintf(w, "event: progress\n")
			fmt.Fprintf(w, "data: %d\n\n", percentage)
			flusher.Flush()
			logger.Debug().Println("Sending update event")
		}
	}

	delete(channelMapping, processID)
	fmt.Fprintf(w, "event: progress\n")
	fmt.Fprintf(w, "data: %d\n\n", 100)
	logger.Debug().Println("Sending update event")
	flusher.Flush()

	logger.Info().Println("Progress update request handled")
}
