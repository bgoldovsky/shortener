package config

import "os"

// ServerAddress Возвращает адрес HTTP сервера
func ServerAddress() string {
	p := os.Getenv("SERVER_ADDRESS")
	if p == "" {
		p = ":8080"
	}

	return p
}

// BaseURL Возвращает хост для генерации сокращенного URL
func BaseURL() string {
	h := os.Getenv("BASE_URL")
	if h == "" {
		h = "http://localhost:8080"
	}

	return h
}
