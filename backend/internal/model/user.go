package model

import (
	"time"
)

// Role определяет роль пользователя
type Role string

const (
	// RoleUser обычный пользователь
	RoleUser Role = "user"
	// RoleAdmin администратор
	RoleAdmin Role = "admin"
)

// User представляет пользователя системы
type User struct {
	ID           int       `json:"id" db:"id" example:"1"`
	Email        string    `json:"email" db:"email" example:"user@example.com"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Name         string    `json:"name" db:"name" example:"Иван Иванов"`
	Role         Role      `json:"role" db:"role" example:"user"`
	AvatarURL    string    `json:"avatar_url" db:"avatar_url" example:"https://example.com/avatar.jpg"`
	CreatedAt    time.Time `json:"created_at" db:"created_at" example:"2026-02-04T00:00:00Z"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at" example:"2026-02-04T00:00:00Z"`
}

// UserRegisterRequest представляет запрос на регистрацию
type UserRegisterRequest struct {
	Email    string `json:"email" binding:"required,email" example:"newuser@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"password"`
	Name     string `json:"name" binding:"required" example:"Егор Егоров"`
}

// UserLoginRequest представляет запрос на вход
type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"password"`
}

// AuthResponse представляет ответ с токеном аутентификации
type AuthResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  *User  `json:"user"`
}
