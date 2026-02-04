package handler

import (
	"fmt"
	"io"
	"jp-ru-dict/backend/internal/model"
	"jp-ru-dict/backend/internal/service"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	authService service.AuthService
}

// NewUserHandler создает новый обработчик пользователей
func NewUserHandler(authService service.AuthService) *UserHandler {
	return &UserHandler{authService: authService}
}

// UploadAvatar handles avatar upload
// @Summary Загрузка аватара
// @Description Загружает аватар пользователя (JPG/PNG) и обновляет профиль.
// @Tags users
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Файл аватара"
// @Success 200 {object} map[string]string "URL аватара"
// @Failure 400 {object} model.ErrorResponse "Ошибка загрузки или неверный формат"
// @Failure 401 {object} model.ErrorResponse "Неавторизован"
// @Failure 500 {object} model.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/users/avatar [post]
func (h *UserHandler) UploadAvatar(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Неавторизован"})
		return
	}
	userModel := user.(*model.User)

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Файл не загружен"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Разрешены только JPG и PNG"})
		return
	}

	cwd, _ := os.Getwd()
	uploadDir := filepath.Join(cwd, "uploads", "avatars")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: fmt.Sprintf("Не удалось создать директорию: %v", err)})
		return
	}

	filename := fmt.Sprintf("%d_%d%s", userModel.ID, time.Now().Unix(), ext)
	filePath := filepath.Join(uploadDir, filename)

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: fmt.Sprintf("Не удалось открыть файл: %v", err)})
		return
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: fmt.Sprintf("Не удалось создать файл: %v", err)})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: fmt.Sprintf("Не удалось сохранить файл: %v", err)})
		return
	}

	avatarURL := fmt.Sprintf("/uploads/avatars/%s", filename)

	if err := h.authService.UpdateAvatar(userModel.ID, avatarURL); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Не удалось обновить профиль"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": avatarURL})
}

// UpdateAvatarRequest запрос на обновление аватара по URL
type UpdateAvatarRequest struct {
	AvatarURL string `json:"avatar_url" binding:"required" example:"/assets/default_avatars/cat.jpg"`
}

// UpdateAvatarURL обновляет аватар по URL (для предустановленных)
// @Summary Обновление аватара по URL
// @Description Устанавливает аватар пользователя по указанному URL (обычно для предустановленных аватаров).
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateAvatarRequest true "URL аватара"
// @Success 200 {object} map[string]bool "Успех"
// @Failure 400 {object} model.ErrorResponse "Неверные данные"
// @Failure 401 {object} model.ErrorResponse "Неавторизован"
// @Failure 500 {object} model.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/users/avatar/url [put]
func (h *UserHandler) UpdateAvatarURL(c *gin.Context) {
	var req UpdateAvatarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Неавторизован"})
		return
	}
	userModel := user.(*model.User)

	if err := h.authService.UpdateAvatar(userModel.ID, req.AvatarURL); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Не удалось обновить профиль"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
