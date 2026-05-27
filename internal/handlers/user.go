package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/nambuitechx/nam-tcp/internal/models"
	"github.com/nambuitechx/nam-tcp/internal/services"
	"github.com/nambuitechx/nam-tcp/internal/utils"
)

type UserHandler struct {
	Service services.UserService
}

func NewUserHandler(srv *services.UserService) *UserHandler {
	return  &UserHandler{Service: *srv}
}

func (h *UserHandler) GetUsers() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limit, offset := utils.ParseLimitAndOffset(r)
		if offset == -1 {
			http.Error(w, "invalid offset", http.StatusBadRequest)
			return
		}

		users, err := h.Service.GetUsers(limit, offset)
		if err != nil {
			log.Printf("error in getting users: %v", err)
			http.Error(w, "failed to get users", http.StatusInternalServerError)
			return
		}

		utils.EncodeJsonResponse(w, "get users successfully", users)
	})
}

func (h *UserHandler) CreateUser() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload models.CreateUserPayload
		err := json.NewDecoder(r.Body).Decode(&payload)

		if err != nil {
			http.Error(w, "Invalid JSON: " + err.Error(), http.StatusBadRequest)
			return
		}

		defer r.Body.Close()

		newUser, err := h.Service.CreateUser(&payload)
		if err != nil {
			log.Printf("error in creating user: %v", err)
			http.Error(w, "failed to create user", http.StatusInternalServerError)
			return
		}

		utils.EncodeJsonResponse(w, "create new user successfully", newUser)
	})
}
