package model

import "strings"

// CSVImportRequest представляет запрос на импорт слов из CSV
type CSVImportRequest struct {
	Content string `json:"content" binding:"required"`
}

// CSVImportResponse представляет ответ на импорт CSV
type CSVImportResponse struct {
	ImportedCount int      `json:"imported_count"`
	FailedCount   int      `json:"failed_count"`
	Errors        []string `json:"errors,omitempty"`
}

// CSVRow представляет строку CSV файла
type CSVRow struct {
	Jp   string `csv:"jp"`
	Ru   string `csv:"ru"`
	On   string `csv:"on,omitempty"`
	Kun  string `csv:"kun,omitempty"`
	ExJp string `csv:"ex_jp,omitempty"`
	ExRu string `csv:"ex_ru,omitempty"`
	Tags string `csv:"tags,omitempty"`
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
