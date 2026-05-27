package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/nambuitechx/nam-tcp/internal/models"
	"github.com/nambuitechx/nam-tcp/internal/services"
	"github.com/nambuitechx/nam-tcp/internal/utils"
)

type UserPATHandler struct {
	Service services.UserPATService
}

func NewUserPATHandler(srv *services.UserPATService) *UserPATHandler {
	return  &UserPATHandler{Service: *srv}
}

func (h *UserPATHandler) GetUserPATs() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")

		limit := -1
		offset := 0
		var err error

		if limitStr != "" {
			limit, err = strconv.Atoi(limitStr)
			if err != nil {
				http.Error(w, "invalid limit", http.StatusBadRequest)
				return
			}
		}

		if offsetStr != "" {
			offset, err = strconv.Atoi(offsetStr)
			if err != nil {
				http.Error(w, "invalid offset", http.StatusBadRequest)
				return
			}
		}

		user_pats, err := h.Service.GetUserPATs(limit, offset)
		if err != nil {
			log.Printf("error in getting user pats: %v", err)
			http.Error(w, "failed to get user pats", http.StatusInternalServerError)
			return
		}

		utils.EncodeJsonResponse(w, "get user pats successfully", user_pats)
	})
}

func (h *UserPATHandler) CreateUserPAT() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload models.CreateUserPATPayload
		err := json.NewDecoder(r.Body).Decode(&payload)

		if err != nil {
			http.Error(w, "Invalid JSON: " + err.Error(), http.StatusBadRequest)
			return
		}

		defer r.Body.Close()

		plaintext, newUserPAT, err := h.Service.CreateUserPAT(&payload)
		if err != nil {
			log.Printf("error in creating user pat: %v", err)
			http.Error(w, "failed to create user pat", http.StatusInternalServerError)
			return
		}

		utils.EncodeJsonResponse(w, "create new user pat successfully", map[string]any{"plaintext": plaintext, "user_pat": newUserPAT})
	})
}
