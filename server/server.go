package server //создаёт Gin Router, вешает middleware (логирование), регистрирует ручки.

import "net/http"

type Server struct {
	mux *http.ServeMux
}

func New() *Server {
	s := &Server{
		mux: http.NewServeMux(),
	}

	s.mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"pong"}`))
	})
	return s
}

func (s *Server) Run(addr string) error {
	return http.ListenAndServe(addr, s.mux)
}
