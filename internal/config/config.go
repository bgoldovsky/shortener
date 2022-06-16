package config

import (
	"flag"
	"os"
)

func init() {
	cfg.serverAddress = serverAddress()
	cfg.baseURL = baseURL()
	cfg.fileStoragePath = fileStoragePath()

	flag.Parse()
}

// cfg Конфигурация сервиса
var cfg struct {
	serverAddress   *string
	baseURL         *string
	fileStoragePath *string
}

// ServerAddress Возвращает адрес HTTP сервера
func ServerAddress() string {
	if cfg.serverAddress == nil {
		panic("server address not specified")
	}

	return *cfg.serverAddress
}

// BaseURL Возвращает хост для генерации сокращенного URL
func BaseURL() string {
	if cfg.baseURL == nil {
		panic("base url not specified")
	}

	return *cfg.baseURL
}

// FileStoragePath Возвращает путь к файлу хранилища
func FileStoragePath() string {
	if cfg.fileStoragePath == nil {
		panic("file storage path not specified")
	}

	return *cfg.fileStoragePath
}

func serverAddress() *string {
	address := os.Getenv("SERVER_ADDRESS")
	if address == "" {
		address = ":8080"
	}

	return flag.String("a", address, "server address")
}

func baseURL() *string {
	url := os.Getenv("BASE_URL")
	if url == "" {
		url = "http://localhost:8080"
	}

	return flag.String("b", url, "base url")
}

func fileStoragePath() *string {
	path := os.Getenv("FILE_STORAGE_PATH")

	return flag.String("f", path, "file storage path")
}
