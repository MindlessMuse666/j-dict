package handler

import (
	"jp-ru-dict/backend/internal/model"
	"jp-ru-dict/backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ImportHandler struct {
	wordsService service.WordsService
}

// NewImportHandler создает новый обработчик импорта
func NewImportHandler(wordsService service.WordsService) *ImportHandler {
	return &ImportHandler{wordsService: wordsService}
}

// ImportCSV обрабатывает импорт слов из CSV
// @Summary Импорт слов из CSV
// @Description Импортирует слова из CSV формата. Формат: jp;ru;on;kun;ex_jp;ex_ru;tags. Поля разделены точкой с запятой, массивы внутри полей разделены запятой.
// @Tags words
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.CSVImportRequest true "CSV данные для импорта"
// @Success 200 {object} model.CSVImportResponseWrapper "Результат импорта"
// @Failure 400 {object} model.ErrorResponse "Неверный формат CSV"
// @Failure 500 {object} model.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/words/import [post]
func (h *ImportHandler) ImportCSV(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req model.CSVImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "неверный формат запроса: " + err.Error()})
		return
	}

	if req.Content == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "CSV контент не может быть пустым"})
		return
	}

	result, err := h.wordsService.ImportCSV(userID, req.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "ошибка импорта: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.CSVImportResponseWrapper{Data: *result})
}
