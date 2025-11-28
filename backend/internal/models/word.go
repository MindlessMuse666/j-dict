package models

import (
	"time"

	"github.com/lib/pq"
)

type Word struct {
	ID        int            `json:"id" db:"id"`
	Ru   	  string         `json:"ru" db:"ru"`
	Jp 		  string         `json:"jp" db:"jp"`
	On        pq.StringArray `json:"on,omitempty" db:"on"`       // онъёми
	Kun       pq.StringArray `json:"kun,omitempty" db:"kun"`     // кунъёми
	ExJP      pq.StringArray `json:"ex_jp,omitempty" db:"ex_jp"` // примеры на японском
	ExRU      pq.StringArray `json:"ex_ru,omitempty" db:"ex_ru"` // примеры на русском
	Tags      pq.StringArray `json:"tags,omitempty" db:"tags"`   // теги для категоризации
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at"`
}

type WordRequest struct {
	Ru 		 string   `json:"ru" validate:"required"`
	Jp 		 string   `json:"jp" validate:"required"`
	On       []string `json:"on,omitempty"`
	Kun      []string `json:"kun,omitempty"`
	ExJP     []string `json:"ex_jp,omitempty"`
	ExRU     []string `json:"ex_ru,omitempty"`
	Tags     []string `json:"tags,omitempty"`
}

type WordsResponse struct {
	Data []Word `json:"data"`
	Meta Meta   `json:"meta,omitempty"`
}

type WordResponse struct {
	Data Word `json:"data"`
}

type Meta struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
