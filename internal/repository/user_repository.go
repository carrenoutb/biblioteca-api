package repository

import (
	"github.com/jmoiron/sqlx"

	"biblioteca-api/internal/models"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(u *models.User) error {
	query := `INSERT INTO users (name, email, password_hash, role) VALUES (?, ?, ?, ?)`
	result, err := r.db.Exec(query, u.Name, u.Email, u.PasswordHash, u.Role)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = int(id)
	return nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var u models.User
	query := `SELECT id, name, email, password_hash, role, created_at FROM users WHERE email = ?`
	if err := r.db.Get(&u, query, email); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindByID(id int) (*models.User, error) {
	var u models.User
	query := `SELECT id, name, email, password_hash, role, created_at FROM users WHERE id = ?`
	if err := r.db.Get(&u, query, id); err != nil {
		return nil, err
	}
	return &u, nil
}
