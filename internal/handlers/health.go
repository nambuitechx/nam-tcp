package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "healthy"}); err != nil {
		log.Printf("encode health: %v", err)
	}
}