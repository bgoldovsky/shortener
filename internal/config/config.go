package config

import "os"

// ServerAddress Возвращает адрес HTTP сервера
func ServerAddress() string {
	address := os.Getenv("SERVER_ADDRESS")
	if address == "" {
		address = ":8080"
	}

	return address
}

// BaseURL Возвращает хост для генерации сокращенного URL
func BaseURL() string {
	url := os.Getenv("BASE_URL")
	if url == "" {
		url = "http://localhost:8080"
	}

	return url
}
