package handlers

import (
	"fmt"
	"log"
	"net/http"
)

func HandleProgressUpdates(w http.ResponseWriter, r *http.Request) {
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
		fmt.Fprintf(w, "data: %d\n\n", percentage)
		flusher.Flush()
		if percentage >= 100 {
			delete(channelMapping, processID)
			break
		}
	}

	fmt.Fprintf(w, "event: progress\n")
	fmt.Fprintf(w, "data: %d\n\n", 100)
	flusher.Flush()

}
