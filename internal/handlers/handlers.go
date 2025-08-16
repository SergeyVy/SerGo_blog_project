package handlers

import "my-notes-app/storage"

type Handler struct {
	store  *storage.Storage
	secret string
}

func NewHandler(store *storage.Storage, secret string) *Handler {
	return &Handler{store: store, secret: secret}
}
