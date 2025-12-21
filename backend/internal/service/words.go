package service

import (
	"errors"
	"jp-ru-dict/backend/internal/model"
	"jp-ru-dict/backend/internal/repository"
	"jp-ru-dict/backend/internal/utils"
)

// WordsService определяет интерфейс сервиса работы со словами
type WordsService interface {
	CreateWord(userID int, req *model.WordCreateRequest) (*model.Word, error)
	GetWords(userID int, limit, cursor int) ([]*model.Word, error)
	GetWord(userID, wordID int) (*model.Word, error)
	UpdateWord(userID, wordID int, req *model.WordUpdateRequest) (*model.Word, error)
	DeleteWord(userID, wordID int) error
	SearchWords(userID int, query string, tags, on, kun []string, limit, cursor int) ([]*model.Word, error)
}

type wordsService struct {
	repo repository.WordsRepository
}

// NewWordsService создает новый экземпляр сервиса слов
func NewWordsService(repo repository.WordsRepository) WordsService {
	return &wordsService{repo: repo}
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
