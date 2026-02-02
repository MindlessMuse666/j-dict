package model

// ResponseWrappers для Swagger документации
// Эти структуры нужны для того, чтобы Swagger корректно отображал структуру ответов,
// так как все успешные ответы обернуты в поле "data"

type AuthResponseWrapper struct {
	Data AuthResponse `json:"data"`
}

type UserResponseWrapper struct {
	Data User `json:"data"`
}

type WordResponseWrapper struct {
	Data Word `json:"data"`
}

type WordsListResponseWrapper struct {
	Data WordsListResponse `json:"data"`
}

type CSVImportResponseWrapper struct {
	Data CSVImportResponse `json:"data"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
