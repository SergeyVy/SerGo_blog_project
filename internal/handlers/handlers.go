package handlers

import "my-notes-app/storage"

type Handler struct {
	Storage *storage.Storage
	secret  string
}

func NewHandler(store *storage.Storage, secret string) *Handler {
	return &Handler{
		Storage: store,
		secret:  secret,
	}
}
