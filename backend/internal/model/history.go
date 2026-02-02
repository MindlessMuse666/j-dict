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
	ID        int               `json:"id" db:"id"`
	WordID    int               `json:"word_id" db:"word_id"`
	UserID    int               `json:"user_id" db:"user_id"`
	Action    WordHistoryAction `json:"action" db:"action"`
	Snapshot  Word              `json:"snapshot" db:"snapshot"` // Сохраняем полное состояние слова
	CreatedAt time.Time         `json:"created_at" db:"created_at"`
}
