package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"my-notes-app/internal/models"
	"strings"
	"time"
)

type Storage struct{ db *sql.DB }

func New(db *sql.DB) *Storage { return &Storage{db: db} }

// Миграция (создаст таблицы, если их нет)
func (s *Storage) Migrate(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS users(
  id BIGSERIAL PRIMARY KEY,
  username TEXT UNIQUE NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS notes(
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  content TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);`)
	return err
}

// Users
func (s *Storage) CreateUser(ctx context.Context, username string) (models.User, error) {
	var u models.User
	err := s.db.QueryRowContext(ctx,
		`INSERT INTO users(username) VALUES($1)
		 RETURNING id, username, created_at`, username).
		Scan(&u.ID, &u.Username, &u.CreatedAt)
	return u, err
}
func (s *Storage) GetUserByID(ctx context.Context, id int64) (models.User, error) {
	var u models.User
	err := s.db.QueryRowContext(ctx,
		`SELECT id, username, created_at FROM users WHERE id=$1`, id).
		Scan(&u.ID, &u.Username, &u.CreatedAt)
	return u, err
}

// Notes
func (s *Storage) CreateNote(ctx context.Context, userID int64, title, content string) (models.Note, error) {
	var n models.Note
	err := s.db.QueryRowContext(ctx,
		`INSERT INTO notes(user_id, title, content)
		 VALUES($1,$2,$3)
		 RETURNING id, user_id, title, content, created_at, updated_at`,
		userID, title, content).
		Scan(&n.ID, &n.UserID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt)
	return n, err
}

type ListNotesParams struct {
	UserID int64
	Limit  int
	Offset int
	Sort   string // asc/desc по created_at
}

func (s *Storage) ListNotes(ctx context.Context, p ListNotesParams) ([]models.Note, error) {
	if p.Limit <= 0 {
		p.Limit = 10
	}
	if strings.ToLower(p.Sort) != "asc" {
		p.Sort = "desc"
	}
	rows, err := s.db.QueryContext(ctx, fmt.Sprintf(`
SELECT id, user_id, title, content, created_at, updated_at
FROM notes
WHERE user_id=$1
ORDER BY created_at %s
LIMIT $2 OFFSET $3`, p.Sort), p.UserID, p.Limit, p.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Note
	for rows.Next() {
		var n models.Note
		if err := rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, rows.Err()
}

func (s *Storage) GetNote(ctx context.Context, userID, noteID int64) (models.Note, error) {
	var n models.Note
	err := s.db.QueryRowContext(ctx, `
SELECT id, user_id, title, content, created_at, updated_at
FROM notes WHERE id=$1 AND user_id=$2`, noteID, userID).
		Scan(&n.ID, &n.UserID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt)
	return n, err
}

func (s *Storage) UpdateNote(ctx context.Context, userID, noteID int64, title, content string) (models.Note, error) {
	var n models.Note
	err := s.db.QueryRowContext(ctx, `
UPDATE notes
SET title=$1, content=$2, updated_at=$3
WHERE id=$4 AND user_id=$5
RETURNING id, user_id, title, content, created_at, updated_at`,
		title, content, time.Now(), noteID, userID).
		Scan(&n.ID, &n.UserID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return n, sql.ErrNoRows
	}
	return n, err
}

func (s *Storage) DeleteNote(ctx context.Context, userID, noteID int64) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM notes WHERE id=$1 AND user_id=$2`, noteID, userID)
	if err != nil {
		return err
	}
	if aff, _ := res.RowsAffected(); aff == 0 {
		return sql.ErrNoRows
	}
	return nil
}
