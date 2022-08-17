package urls

import (
	"context"
	"fmt"

	"github.com/bgoldovsky/shortener/internal/app/models"
	"github.com/bgoldovsky/shortener/internal/app/repositories/urls/file"
	"github.com/bgoldovsky/shortener/internal/app/repositories/urls/inmemory"
	"github.com/bgoldovsky/shortener/internal/app/repositories/urls/postgres"
)

type Repository interface {
	Add(ctx context.Context, urlID, url, userID string) error
	AddBatch(ctx context.Context, urls []models.URL, userID string) error
	Get(ctx context.Context, urlID string) (string, error)
	GetList(ctx context.Context, userID string) ([]models.URL, error)
	Ping(ctx context.Context) error
	Close() error
}

// Factory Инициализирует новый репозиторий
func Factory(filePath, databaseDSN string) (Repository, error) {
	switch {
	case databaseDSN != "":
		r, err := postgres.NewRepository(databaseDSN)
		if err != nil {
			return nil, fmt.Errorf("initialize postgres repo error: %w", err)
		}
		return r, nil
	case filePath != "":
		r, err := file.NewRepository(filePath)
		if err != nil {
			return nil, fmt.Errorf("initialize file repo error: %w", err)
		}
		return r, nil
	}

	return inmemory.NewRepository(), nil
}
