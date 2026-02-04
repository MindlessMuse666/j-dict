package model

import (
	"errors"
	"time"
)

// МОДЕЛИ

// Word представляет слово в словаре пользователя
// Содержит основную информацию о японском слове и его переводе
type Word struct {
	ID        int       `json:"id" db:"id" example:"1"`
	UserID    int       `json:"user_id" db:"user_id" example:"1"`
	Jp        []string  `json:"jp" db:"jp" validate:"required,min=1" example:"猫"`
	Ru        []string  `json:"ru" db:"ru" validate:"required,min=1" example:"кот,кошка"`
	On        []string  `json:"on,omitempty" db:"on" example:"ビョウ"`
	Kun       []string  `json:"kun,omitempty" db:"kun" example:"ねこ"`
	ExJp      []string  `json:"ex_jp,omitempty" db:"ex_jp" validate:"max=3" example:"猫が好き"`
	ExRu      []string  `json:"ex_ru,omitempty" db:"ex_ru" validate:"max=3" example:"Мне нравятся кошки"`
	Tags      []string  `json:"tags,omitempty" db:"tags" validate:"max=5" example:"животные,n5"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" example:"2026-02-04T00:00:00Z"`
	CreatedAt time.Time `json:"created_at" db:"created_at" example:"2026-02-04T00:00:00Z"`
}

// WordCreateRequest представляет запрос на создание слова
type WordCreateRequest struct {
	Jp   []string `json:"jp" binding:"required,min=1" example:"猫"`
	Ru   []string `json:"ru" binding:"required,min=1" example:"кот,кошка"`
	On   []string `json:"on,omitempty" example:"ビョウ"`
	Kun  []string `json:"kun,omitempty" example:"ねこ"`
	ExJp []string `json:"ex_jp,omitempty" binding:"max=3" example:"猫が好き"`
	ExRu []string `json:"ex_ru,omitempty" binding:"max=3" example:"Мне нравятся кошки"`
	Tags []string `json:"tags,omitempty" binding:"max=5" example:"животные"`
}

// WordUpdateRequest представляет запрос на обновление слова
// Все поля опциональны - если поле не передано, оно сохраняет текущее значение
type WordUpdateRequest struct {
	Jp   []string `json:"jp,omitempty" binding:"omitempty,min=1" example:"猫"`
	Ru   []string `json:"ru,omitempty" binding:"omitempty,min=1" example:"кот,кошка"`
	On   []string `json:"on,omitempty" example:"ビョウ"`
	Kun  []string `json:"kun,omitempty" example:"ねこ"`
	ExJp []string `json:"ex_jp,omitempty" binding:"omitempty,max=3" example:"猫が好き"`
	ExRu []string `json:"ex_ru,omitempty" binding:"omitempty,max=3" example:"Мне нравятся кошки"`
	Tags []string `json:"tags,omitempty" binding:"omitempty,max=5" example:"животные"`
}

// WordsListResponse представляет ответ со списком слов
type WordsListResponse struct {
	Words      []*Word `json:"words"`
	NextCursor int     `json:"next_cursor,omitempty" example:"10"`
	HasMore    bool    `json:"has_more" example:"true"`
	TotalCount int     `json:"total_count" example:"100"`
}

// СЕРВИСНЫЕ МЕТОДЫ

// Validate проверяет, что хотя бы одно поле было передано для обновления
func (r *WordUpdateRequest) Validate() error {
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
