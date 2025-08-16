package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	_ "my-notes-app/internal/models"
	"my-notes-app/storage"
	"net/http"
	"strconv"
	"strings"
	_ "time"
)

func (h *Handler) NotesHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := strings.Split(r.URL.Path, "/")[2]
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPost:
		var req struct {
			Title   string `json:"title"`
			Content string `json:"content"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		note, err := h.Storage.CreateNote(r.Context(), userID, req.Title, req.Content)
		if err != nil {
			http.Error(w, "failed to create note", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusCreated, note)

	case http.MethodGet:
		notes, err := h.Storage.ListNotes(r.Context(), storage.ListNotesParams{
			UserID: userID,
			Limit:  50,
			Offset: 0,
			Sort:   "desc",
		})
		if err != nil {
			http.Error(w, "failed to get notes", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, notes)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) NoteByIDHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	userID, _ := strconv.ParseInt(parts[2], 10, 64)
	noteID, _ := strconv.ParseInt(parts[4], 10, 64)

	switch r.Method {
	case http.MethodGet:
		note, err := h.Storage.GetNote(r.Context(), userID, noteID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "note not found", http.StatusNotFound)
				return
			}
			http.Error(w, "failed to get note", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, note)

	case http.MethodPut:
		var req struct {
			Title   string `json:"title"`
			Content string `json:"content"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		note, err := h.Storage.UpdateNote(r.Context(), userID, noteID, req.Title, req.Content)
		if err != nil {
			http.Error(w, "failed to update note", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, note)

	case http.MethodDelete:
		err := h.Storage.DeleteNote(r.Context(), userID, noteID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "note not found", http.StatusNotFound)
				return
			}
			http.Error(w, "failed to delete note", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Println("failed to write json:", err)
	}
}
