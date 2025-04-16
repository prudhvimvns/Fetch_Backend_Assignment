package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Route registrations
	r.HandleFunc("/receipts/process", ProcessReceiptHandler).Methods("POST")
	r.HandleFunc("/receipts/{id}/points", GetPointsHandler).Methods("GET")

	log.Println("🚀 Server starting on http://localhost:9080")
	if err := http.ListenAndServe(":9080", r); err != nil {
		log.Fatalf("❌ Server failed to start: %v", err)
	}
}
