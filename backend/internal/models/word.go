package models

import (
	"time"

	"github.com/lib/pq"
)

// Word представляет слово в словаре пользователя
// Содержит основную информацию о японском слове и его переводе
type Word struct {
	ID        int            `json:"id" db:"id"`
	Jp        pq.StringArray `json:"jp" db:"jp"`                 // Написания на японском
	Ru        pq.StringArray `json:"ru" db:"ru"`                 // Переводы на русский
	On        pq.StringArray `json:"on,omitempty" db:"on"`       // Онъёми тения
	Kun       pq.StringArray `json:"kun,omitempty" db:"kun"`     // Кунъёми чтения
	ExJP      pq.StringArray `json:"ex_jp,omitempty" db:"ex_jp"` // Примеры на японском
	ExRU      pq.StringArray `json:"ex_ru,omitempty" db:"ex_ru"` // Переводы примеров
	Tags      pq.StringArray `json:"tags,omitempty" db:"tags"`   // Теги для категоризации
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at"`
}

// WordRequest представляет запрос на создание или обновление слова
type WordRequest struct {
	Jp   []string `json:"jp" validate:"required"`
	Ru   []string `json:"ru" validate:"required"`
	On   []string `json:"on,omitempty"`
	Kun  []string `json:"kun,omitempty"`
	ExJP []string `json:"ex_jp,omitempty"`
	ExRU []string `json:"ex_ru,omitempty"`
	Tags []string `json:"tags,omitempty"`
}

// WordsResponse представляет ответ со списком слов
type WordsResponse struct {
	Data []Word `json:"data"`
	Meta Meta   `json:"meta"`
}

// WordResponse представляет ответ с одним словом
type WordResponse struct {
	Data Word `json:"data"`
}

// Meta содержит метаинформацию для пагинации
type Meta struct {
	Total      int    `json:"total"`
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
	NextCursor string `json:"next_cursor,omitempty"` // TODO: Добавить для курсорной пагинации
}

// ErrorResponse представляет стандартный ответ об ошибке
type ErrorResponse struct {
	Error string `json:"error"`
}
