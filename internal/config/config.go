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
	DatabaseDSN     string
	Secret          []byte
}

func NewConfig() (*appConfig, error) {
	serverAddress := getServerAddress()
	baseURL := getBaseURL()
	fileStoragePath := getFileStoragePath()
	databaseDSN := getDatabaseDSN()
	secret := getSecret()
	flag.Parse()

	if serverAddress == nil {
		return nil, errors.New("server address not specified")
	}

	if baseURL == nil {
		return nil, errors.New("base url not specified")
	}

	if fileStoragePath == nil {
		return nil, errors.New("file storage path not specified")
	}

	if databaseDSN == nil {
		return nil, errors.New("database dsn not specified")
	}

	if secret == nil {
		return nil, errors.New("secret key not specified")
	}

	return &appConfig{
		ServerAddress:   *serverAddress,
		BaseURL:         *baseURL,
		FileStoragePath: *fileStoragePath,
		DatabaseDSN:     *databaseDSN,
		Secret:          []byte(*secret),
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

func getDatabaseDSN() *string {
	dsn := os.Getenv("DATABASE_DSN")

	return flag.String("d", dsn, "data source name")
}

func getSecret() *string {
	url := os.Getenv("SECRET")
	if url == "" {
		url = "my-secret-key"
	}

	return flag.String("s", url, "secret")
}
