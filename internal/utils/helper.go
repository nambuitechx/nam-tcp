package utils

import (
	"encoding/json"
	"net/http"
)

func EncodeJsonResponse(w http.ResponseWriter, message string, data any) {
	json.NewEncoder(w).Encode(map[string]any{
		"message": message,
		"data": data,
	})
}
