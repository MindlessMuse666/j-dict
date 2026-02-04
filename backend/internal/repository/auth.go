package repository

import (
	"database/sql"
	"jp-ru-dict/backend/internal/model"
)

// AuthRepository определяет интерфейс для работы с данными пользователей
type AuthRepository interface {
	CreateUser(user *model.User) error
	GetUserByEmail(email string) (*model.User, error)
	GetUserByID(id int) (*model.User, error)
	UpdateUserAvatar(userID int, avatarURL string) error
}

type authRepository struct {
	db *sql.DB
}

// NewAuthRepository создает новый экземпляр репозитория аутентификации
func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepository{db: db}
}

// CreateUser создает нового пользователя в базе данных
func (r *authRepository) CreateUser(user *model.User) error {
	query := `INSERT INTO users (email, password_hash, name, avatar_url) 
			  VALUES ($1, $2, $3, $4) 
			  RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(query, user.Email, user.PasswordHash, user.Name, user.AvatarURL).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	return err
}

// GetUserByEmail получает пользователя по email
func (r *authRepository) GetUserByEmail(email string) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id, email, password_hash, name, avatar_url, created_at, updated_at 
			  FROM users 
			  WHERE email = $1`

	err := r.db.QueryRow(query, email).
		Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.AvatarURL,
			&user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return user, err
}

// GetUserByID получает пользователя по ID
func (r *authRepository) GetUserByID(id int) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id, email, name, avatar_url, created_at, updated_at 
			  FROM users 
			  WHERE id = $1`

	err := r.db.QueryRow(query, id).
		Scan(&user.ID, &user.Email, &user.Name, &user.AvatarURL,
			&user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return user, err
}

// UpdateUserAvatar обновляет аватар пользователя
func (r *authRepository) UpdateUserAvatar(userID int, avatarURL string) error {
	query := `UPDATE users SET avatar_url = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(query, avatarURL, userID)
	return err
}
