package models

import "time"

// Book representa un libro del catálogo de la biblioteca.
type Book struct {
	ID               int       `db:"id" json:"id"`
	Title            string    `db:"title" json:"title"`
	Author           string    `db:"author" json:"author"`
	ISBN             string    `db:"isbn" json:"isbn"`
	Category         string    `db:"category" json:"category"`
	CoverURL         string    `db:"cover_url" json:"cover_url"`
	TotalCopies      int       `db:"total_copies" json:"total_copies"`
	AvailableCopies  int       `db:"available_copies" json:"available_copies"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
}

// BookRequest es el payload para crear/actualizar un libro.
type BookRequest struct {
	Title       string `json:"title" binding:"required" example:"Clean Code"`
	Author      string `json:"author" binding:"required" example:"Robert C. Martin"`
	ISBN        string `json:"isbn" example:"9780132350884"`
	Category    string `json:"category" example:"Tecnología"`
	CoverURL    string `json:"cover_url" example:"https://ejemplo.com/portada.jpg"`
	TotalCopies int    `json:"total_copies" binding:"required,min=1" example:"3"`
}
