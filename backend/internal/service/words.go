package service

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"jp-ru-dict/backend/internal/kafka"
	"jp-ru-dict/backend/internal/model"
	"jp-ru-dict/backend/internal/repository"
	"jp-ru-dict/backend/internal/utils"
	"strings"
)

// WordsService определяет интерфейс сервиса работы со словами
type WordsService interface {
	CreateWord(userID int, req *model.WordCreateRequest) (*model.Word, error)
	GetWords(userID int, limit, cursor int) ([]*model.Word, error)
	GetWord(userID, wordID int) (*model.Word, error)
	UpdateWord(userID, wordID int, req *model.WordUpdateRequest) (*model.Word, error)
	DeleteWord(userID, wordID int) error
	SearchWords(userID int, query string, tags, on, kun []string, limit, cursor int) ([]*model.Word, error)
	CountWords(userID int) (int, error)
	CountSearchWords(userID int, query string, tags, on, kun []string) (int, error)
	ImportCSV(userID int, csvContent string) (*model.CSVImportResponse, error) // Новый метод
}

type wordsService struct {
	repo     repository.WordsRepository
	producer kafka.Producer
}

// NewWordsService создает новый экземпляр сервиса слов
func NewWordsService(repo repository.WordsRepository, producer kafka.Producer) WordsService {
	return &wordsService{
		repo:     repo,
		producer: producer,
	}
}

// CreateWord создает новое слово в словаре пользователя
func (s *wordsService) CreateWord(userID int, req *model.WordCreateRequest) (*model.Word, error) {
	word := &model.Word{
		UserID: userID,
		Jp:     req.Jp,
		Ru:     req.Ru,
		On:     req.On,
		Kun:    req.Kun,
		ExJp:   req.ExJp,
		ExRu:   req.ExRu,
		Tags:   req.Tags,
	}

	if err := s.repo.CreateWord(word); err != nil {
		return nil, err
	}

	return word, nil
}

// GetWords получает список слов пользователя с пагинацией
func (s *wordsService) GetWords(userID int, limit, cursor int) ([]*model.Word, error) {
	return s.repo.GetWordsByUserID(userID, limit, cursor)
}

// GetWord получает слово по ID с проверкой владельца
func (s *wordsService) GetWord(userID, wordID int) (*model.Word, error) {
	return s.repo.GetWordByIDAndUserID(wordID, userID)
}

// UpdateWord обновляет существующее слово
func (s *wordsService) UpdateWord(userID, wordID int, req *model.WordUpdateRequest) (*model.Word, error) {
	// Валидируем, что хотя бы одно поле было передано
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Получаем текущее слово
	word, err := s.repo.GetWordByIDAndUserID(wordID, userID)
	if err != nil {
		return nil, err
	}
	if word == nil {
		return nil, errors.New("слово не найдено")
	}

	// Обновляем только переданные поля
	// Для массивов nil означает "не передано", пустой массив [] означает "очистить поле"
	if req.Jp != nil {
		word.Jp = req.Jp
	}
	if req.Ru != nil {
		word.Ru = req.Ru
	}
	if req.On != nil {
		word.On = req.On
	}
	if req.Kun != nil {
		word.Kun = req.Kun
	}
	if req.ExJp != nil {
		word.ExJp = req.ExJp
	}
	if req.ExRu != nil {
		word.ExRu = req.ExRu
	}
	if req.Tags != nil {
		word.Tags = req.Tags
	}

	if err := s.repo.UpdateWord(word); err != nil {
		return nil, err
	}

	// Сохраняем историю
	go func() {
		history := &model.WordHistory{
			WordID:   word.ID,
			UserID:   userID,
			Action:   model.ActionUpdate,
			Snapshot: *word,
		}
		_ = s.repo.CreateHistory(history)
	}()

	_ = s.producer.SendEvent(kafka.EventWordUpdated, userID, word)

	return word, nil
}

// DeleteWord удаляет слово
func (s *wordsService) DeleteWord(userID, wordID int) error {
	// Получаем слово для истории и события
	word, err := s.repo.GetWordByIDAndUserID(wordID, userID)
	if err != nil {
		return err
	}
	if word == nil {
		return errors.New("слово не найдено")
	}

	if err := s.repo.DeleteWord(wordID, userID); err != nil {
		return err
	}

	// Сохраняем историю
	go func() {
		history := &model.WordHistory{
			WordID:   word.ID,
			UserID:   userID,
			Action:   model.ActionDelete,
			Snapshot: *word,
		}
		_ = s.repo.CreateHistory(history)
	}()

	_ = s.producer.SendEvent(kafka.EventWordDeleted, userID, word)

	return nil
}

// SearchWords выполняет поиск слов по различным критериям
func (s *wordsService) SearchWords(userID int, query string, tags, on, kun []string, limit, cursor int) ([]*model.Word, error) {
	// Нормализуем онъёми: собираем все варианты (катакана + хирагана)
	var normalizedOn []string
	for _, reading := range on {
		variants := utils.NormalizeOnReading(reading)
		normalizedOn = append(normalizedOn, variants...)
	}

	return s.repo.SearchWords(userID, query, tags, normalizedOn, kun, limit, cursor)
}

// CountWords возвращает общее количество слов пользователя
func (s *wordsService) CountWords(userID int) (int, error) {
	return s.repo.CountWords(userID)
}

// CountSearchWords возвращает количество слов, соответствующих критериям поиска
func (s *wordsService) CountSearchWords(userID int, query string, tags, on, kun []string) (int, error) {
	// Нормализуем онъёми
	var normalizedOn []string
	for _, reading := range on {
		variants := utils.NormalizeOnReading(reading)
		normalizedOn = append(normalizedOn, variants...)
	}
	return s.repo.CountSearchWords(userID, query, tags, normalizedOn, kun)
}

// ImportCSV импортирует слова из CSV
func (s *wordsService) ImportCSV(userID int, csvContent string) (*model.CSVImportResponse, error) {
	response := &model.CSVImportResponse{
		Errors: []string{},
	}

	// Создаем CSV reader
	reader := csv.NewReader(strings.NewReader(csvContent))
	reader.Comma = ',' // Используем запятую как разделитель
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	// Читаем заголовок (если есть)
	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения CSV: %v", err)
	}

	// Проверяем наличие обязательных колонок
	hasJp := false
	hasRu := false
	for _, header := range headers {
		if strings.ToLower(header) == "jp" {
			hasJp = true
		}
		if strings.ToLower(header) == "ru" {
			hasRu = true
		}
	}

	if !hasJp || !hasRu {
		return nil, errors.New("CSV должен содержать колонки 'jp' и 'ru'")
	}

	// Читаем строки
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			response.FailedCount++
			response.Errors = append(response.Errors, fmt.Sprintf("Ошибка чтения строки: %v", err))
			continue
		}

		// Парсим строку в структуру
		row := model.CSVRow{}
		for i, header := range headers {
			if i >= len(record) {
				break
			}

			header = strings.ToLower(strings.TrimSpace(header))
			value := strings.TrimSpace(record[i])

			switch header {
			case "jp":
				row.Jp = value
			case "ru":
				row.Ru = value
			case "on":
				row.On = value
			case "kun":
				row.Kun = value
			case "ex_jp":
				row.ExJp = value
			case "ex_ru":
				row.ExRu = value
			case "tags":
				row.Tags = value
			}
		}

		// Создаем запрос на создание слова
		wordReq := row.ParseCSVRow()

		// Валидация (дополнительно к CreateWord)
		// Проверка обязательных полей
		if len(wordReq.Jp) == 0 {
			response.FailedCount++
			response.Errors = append(response.Errors, "Отсутствует японское написание")
			continue
		}
		if len(wordReq.Ru) == 0 {
			response.FailedCount++
			response.Errors = append(response.Errors, "Отсутствует перевод")
			continue
		}

		// Валидация длины массивов (BR-006 / Requirement Specification)
		if len(wordReq.ExJp) > 3 {
			response.FailedCount++
			response.Errors = append(response.Errors, fmt.Sprintf("Превышен лимит примеров (яп): %d > 3", len(wordReq.ExJp)))
			continue
		}
		if len(wordReq.ExRu) > 3 {
			response.FailedCount++
			response.Errors = append(response.Errors, fmt.Sprintf("Превышен лимит примеров (рус): %d > 3", len(wordReq.ExRu)))
			continue
		}
		if len(wordReq.Tags) > 5 {
			response.FailedCount++
			response.Errors = append(response.Errors, fmt.Sprintf("Превышен лимит тегов: %d > 5", len(wordReq.Tags)))
			continue
		}

		// Создаем слово
		_, err = s.CreateWord(userID, wordReq)
		if err != nil {
			response.FailedCount++
			response.Errors = append(response.Errors, fmt.Sprintf("Ошибка создания слова: %v", err))
		} else {
			response.ImportedCount++
		}
	}

	return response, nil
}
