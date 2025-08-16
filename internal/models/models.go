package models

import "time"

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type Note struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Входящие DTO
type CreateUserRequest struct {
	Username string `json:"username"`
}

type CreateNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdateNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
