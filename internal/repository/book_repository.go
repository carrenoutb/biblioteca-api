package repository

import (
	"github.com/jmoiron/sqlx"

	"biblioteca-api/internal/models"
)

type BookRepository struct {
	db *sqlx.DB
}

func NewBookRepository(db *sqlx.DB) *BookRepository {
	return &BookRepository{db: db}
}

// FindAll devuelve los libros, opcionalmente filtrados por categoría y/o texto de búsqueda en el título.
func (r *BookRepository) FindAll(category, search string) ([]models.Book, error) {
	query := `SELECT * FROM books WHERE 1=1`
	args := []interface{}{}

	if category != "" {
		query += ` AND category = ?`
		args = append(args, category)
	}
	if search != "" {
		query += ` AND title LIKE ?`
		args = append(args, "%"+search+"%")
	}
	query += ` ORDER BY title ASC`

	var books []models.Book
	if err := r.db.Select(&books, query, args...); err != nil {
		return nil, err
	}
	return books, nil
}

func (r *BookRepository) FindByID(id int) (*models.Book, error) {
	var b models.Book
	query := `SELECT * FROM books WHERE id = ?`
	if err := r.db.Get(&b, query, id); err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BookRepository) Create(b *models.Book) error {
	query := `INSERT INTO books (title, author, isbn, category, cover_url, total_copies, available_copies)
	          VALUES (?, ?, ?, ?, ?, ?, ?)`
	result, err := r.db.Exec(query, b.Title, b.Author, b.ISBN, b.Category, b.CoverURL, b.TotalCopies, b.AvailableCopies)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	b.ID = int(id)
	return nil
}

func (r *BookRepository) Update(b *models.Book) error {
	query := `UPDATE books SET title = ?, author = ?, isbn = ?, category = ?, cover_url = ?, total_copies = ?
	          WHERE id = ?`
	_, err := r.db.Exec(query, b.Title, b.Author, b.ISBN, b.Category, b.CoverURL, b.TotalCopies, b.ID)
	return err
}

func (r *BookRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM books WHERE id = ?`, id)
	return err
}

// DecrementAvailable reduce en 1 la cantidad de copias disponibles (al prestar un libro).
func (r *BookRepository) DecrementAvailable(id int) error {
	_, err := r.db.Exec(`UPDATE books SET available_copies = available_copies - 1 WHERE id = ? AND available_copies > 0`, id)
	return err
}

// IncrementAvailable aumenta en 1 la cantidad de copias disponibles (al devolver un libro).
func (r *BookRepository) IncrementAvailable(id int) error {
	_, err := r.db.Exec(`UPDATE books SET available_copies = available_copies + 1 WHERE id = ?`, id)
	return err
}
