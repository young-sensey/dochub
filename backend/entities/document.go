package entities

import "time"

type Document struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Author     string    `json:"author"`
	FilePath   string    `json:"file_path"`
	CategoryID *int      `json:"category_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateDocumentRequest struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	Author     string `json:"author"`
	CategoryID *int   `json:"category_id"`
}

type UpdateDocumentRequest struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	Author     string `json:"author"`
	CategoryID *int   `json:"category_id"`
}
