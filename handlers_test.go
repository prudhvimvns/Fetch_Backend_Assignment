package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

// Test a normal, valid receipt submission
func TestProcessReceiptHandler_Valid(t *testing.T) {
	body := `{
		"retailer": "M&M Corner Market",
		"purchaseDate": "2022-03-20",
		"purchaseTime": "14:33",
		"items": [
			{ "shortDescription": "Gatorade", "price": "2.25" },
			{ "shortDescription": "Gatorade", "price": "2.25" }
		],
		"total": "9.00"
	}`

	req := httptest.NewRequest("POST", "/receipts/process", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	ProcessReceiptHandler(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", res.Code)
	}

	var data map[string]string
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if data["id"] == "" {
		t.Error("expected receipt ID in response")
	}
}

// Test a GET request for an existing receipt ID
func TestGetPointsHandler_ExistingID(t *testing.T) {
	receiptStore["mock-id"] = 109

	req := httptest.NewRequest("GET", "/receipts/mock-id/points", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "mock-id"})
	res := httptest.NewRecorder()

	GetPointsHandler(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", res.Code)
	}

	expected := `{"points":109}`
	actual := strings.TrimSpace(res.Body.String())
	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

// Should return 404 for unknown receipt ID
func TestGetPointsHandler_InvalidID(t *testing.T) {
	req := httptest.NewRequest("GET", "/receipts/fake-id/points", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "fake-id"})
	res := httptest.NewRecorder()

	GetPointsHandler(res, req)

	if res.Code != http.StatusNotFound {
		t.Errorf("expected 404 Not Found, got %d", res.Code)
	}
}

// Submit a receipt with missing fields (empty retailer and no items)
func TestProcessReceiptHandler_MissingFields(t *testing.T) {
	body := `{
		"retailer": "",
		"purchaseDate": "2022-03-20",
		"purchaseTime": "14:33",
		"items": [],
		"total": "10.00"
	}`

	req := httptest.NewRequest("POST", "/receipts/process", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	ProcessReceiptHandler(res, req)

	if res.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d", res.Code)
	}
}

// Handle invalid JSON payload gracefully
func TestProcessReceiptHandler_InvalidJSON(t *testing.T) {
	body := `{
		"retailer": "Test Store",
		"purchaseDate": "2022-03-20",
		"items": [` // incomplete JSON

	req := httptest.NewRequest("POST", "/receipts/process", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	ProcessReceiptHandler(res, req)

	if res.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request for malformed JSON, got %d", res.Code)
	}
}

// Ensure empty body is rejected
func TestProcessReceiptHandler_EmptyBody(t *testing.T) {
	req := httptest.NewRequest("POST", "/receipts/process", nil)
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	ProcessReceiptHandler(res, req)

	if res.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request for empty body, got %d", res.Code)
	}
}

// Receipt where item description should trigger %3 length rule
func TestProcessReceiptHandler_DescLengthMultipleOf3(t *testing.T) {
	body := `{
		"retailer": "Target",
		"purchaseDate": "2022-01-01",
		"purchaseTime": "14:01",
		"items": [
			{ "shortDescription": "Emils Cheese Pizza", "price": "12.25" },
			{ "shortDescription": "Klarbrunn 12-PK 12 FL OZ", "price": "12.00" }
		],
		"total": "24.25"
	}`

	req := httptest.NewRequest("POST", "/receipts/process", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	ProcessReceiptHandler(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", res.Code)
	}

	var data map[string]string
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if data["id"] == "" {
		t.Error("expected receipt ID in response")
	}
}

// Item with empty shortDescription field should be invalid
func TestProcessReceiptHandler_EmptyItemDescription(t *testing.T) {
	body := `{
		"retailer": "Test Store",
		"purchaseDate": "2022-01-01",
		"purchaseTime": "13:00",
		"items": [
			{ "shortDescription": "", "price": "5.00" }
		],
		"total": "5.00"
	}`

	req := httptest.NewRequest("POST", "/receipts/process", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	ProcessReceiptHandler(res, req)

	if res.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request for empty item description, got %d", res.Code)
	}
}
