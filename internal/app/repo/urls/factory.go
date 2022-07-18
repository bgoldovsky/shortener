package urls

import (
	"context"
	"fmt"

	"github.com/bgoldovsky/shortener/internal/app/models"
	"github.com/bgoldovsky/shortener/internal/app/repo/urls/file"
	"github.com/bgoldovsky/shortener/internal/app/repo/urls/inmemory"
	"github.com/bgoldovsky/shortener/internal/app/repo/urls/postgres"
)

type Repo interface {
	Add(ctx context.Context, urlID, url, userID string) error
	AddBatch(ctx context.Context, urls []models.URL, userID string) error
	Get(ctx context.Context, urlID string) (string, error)
	GetList(ctx context.Context, userID string) ([]models.URL, error)
	Ping(ctx context.Context) error
	Close() error
}

// Factory Инициализирует новый репозиторий
func Factory(filePath, databaseDSN string) (Repo, error) {
	switch {
	case databaseDSN != "":
		r, err := postgres.NewRepo(databaseDSN)
		if err != nil {
			return nil, fmt.Errorf("initialize postgres repo error: %w", err)
		}
		return r, nil
	case filePath != "":
		r, err := file.NewRepo(filePath)
		if err != nil {
			return nil, fmt.Errorf("initialize file repo error: %w", err)
		}
		return r, nil
	}

	return inmemory.NewRepo(), nil
}
