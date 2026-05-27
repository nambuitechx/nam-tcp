package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/nambuitechx/nam-tcp/internal/models"
	"github.com/nambuitechx/nam-tcp/internal/services"
	"github.com/nambuitechx/nam-tcp/internal/utils"
)

type TargetHandler struct {
	Service services.TargetService
}

func NewTargetHandler(srv *services.TargetService) *TargetHandler {
	return  &TargetHandler{Service: *srv}
}

func (h *TargetHandler) GetTargets() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limit, offset := utils.ParseLimitAndOffset(r)
		if offset == -1 {
			http.Error(w, "invalid offset", http.StatusBadRequest)
			return
		}

		targets, err := h.Service.GetTargets(limit, offset)
		if err != nil {
			log.Printf("error in getting targets: %v", err)
			http.Error(w, "failed to get targets", http.StatusInternalServerError)
			return
		}

		utils.EncodeJsonResponse(w, "get targets successfully", targets)
	})
}

func (h *TargetHandler) CreateTarget() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload models.CreateTargetPayload
		err := json.NewDecoder(r.Body).Decode(&payload)

		if err != nil {
			http.Error(w, "Invalid JSON: " + err.Error(), http.StatusBadRequest)
			return
		}

		defer r.Body.Close()

		newTarget, err := h.Service.CreateTarget(&payload)
		if err != nil {
			log.Printf("error in creating target: %v", err)
			http.Error(w, "failed to create target", http.StatusInternalServerError)
			return
		}

		utils.EncodeJsonResponse(w, "create new target successfully", newTarget)
	})
}
