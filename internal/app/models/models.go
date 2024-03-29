package models

type OriginalURL struct {
	CorrelationID string // Строковый идентификатор для пакетного запроса
	URL           string // Исходный URL
}

type URL struct {
	CorrelationID string // Строковый идентификатор для пакетного запроса
	ShortURL      string // Сокращенный URL
	OriginalURL   string // Исходный URL
}

type UserCollection struct {
	UserID string   // Идентификатор пользователя
	URLIDs []string // Идентификаторы URL пользователя
}
