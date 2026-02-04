package service

import (
	"errors"
	"jp-ru-dict/backend/internal/model"
	"jp-ru-dict/backend/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// AuthService определяет интерфейс сервиса аутентификации
type AuthService interface {
	Register(userReq *model.UserRegisterRequest) (*model.AuthResponse, error)
	Login(userReq *model.UserLoginRequest) (*model.AuthResponse, error)
	GenerateToken(user *model.User) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
	GetUserFromToken(token *jwt.Token) (*model.User, error)
	UpdateAvatar(userID int, avatarURL string) error
}

type authService struct {
	repo      repository.AuthRepository
	jwtSecret string
	jwtExpiry time.Duration
}

// NewAuthService создает новый экземпляр сервиса аутентификации
func NewAuthService(repo repository.AuthRepository, jwtSecret string, jwtExpiry time.Duration) AuthService {
	return &authService{
		repo:      repo,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
}

// Register регистрирует нового пользователя
func (s *authService) Register(userReq *model.UserRegisterRequest) (*model.AuthResponse, error) {
	// Проверяем, существует ли пользователь с таким email
	existingUser, err := s.repo.GetUserByEmail(userReq.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("пользователь с таким email уже существует")
	}

	// Хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Создаем пользователя
	user := &model.User{
		Email:        userReq.Email,
		PasswordHash: string(hashedPassword),
		Name:         userReq.Name,
		Role:         model.RoleUser,
		AvatarURL:    "/assets/default_avatars/~catishe~cat~.jpg",
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}

	// Генерируем JWT токен
	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

// Login выполняет вход пользователя
func (s *authService) Login(userReq *model.UserLoginRequest) (*model.AuthResponse, error) {
	// Получаем пользователя по email
	user, err := s.repo.GetUserByEmail(userReq.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("неверный email или пароль")
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(userReq.Password)); err != nil {
		return nil, errors.New("неверный email или пароль")
	}

	// Генерируем JWT токен
	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

// GenerateToken генерирует JWT токен для пользователя
func (s *authService) GenerateToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"exp":   time.Now().Add(s.jwtExpiry).Unix(),
		"iat":   time.Now().Unix(),
		"email": user.Email,
		"name":  user.Name,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// ValidateToken проверяет валидность JWT токена
func (s *authService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("неожиданный метод подписи")
		}
		return []byte(s.jwtSecret), nil
	})
}

// GetUserFromToken извлекает пользователя из JWT токена
func (s *authService) GetUserFromToken(token *jwt.Token) (*model.User, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("неверные claims токена")
	}

	userIDFloat, ok := claims["sub"].(float64)
	if !ok {
		return nil, errors.New("неверный ID пользователя в токене")
	}

	return s.repo.GetUserByID(int(userIDFloat))
}

// UpdateAvatar обновляет аватар пользователя
func (s *authService) UpdateAvatar(userID int, avatarURL string) error {
	return s.repo.UpdateUserAvatar(userID, avatarURL)
}
