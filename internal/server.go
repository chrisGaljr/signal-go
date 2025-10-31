package internal

import (
	"log"
	"net/http"
	"signal/main/internal/handlers"
)

func StartServer() {
	http.HandleFunc("/", handlers.HomePageHandler)
	http.HandleFunc("/about", handlers.AboutPageHandler)
	http.HandleFunc("/status", handlers.StatusPageHandler)
	http.HandleFunc("/encrypt", handlers.EncryptionRequestHandler)
	http.HandleFunc("/incident-history", handlers.IncidentHistoryPageHandler)

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Server started on port 7777")
	http.ListenAndServe(":7777", nil)
}
