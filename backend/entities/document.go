package entities

import "time"

type Document struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	FilePath   string    `json:"file_path"`
	CategoryID *int      `json:"category_id"`
	UserID     int       `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateDocumentRequest struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	CategoryID *int   `json:"category_id"`
}

type UpdateDocumentRequest struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	CategoryID *int   `json:"category_id"`
}
