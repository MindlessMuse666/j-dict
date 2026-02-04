package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

// NewHealthHandler создает новый обработчик проверки здоровья
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck проверяет работоспособность сервиса
// @Summary Проверка здоровья
// @Description Проверяет, что сервис работает.
// @Tags system
// @Produce json
// @Success 200 {object} map[string]string "Сервис работает"
// @Router /api/health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
