package server

import (
	"net/http"
	"strings"

	"my-notes-app/internal/handlers"
)

type Server struct {
	mux *http.ServeMux
}

func New(h *handlers.Handler) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"message":"pong"}`))
	})

	mux.HandleFunc("/users", h.RegisterUser)
	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/users/")
		if strings.Count(path, "/") == 1 {
			h.NotesHandler(w, r)
		} else {
			h.NoteByIDHandler(w, r)
		}
	})

	return &Server{mux: mux}
}

func (s *Server) Run(addr string) error {
	return http.ListenAndServe(addr, s.mux)
}
