package model

import "time"

// WordHistoryAction определяет тип действия над словом
type WordHistoryAction string

const (
	ActionCreate WordHistoryAction = "CREATE"
	ActionUpdate WordHistoryAction = "UPDATE"
	ActionDelete WordHistoryAction = "DELETE"
)

// WordHistory представляет запись в истории изменений слова
type WordHistory struct {
	ID        int               `json:"id" db:"id" example:"1"`
	WordID    int               `json:"word_id" db:"word_id" example:"10"`
	UserID    int               `json:"user_id" db:"user_id" example:"1"`
	Action    WordHistoryAction `json:"action" db:"action" example:"CREATE"`
	Snapshot  Word              `json:"snapshot" db:"snapshot"`
	CreatedAt time.Time         `json:"created_at" db:"created_at" example:"2026-02-04T00:00:00Z"`
}
