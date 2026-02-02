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

	_ = s.producer.SendEvent(kafka.EventWordUpdated, userID, word)

	return word, nil
}

// DeleteWord удаляет слово
func (s *wordsService) DeleteWord(userID, wordID int) error {
	return s.repo.DeleteWord(wordID, userID)
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

// ImportCSV импортирует слова из CSV строки
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

		// Преобразуем в WordCreateRequest
		wordReq := row.ParseCSVRow()

		// Валидация
		if len(wordReq.Jp) == 0 {
			response.FailedCount++
			response.Errors = append(response.Errors, "Поле 'jp' не может быть пустым")
			continue
		}

		if len(wordReq.Ru) == 0 {
			response.FailedCount++
			response.Errors = append(response.Errors, "Поле 'ru' не может быть пустым")
			continue
		}

		// Создаем слово
		_, err = s.CreateWord(userID, wordReq)
		if err != nil {
			response.FailedCount++
			response.Errors = append(response.Errors, fmt.Sprintf("Ошибка создания слова: %v", err))
			continue
		}

		response.ImportedCount++
	}

	return response, nil
}
