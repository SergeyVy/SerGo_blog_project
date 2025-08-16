package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if req.Username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	user, err := h.Storage.CreateUser(r.Context(), req.Username)
	if err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Println("error writing response:", err)
	}
}
