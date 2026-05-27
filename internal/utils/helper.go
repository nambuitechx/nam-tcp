package utils

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func EncodeJsonResponse(w http.ResponseWriter, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"message": message,
		"data": data,
	})
}

func ParseLimitAndOffset(r *http.Request) (int, int) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := -1
	offset := 0
	var err error

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return -1, 0
		}
		if limit <= 0 {
			return -1, 0
		}
	}

	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			return -1, 0
		}
		if offset < 0 {
			return -1, 0
		}
	}

	return limit, offset
}
