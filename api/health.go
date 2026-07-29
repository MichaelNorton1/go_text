package handler

import (
	"encoding/json"
	"net/http"
	"time"
)

type Health struct {
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := Health{Status: "OK", Message: "OK", Timestamp: time.Now().UTC()}
	json.NewEncoder(w).Encode(response)

}
