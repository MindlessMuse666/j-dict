package handler

import (
	"jp-ru-dict/backend/internal/model"
	"jp-ru-dict/backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WordsHandler struct {
	wordsService service.WordsService
}

// NewWordsHandler создает новый обработчик слов
func NewWordsHandler(wordsService service.WordsService) *WordsHandler {
	return &WordsHandler{wordsService: wordsService}
}

// CreateWord создает новое слово
// @Summary Создать слово
// @Description Создает новое слово в словаре пользователя. Обязательны поля Jp (японский) и Ru (русский).
// @Tags words
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.WordCreateRequest true "Данные слова"
// @Success 201 {object} model.WordResponseWrapper "Слово создано"
// @Failure 400 {object} model.ErrorResponse "Неверные данные"
// @Failure 500 {object} model.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/words [post]
func (h *WordsHandler) CreateWord(c *gin.Context) {
	userID := c.GetInt("user_id")
	var req model.WordCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "неверные данные: " + err.Error()})
		return
	}

	word, err := h.wordsService.CreateWord(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, model.WordResponseWrapper{Data: *word})
}

// GetWords получает список слов с пагинацией
// @Summary Получить список слов
// @Description Возвращает список слов пользователя с курсорной пагинацией. Поддерживает фильтрацию и сортировку.
// @Tags words
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Лимит слов (макс. 100)" default(20) minimum(1) maximum(100)
// @Param cursor query int false "Курсор для пагинации (ID последнего загруженного слова)" default(0) minimum(0)
// @Success 200 {object} model.WordsListResponseWrapper "Список слов"
// @Failure 500 {object} model.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/words [get]
func (h *WordsHandler) GetWords(c *gin.Context) {
	userID := c.GetInt("user_id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	cursor, _ := strconv.Atoi(c.DefaultQuery("cursor", "0"))

	if limit > 100 {
		limit = 100
	}

	words, err := h.wordsService.GetWords(userID, limit, cursor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	totalCount, err := h.wordsService.CountWords(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	hasMore := len(words) == limit
	nextCursor := 0
	if hasMore && len(words) > 0 {
		nextCursor = words[len(words)-1].ID
	}

	response := model.WordsListResponse{
		Words:      words,
		NextCursor: nextCursor,
		HasMore:    hasMore,
		TotalCount: totalCount,
	}

	c.JSON(http.StatusOK, model.WordsListResponseWrapper{Data: response})
}

// GetWord получает слово по ID
// @Summary Получить слово по ID
// @Description Возвращает слово по ID с проверкой владельца
// @Tags words
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID слова"
// @Success 200 {object} model.WordResponseWrapper "Слово"
// @Failure 404 {object} model.ErrorResponse "Слово не найдено"
// @Failure 400 {object} model.ErrorResponse "Неверный ID слова"
// @Router /api/words/{id} [get]
func (h *WordsHandler) GetWord(c *gin.Context) {
	userID := c.GetInt("user_id")
	wordID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID слова"})
		return
	}

	word, err := h.wordsService.GetWord(userID, wordID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if word == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "слово не найдено"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": word})
}

// UpdateWord обновляет слово
// @Summary Обновить слово
// @Description Обновляет существующее слово (частичное обновление).
// @Description Все поля опциональны: если поле не передано, оно сохраняет текущее значение.
// @Description Передача пустого массива [] очищает поле.
// @Description Передача null или отсутствие поля оставляет значение без изменений.
// @Tags words
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID слова"
// @Param request body model.WordUpdateRequest true "Обновленные данные (все поля опциональны)"
// @Success 200 {object} model.WordResponseWrapper "Слово обновлено"
// @Failure 400 {object} model.ErrorResponse "Неверные данные или пустой запрос"
// @Failure 404 {object} model.ErrorResponse "Слово не найдено"
// @Router /api/words/{id} [patch]
func (h *WordsHandler) UpdateWord(c *gin.Context) {
	userID := c.GetInt("user_id")
	wordID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID слова"})
		return
	}

	var req model.WordUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверные данные: " + err.Error()})
		return
	}

	word, err := h.wordsService.UpdateWord(userID, wordID, &req)
	if err != nil {
		// Проверяем тип ошибки для более точного HTTP статуса
		if err.Error() == "слово не найдено" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if err.Error() == "хотя бы одно поле должно быть передано для обновления" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": word})
}

// DeleteWord удаляет слово
// @Summary Удалить слово
// @Description Удаляет слово по ID с проверкой владельца
// @Tags words
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID слова"
// @Success 204 "Слово удалено"
// @Failure 404 {object} model.ErrorResponse "Слово не найдено"
// @Failure 400 {object} model.ErrorResponse "Неверный ID слова"
// @Router /api/words/{id} [delete]
func (h *WordsHandler) DeleteWord(c *gin.Context) {
	userID := c.GetInt("user_id")
	wordID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID слова"})
		return
	}

	if err := h.wordsService.DeleteWord(userID, wordID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// SearchWords выполняет поиск слов
// @Summary Поиск слов
// @Description Выполняет поиск слов по различным критериям
// @Tags words
// @Produce json
// @Security BearerAuth
// @Param q query string false "Поисковый запрос (подстрока в jp или ru)"
// @Param tags query []string false "Теги (точное совпадение)"
// @Param on query []string false "Онъёми (точное совпадение)"
// @Param kun query []string false "Кунъёми (точное совпадение)"
// @Param limit query int false "Лимит слов" default(20)
// @Param cursor query int false "Курсор для пагинации" default(0)
// @Success 200 {object} model.WordsListResponseWrapper "Результаты поиска"
// @Failure 500 {object} model.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/words/search [get]
func (h *WordsHandler) SearchWords(c *gin.Context) {
	userID := c.GetInt("user_id")
	query := c.Query("q")
	tags := c.QueryArray("tags")
	on := c.QueryArray("on")
	kun := c.QueryArray("kun")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	cursor, _ := strconv.Atoi(c.DefaultQuery("cursor", "0"))

	// Ограничиваем максимальный лимит
	if limit > 100 {
		limit = 100
	}

	words, err := h.wordsService.SearchWords(userID, query, tags, on, kun, limit, cursor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalCount, err := h.wordsService.CountSearchWords(userID, query, tags, on, kun)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Определяем есть ли еще слова
	hasMore := len(words) == limit
	nextCursor := 0
	if hasMore && len(words) > 0 {
		nextCursor = words[len(words)-1].ID
	}

	response := model.WordsListResponse{
		Words:      words,
		NextCursor: nextCursor,
		HasMore:    hasMore,
		TotalCount: totalCount,
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}
