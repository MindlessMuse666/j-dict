package model

import (
	"errors"
	"time"
)

// МОДЕЛИ

// Word представляет слово в словаре пользователя
// Содержит основную информацию о японском слове и его переводе
type Word struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Jp        []string  `json:"jp" db:"jp" validate:"required,min=1"`
	Ru        []string  `json:"ru" db:"ru" validate:"required,min=1"`
	On        []string  `json:"on,omitempty" db:"on"`
	Kun       []string  `json:"kun,omitempty" db:"kun"`
	ExJp      []string  `json:"ex_jp,omitempty" db:"ex_jp" validate:"max=3"`
	ExRu      []string  `json:"ex_ru,omitempty" db:"ex_ru" validate:"max=3"`
	Tags      []string  `json:"tags,omitempty" db:"tags" validate:"max=5"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// WordCreateRequest представляет запрос на создание слова
type WordCreateRequest struct {
	Jp   []string `json:"jp" binding:"required,min=1"`
	Ru   []string `json:"ru" binding:"required,min=1"`
	On   []string `json:"on,omitempty"`
	Kun  []string `json:"kun,omitempty"`
	ExJp []string `json:"ex_jp,omitempty" binding:"max=3"`
	ExRu []string `json:"ex_ru,omitempty" binding:"max=3"`
	Tags []string `json:"tags,omitempty" binding:"max=5"`
}

// WordUpdateRequest представляет запрос на обновление слова
// Все поля опциональны - если поле не передано, оно сохраняет текущее значение
type WordUpdateRequest struct {
	Jp   []string `json:"jp,omitempty" binding:"omitempty,min=1"`
	Ru   []string `json:"ru,omitempty" binding:"omitempty,min=1"`
	On   []string `json:"on,omitempty"`
	Kun  []string `json:"kun,omitempty"`
	ExJp []string `json:"ex_jp,omitempty" binding:"omitempty,max=3"`
	ExRu []string `json:"ex_ru,omitempty" binding:"omitempty,max=3"`
	Tags []string `json:"tags,omitempty" binding:"omitempty,max=5"`
}

// WordsListResponse представляет ответ со списком слов
type WordsListResponse struct {
	Words      []*Word `json:"words"`
	NextCursor int     `json:"next_cursor,omitempty"`
	HasMore    bool    `json:"has_more"`
	TotalCount int     `json:"total_count"`
}

// СЕРВИСНЫЕ МЕТОДЫ

// Validate проверяет, что хотя бы одно поле было передано для обновления
func (r *WordUpdateRequest) Validate() error {
	// Проверяем, что хотя бы одно поле было передано (не nil)
	// В Go JSON unmarshal устанавливает срезы в nil, если они не переданы
	// или были явно установлены в null
	if r.Jp == nil && r.Ru == nil && r.On == nil && r.Kun == nil &&
		r.ExJp == nil && r.ExRu == nil && r.Tags == nil {
		return errors.New("хотя бы одно поле должно быть передано для обновления")
	}
	return nil
}

// IsFieldPresent проверяет, было ли поле передано в запросе
func (r *WordUpdateRequest) IsFieldPresent(fieldName string) bool {
	switch fieldName {
	case "jp":
		return r.Jp != nil
	case "ru":
		return r.Ru != nil
	case "on":
		return r.On != nil
	case "kun":
		return r.Kun != nil
	case "ex_jp":
		return r.ExJp != nil
	case "ex_ru":
		return r.ExRu != nil
	case "tags":
		return r.Tags != nil
	default:
		return false
	}
}
