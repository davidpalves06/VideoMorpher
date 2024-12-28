package handlers

import (
	"net/http"
	"text/template"

	"github.com/davidpalves06/GodeoEffects/internal/logger"
)

func HandleIndexPage(w http.ResponseWriter, r *http.Request) {
	logger.Info().Println("Index Page request received")
	w.Header().Set("Content-type", "text/html")
	html, err := template.ParseFiles("static/index.html")

	if err != nil {
		logger.Error().Println("Error while loading HTML template. Request failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	html.Execute(w, nil)
	logger.Info().Println("Index Page request handled")
}
