package model

import "strings"

// CSVImportRequest представляет запрос на импорт слов из CSV
type CSVImportRequest struct {
	Content string `json:"content" binding:"required" example:"jp,ru\n猫,кот"`
}

// CSVImportResponse представляет ответ на импорт CSV
type CSVImportResponse struct {
	ImportedCount int      `json:"imported_count" example:"10"`
	FailedCount   int      `json:"failed_count" example:"0"`
	Errors        []string `json:"errors,omitempty" example:"строка 1: неверный формат"`
}

// CSVRow представляет строку CSV файла
type CSVRow struct {
	Jp   string `csv:"jp" example:"猫"`
	Ru   string `csv:"ru" example:"кот"`
	On   string `csv:"on,omitempty" example:"ビョウ"`
	Kun  string `csv:"kun,omitempty" example:"ねこ"`
	ExJp string `csv:"ex_jp,omitempty" example:"猫が好き"`
	ExRu string `csv:"ex_ru,omitempty" example:"Мне нравятся кошки"`
	Tags string `csv:"tags,omitempty" example:"животные"`
}

// ParseCSVRow парсит CSV строку в структуру WordCreateRequest
func (r *CSVRow) ParseCSVRow() *WordCreateRequest {
	return &WordCreateRequest{
		Jp:   splitCSVField(r.Jp),
		Ru:   splitCSVField(r.Ru),
		On:   splitCSVField(r.On),
		Kun:  splitCSVField(r.Kun),
		ExJp: splitCSVField(r.ExJp),
		ExRu: splitCSVField(r.ExRu),
		Tags: splitCSVField(r.Tags),
	}
}

// splitCSVField разделяет строку CSV поля на массив строк
func splitCSVField(field string) []string {
	if field == "" {
		return []string{}
	}

	// Используем точку с запятой как разделитель
	// Убираем лишние пробелы
	var result []string
	// Простой сплит, можно улучшить для обработки кавычек
	parts := splitPreservingQuotes(field, ';')
	for _, part := range parts {
		trimmed := trimQuotes(strings.TrimSpace(part))
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// splitPreservingQuotes разделяет строку с сохранением кавычек
func splitPreservingQuotes(s string, sep rune) []string {
	var result []string
	var current strings.Builder
	inQuotes := false
	escape := false

	for _, r := range s {
		if escape {
			current.WriteRune(r)
			escape = false
			continue
		}

		if r == '\\' {
			escape = true
			continue
		}

		if r == '"' {
			inQuotes = !inQuotes
			current.WriteRune(r)
			continue
		}

		if r == sep && !inQuotes {
			result = append(result, current.String())
			current.Reset()
			continue
		}

		current.WriteRune(r)
	}

	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}

// trimQuotes убирает кавычки с начала и конца строки
func trimQuotes(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}
