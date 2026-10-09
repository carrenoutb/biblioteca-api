package models

import "time"

// User representa un usuario del sistema (estudiante o administrador).
type User struct {
	ID           int       `db:"id" json:"id"`
	Name         string    `db:"name" json:"name"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Role         string    `db:"role" json:"role"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

// RegisterRequest es el payload esperado en POST /api/auth/register
type RegisterRequest struct {
	Name     string `json:"name" binding:"required" example:"Ana Pérez"`
	Email    string `json:"email" binding:"required,email" example:"ana@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"secreta123"`
}

// LoginRequest es el payload esperado en POST /api/auth/login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"ana@example.com"`
	Password string `json:"password" binding:"required" example:"secreta123"`
}

// AuthResponse es la respuesta de login/register exitoso.
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
