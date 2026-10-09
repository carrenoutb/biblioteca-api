package models

import "time"

// Loan representa el préstamo de un libro a un usuario.
type Loan struct {
	ID         int        `db:"id" json:"id"`
	UserID     int        `db:"user_id" json:"user_id"`
	BookID     int        `db:"book_id" json:"book_id"`
	LoanedAt   time.Time  `db:"loaned_at" json:"loaned_at"`
	DueAt      time.Time  `db:"due_at" json:"due_at"`
	ReturnedAt *time.Time `db:"returned_at" json:"returned_at,omitempty"`
	Status     string     `db:"status" json:"status"`
}

// LoanRequest es el payload para pedir prestado un libro.
type LoanRequest struct {
	BookID int `json:"book_id" binding:"required" example:"1"`
}

// LoanDetail combina el préstamo con la información básica del libro,
// pensado para que el frontend no tenga que hacer una segunda petición.
type LoanDetail struct {
	Loan
	BookTitle  string `db:"book_title" json:"book_title"`
	BookAuthor string `db:"book_author" json:"book_author"`
}
