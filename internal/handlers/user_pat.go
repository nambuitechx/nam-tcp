package handlers

import (
	"encoding/json"
	"log"
	"net/http"

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
		limit, offset := utils.ParseLimitAndOffset(r)
		if offset == -1 {
			http.Error(w, "invalid offset", http.StatusBadRequest)
			return
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


func (h *UserPATHandler) RevokeUserPAT() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		err := h.Service.RevokeUserPAT(id)
		if err != nil {
			log.Printf("error in revoking user pat: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		utils.EncodeJsonResponse(w, "revoke user pat successfully", nil)
	})
}

func (h *UserPATHandler) RevokeExpiredUserPATs() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := h.Service.RevokeExpiredUserPATs()
		if err != nil {
			log.Printf("error in revoking expired user pats: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		utils.EncodeJsonResponse(w, "revoke expired user pats successfully", nil)
	})
}