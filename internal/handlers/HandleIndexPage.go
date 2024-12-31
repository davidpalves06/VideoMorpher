package handlers

import (
	"net/http"
	"text/template"

	"github.com/davidpalves06/GodeoEffects/internal/logger"
)

func HandleIndexPage(w http.ResponseWriter, r *http.Request) {
	logger.Info().Println("Index Page request received")

	if r.Method != "GET" {
		logger.Error().Println("Method not allowed. Request failed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-type", "text/html")
	html, err := template.ParseFiles("static/index.html")

	if err != nil {
		logger.Error().Println("Error while loading HTML template. Request failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	html.Execute(w, nil)
	logger.Info().Println("Index Page request handled")
}
