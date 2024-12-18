package handlers

import (
	"log"
	"net/http"
	"text/template"
)

func HandleIndexPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	log.Println("Serving HTML to index route")
	html, err := template.ParseFiles("static/index.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal("HTML template failed to load!")
	}
	html.Execute(w, nil)
}
