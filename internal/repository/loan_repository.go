package repository

import (
	"time"

	"github.com/jmoiron/sqlx"

	"biblioteca-api/internal/models"
)

type LoanRepository struct {
	db *sqlx.DB
}

func NewLoanRepository(db *sqlx.DB) *LoanRepository {
	return &LoanRepository{db: db}
}

func (r *LoanRepository) Create(userID, bookID int, dueAt time.Time) (*models.Loan, error) {
	query := `INSERT INTO loans (user_id, book_id, due_at, status) VALUES (?, ?, ?, 'active')`
	result, err := r.db.Exec(query, userID, bookID, dueAt)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return r.FindByID(int(id))
}

func (r *LoanRepository) FindByID(id int) (*models.Loan, error) {
	var l models.Loan
	query := `SELECT * FROM loans WHERE id = ?`
	if err := r.db.Get(&l, query, id); err != nil {
		return nil, err
	}
	return &l, nil
}

// FindByUser devuelve los préstamos de un usuario, incluyendo el título y autor del libro.
func (r *LoanRepository) FindByUser(userID int) ([]models.LoanDetail, error) {
	query := `
		SELECT l.*, b.title AS book_title, b.author AS book_author
		FROM loans l
		JOIN books b ON b.id = l.book_id
		WHERE l.user_id = ?
		ORDER BY l.loaned_at DESC`

	var loans []models.LoanDetail
	if err := r.db.Select(&loans, query, userID); err != nil {
		return nil, err
	}
	return loans, nil
}

// FindAll devuelve todos los préstamos del sistema (uso administrativo).
func (r *LoanRepository) FindAll() ([]models.LoanDetail, error) {
	query := `
		SELECT l.*, b.title AS book_title, b.author AS book_author
		FROM loans l
		JOIN books b ON b.id = l.book_id
		ORDER BY l.loaned_at DESC`

	var loans []models.LoanDetail
	if err := r.db.Select(&loans, query); err != nil {
		return nil, err
	}
	return loans, nil
}

func (r *LoanRepository) MarkReturned(id int) error {
	query := `UPDATE loans SET status = 'returned', returned_at = NOW() WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

// HasActiveLoan indica si el usuario ya tiene un préstamo activo del mismo libro.
func (r *LoanRepository) HasActiveLoan(userID, bookID int) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM loans WHERE user_id = ? AND book_id = ? AND status = 'active'`
	if err := r.db.Get(&count, query, userID, bookID); err != nil {
		return false, err
	}
	return count > 0, nil
}
