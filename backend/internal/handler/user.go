package handler

import (
	"fmt"
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

func NewUserHandler(authService service.AuthService) *UserHandler {
	return &UserHandler{authService: authService}
}

// UploadAvatar handles avatar upload
// @Summary Загрузка аватара
// @Description Загружает аватар пользователя и обновляет профиль
// @Tags users
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Avatar file"
// @Success 200 {object} map[string]string "Avatar URL"
// @Router /api/users/avatar [post]
func (h *UserHandler) UploadAvatar(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userModel := user.(*model.User)

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only JPG and PNG are allowed"})
		return
	}

	// Create directory if not exists
	uploadDir := "uploads/avatars"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	// Generate filename
	filename := fmt.Sprintf("%d_%d%s", userModel.ID, time.Now().Unix(), ext)
	filePath := filepath.Join(uploadDir, filename)

	// Save file
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Update user profile
	// URL should be /uploads/avatars/filename
	avatarURL := "/" + filePath // e.g. /uploads/avatars/123.jpg
	avatarURL = strings.ReplaceAll(avatarURL, "\\", "/")

	if err := h.authService.UpdateAvatar(userModel.ID, avatarURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": avatarURL})
}

type UpdateAvatarRequest struct {
	AvatarURL string `json:"avatar_url" binding:"required"`
}

// UpdateAvatarURL handles updating avatar by URL (presets)
func (h *UserHandler) UpdateAvatarURL(c *gin.Context) {
	var req UpdateAvatarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userModel := user.(*model.User)

	if err := h.authService.UpdateAvatar(userModel.ID, req.AvatarURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
