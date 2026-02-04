package model

import (
	"time"
)

// User представляет пользователя системы
type User struct {
	ID           int       `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Name         string    `json:"name" db:"name"`
	Role         string    `json:"role" db:"role"`
	AvatarURL    string    `json:"avatar_url" db:"avatar_url"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// UserRegisterRequest представляет запрос на регистрацию
type UserRegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

// UserLoginRequest представляет запрос на вход
type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse представляет ответ с токеном аутентификации
type AuthResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}
