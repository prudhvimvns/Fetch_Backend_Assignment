package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var receiptStore = make(map[string]int)

func ProcessReceiptHandler(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "Request body is empty. Please verify input.", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var receipt Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		log.Printf("❌ JSON decode error: %v", err)
		http.Error(w, "Invalid JSON format. Please verify input.", http.StatusBadRequest)
		return
	}

	if err := validateReceipt(receipt); err != nil {
		log.Printf("❌ Receipt validation failed: %v", err)
		http.Error(w, "Invalid receipt payload. Please verify input.", http.StatusBadRequest)
		return
	}

	points := CalculatePoints(receipt)
	id := uuid.New().String()
	receiptStore[id] = points

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"id": id}); err != nil {
		log.Printf("❌ Response encoding failed: %v", err)
		http.Error(w, "Could not encode response.", http.StatusInternalServerError)
	}
}

func GetPointsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, exists := vars["id"]
	if !exists || id == "" {
		http.Error(w, "Missing receipt ID in path.", http.StatusBadRequest)
		return
	}

	points, found := receiptStore[id]
	if !found {
		http.Error(w, "No receipt found for that ID.", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]int{"points": points}); err != nil {
		log.Printf("❌ Response encoding failed: %v", err)
		http.Error(w, "Could not encode response.", http.StatusInternalServerError)
	}
}

func validateReceipt(r Receipt) error {
	if r.Retailer == "" || r.PurchaseDate == "" || r.PurchaseTime == "" || r.Total == "" {
		return errors.New("missing required fields")
	}
	if len(r.Items) == 0 {
		return errors.New("receipt must contain at least one item")
	}
	for _, item := range r.Items {
		if item.ShortDescription == "" || item.Price == "" {
			return errors.New("each item must have a short description and price")
		}
	}
	return nil
}
