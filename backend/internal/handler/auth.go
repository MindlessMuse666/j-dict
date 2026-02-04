package handler

import (
	"jp-ru-dict/backend/internal/model"
	"jp-ru-dict/backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler создает новый обработчик аутентификации
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register обрабатывает регистрацию нового пользователя
// @Summary Регистрация пользователя
// @Description Создает нового пользователя и возвращает JWT токен. Требует email, пароль и имя.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.UserRegisterRequest true "Данные для регистрации"
// @Success 201 {object} model.AuthResponseWrapper "Успешная регистрация"
// @Failure 400 {object} model.ErrorResponse "Неверные данные"
// @Failure 409 {object} model.ErrorResponse "Пользователь уже существует"
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req model.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "неверные данные: " + err.Error()})
		return
	}

	authResponse, err := h.authService.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, model.AuthResponseWrapper{Data: *authResponse})
}

// Login обрабатывает вход пользователя
// @Summary Вход пользователя
// @Description Аутентифицирует пользователя и возвращает JWT токен.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.UserLoginRequest true "Данные для входа"
// @Success 200 {object} model.AuthResponseWrapper "Успешный вход"
// @Failure 400 {object} model.ErrorResponse "Неверные данные"
// @Failure 401 {object} model.ErrorResponse "Неверные учетные данные"
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "неверные данные: " + err.Error()})
		return
	}

	authResponse, err := h.authService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.AuthResponseWrapper{Data: *authResponse})
}

// Me возвращает информацию о текущем пользователе
// @Summary Информация о текущем пользователе
// @Description Возвращает информацию о пользователе из JWT токена. Требует Bearer токен.
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.UserResponseWrapper "Информация о пользователе"
// @Failure 401 {object} model.ErrorResponse "Неавторизован"
// @Router /api/auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "пользователь не найден"})
		return
	}

	c.JSON(http.StatusOK, model.UserResponseWrapper{Data: *user.(*model.User)})
}

// Logout обрабатывает выход пользователя
// @Summary Выход пользователя
// @Description Валидирует выход пользователя (на клиенте удаляется токен)
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string "Успешный выход"
// @Router /api/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "выход выполнен успешно"})
}
