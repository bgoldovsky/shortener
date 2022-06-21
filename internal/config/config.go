package config

import (
	"errors"
	"flag"
	"os"
)

type appConfig struct {
	ServerAddress   string
	BaseURL         string
	FileStoragePath string
}

func New() (*appConfig, error) {
	serverAddress := getServerAddress()
	baseURL := getBaseURL()
	fileStoreagePath := getFileStoragePath()
	flag.Parse()

	if serverAddress == nil {
		return nil, errors.New("server address not specified")
	}

	if baseURL == nil {
		return nil, errors.New("base url not specified")
	}

	if fileStoreagePath == nil {
		return nil, errors.New("file storage path not specified")
	}

	return &appConfig{
		ServerAddress:   *serverAddress,
		BaseURL:         *baseURL,
		FileStoragePath: *fileStoreagePath,
	}, nil
}

func getServerAddress() *string {
	address := os.Getenv("SERVER_ADDRESS")
	if address == "" {
		address = ":8080"
	}

	return flag.String("a", address, "server address")
}

func getBaseURL() *string {
	url := os.Getenv("BASE_URL")
	if url == "" {
		url = "http://localhost:8080"
	}

	return flag.String("b", url, "base url")
}

func getFileStoragePath() *string {
	path := os.Getenv("FILE_STORAGE_PATH")

	return flag.String("f", path, "file storage path")
}
