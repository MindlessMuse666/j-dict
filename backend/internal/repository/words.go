package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"jp-ru-dict/backend/internal/model"
	"strings"

	"github.com/lib/pq"
)

// WordsRepository определяет интерфейс для работы со словами
type WordsRepository interface {
	CreateWord(word *model.Word) error
	GetWordsByUserID(userID int, limit int, cursor int) ([]*model.Word, error)
	GetWordByIDAndUserID(wordID, userID int) (*model.Word, error)
	UpdateWord(word *model.Word) error
	DeleteWord(wordID, userID int) error
	SearchWords(userID int, query string, tags, on, kun []string, limit, cursor int) ([]*model.Word, error)
	CountWords(userID int) (int, error)
	CountSearchWords(userID int, query string, tags, on, kun []string) (int, error)
	CreateHistory(history *model.WordHistory) error
}

type wordsRepository struct {
	db *sql.DB
}

// NewWordsRepository создает новый экземпляр репозитория слов
func NewWordsRepository(db *sql.DB) WordsRepository {
	return &wordsRepository{db: db}
}

// CreateWord создает новое слово в словаре пользователя
func (r *wordsRepository) CreateWord(word *model.Word) error {
	query := `INSERT INTO words (user_id, jp, ru, "on", kun, ex_jp, ex_ru, tags)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
              RETURNING id, created_at, updated_at`

	// Используем pq.Array для конвертации []string в формат PostgreSQL массива
	err := r.db.QueryRow(query,
		word.UserID,
		pq.Array(word.Jp),
		pq.Array(word.Ru),
		pq.Array(word.On),
		pq.Array(word.Kun),
		pq.Array(word.ExJp),
		pq.Array(word.ExRu),
		pq.Array(word.Tags)).
		Scan(&word.ID, &word.CreatedAt, &word.UpdatedAt)

	return err
}

// GetWordsByUserID получает список слов пользователя с пагинацией
func (r *wordsRepository) GetWordsByUserID(userID int, limit int, cursor int) ([]*model.Word, error) {
	query := `SELECT id, user_id, jp, ru, "on", kun, ex_jp, ex_ru, tags, created_at, updated_at
              FROM words 
              WHERE user_id = $1 AND id > $2
              ORDER BY id ASC
              LIMIT $3`

	rows, err := r.db.Query(query, userID, cursor, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []*model.Word
	for rows.Next() {
		word := &model.Word{}
		err := rows.Scan(
			&word.ID,
			&word.UserID,
			pq.Array(&word.Jp),
			pq.Array(&word.Ru),
			pq.Array(&word.On),
			pq.Array(&word.Kun),
			pq.Array(&word.ExJp),
			pq.Array(&word.ExRu),
			pq.Array(&word.Tags),
			&word.CreatedAt,
			&word.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		words = append(words, word)
	}

	return words, nil
}

// GetWordByIDAndUserID получает слово по ID с проверкой владельца
func (r *wordsRepository) GetWordByIDAndUserID(wordID, userID int) (*model.Word, error) {
	word := &model.Word{}
	query := `SELECT id, user_id, jp, ru, "on", kun, ex_jp, ex_ru, tags, created_at, updated_at
              FROM words 
              WHERE id = $1 AND user_id = $2`

	err := r.db.QueryRow(query, wordID, userID).Scan(
		&word.ID,
		&word.UserID,
		pq.Array(&word.Jp),
		pq.Array(&word.Ru),
		pq.Array(&word.On),
		pq.Array(&word.Kun),
		pq.Array(&word.ExJp),
		pq.Array(&word.ExRu),
		pq.Array(&word.Tags),
		&word.CreatedAt,
		&word.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return word, err
}

// UpdateWord обновляет существующее слово
func (r *wordsRepository) UpdateWord(word *model.Word) error {
	query := `UPDATE words 
			  SET jp = $1, ru = $2, "on" = $3, kun = $4, ex_jp = $5, ex_ru = $6, tags = $7
			  WHERE id = $8 AND user_id = $9
			  RETURNING updated_at`

	// Конвертация для корректной обработки nil и пустых массивов
	err := r.db.QueryRow(query,
		pq.Array(word.Jp),
		pq.Array(word.Ru),
		pq.Array(word.On),
		pq.Array(word.Kun),
		pq.Array(word.ExJp),
		pq.Array(word.ExRu),
		pq.Array(word.Tags),
		word.ID, word.UserID).
		Scan(&word.UpdatedAt)

	return err
}

// DeleteWord удаляет слово
func (r *wordsRepository) DeleteWord(wordID, userID int) error {
	query := `DELETE FROM words WHERE id = $1 AND user_id = $2`
	_, err := r.db.Exec(query, wordID, userID)
	return err
}

// SearchWords выполняет поиск слов по различным критериям
func (r *wordsRepository) SearchWords(userID int, query string, tags, on, kun []string, limit, cursor int) ([]*model.Word, error) {
	baseQuery := `SELECT id, user_id, jp, ru, "on", kun, ex_jp, ex_ru, tags, created_at, updated_at
                  FROM words 
                  WHERE user_id = $1 AND id > $2`

	args := []any{userID, cursor}
	argCounter := 3

	var conditions []string

	// Поиск по тексту (подстрока в массивах jp или ru)
	if query != "" {
		conditions = append(conditions,
			fmt.Sprintf(`(EXISTS (SELECT 1 FROM unnest(jp) AS elem WHERE elem ILIKE '%%' || $%d || '%%') OR 
                          EXISTS (SELECT 1 FROM unnest(ru) AS elem WHERE elem ILIKE '%%' || $%d || '%%'))`,
				argCounter, argCounter))
		args = append(args, query)
		argCounter++
	}

	// Поиск по тегам (хотя бы один совпадающий элемент массива)
	if len(tags) > 0 {
		conditions = append(conditions, fmt.Sprintf(`tags && $%d`, argCounter))
		args = append(args, pq.Array(tags))
		argCounter++
	}

	// Поиск по онъёми
	if len(on) > 0 {
		conditions = append(conditions, fmt.Sprintf(`"on" && $%d`, argCounter))
		args = append(args, pq.Array(on))
		argCounter++
	}

	// Поиск по кунъёми
	if len(kun) > 0 {
		conditions = append(conditions, fmt.Sprintf(`kun && $%d`, argCounter))
		args = append(args, pq.Array(kun))
		argCounter++
	}

	// Добавляем условия к запросу
	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	}

	// Добавляем сортировку и лимит
	baseQuery += fmt.Sprintf(` ORDER BY id ASC LIMIT $%d`, argCounter)
	args = append(args, limit)

	rows, err := r.db.Query(baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []*model.Word
	for rows.Next() {
		word := &model.Word{}
		err := rows.Scan(
			&word.ID,
			&word.UserID,
			pq.Array(&word.Jp),
			pq.Array(&word.Ru),
			pq.Array(&word.On),
			pq.Array(&word.Kun),
			pq.Array(&word.ExJp),
			pq.Array(&word.ExRu),
			pq.Array(&word.Tags),
			&word.CreatedAt,
			&word.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		words = append(words, word)
	}

	return words, nil
}

// CountWords возвращает общее количество слов пользователя
func (r *wordsRepository) CountWords(userID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM words WHERE user_id = $1`
	err := r.db.QueryRow(query, userID).Scan(&count)
	return count, err
}

// CountSearchWords возвращает количество слов, соответствующих критериям поиска
func (r *wordsRepository) CountSearchWords(userID int, query string, tags, on, kun []string) (int, error) {
	baseQuery := `SELECT COUNT(*) FROM words WHERE user_id = $1`
	args := []any{userID}
	argCounter := 2
	var conditions []string

	if query != "" {
		conditions = append(conditions,
			fmt.Sprintf(`(EXISTS (SELECT 1 FROM unnest(jp) AS elem WHERE elem ILIKE '%%' || $%d || '%%') OR 
                          EXISTS (SELECT 1 FROM unnest(ru) AS elem WHERE elem ILIKE '%%' || $%d || '%%'))`,
				argCounter, argCounter))
		args = append(args, query)
		argCounter++
	}

	if len(tags) > 0 {
		conditions = append(conditions, fmt.Sprintf(`tags && $%d`, argCounter))
		args = append(args, pq.Array(tags))
		argCounter++
	}

	if len(on) > 0 {
		conditions = append(conditions, fmt.Sprintf(`"on" && $%d`, argCounter))
		args = append(args, pq.Array(on))
		argCounter++
	}

	if len(kun) > 0 {
		conditions = append(conditions, fmt.Sprintf(`"kun" && $%d`, argCounter))
		args = append(args, pq.Array(kun))
		argCounter++
	}

	fullQuery := baseQuery
	if len(conditions) > 0 {
		fullQuery += " AND " + strings.Join(conditions, " AND ")
	}

	var count int
	err := r.db.QueryRow(fullQuery, args...).Scan(&count)
	return count, err
}

// CreateHistory создает запись в истории изменений
func (r *wordsRepository) CreateHistory(history *model.WordHistory) error {
	query := `INSERT INTO word_history (word_id, user_id, action, snapshot)
			  VALUES ($1, $2, $3, $4)`

	// Сериализуем snapshot в JSON
	snapshotJSON, err := json.Marshal(history.Snapshot)
	if err != nil {
		return fmt.Errorf("ошибка маршалинга snapshot: %v", err)
	}

	_, err = r.db.Exec(query, history.WordID, history.UserID, history.Action, snapshotJSON)
	return err
}
