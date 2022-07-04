package urls

import (
	"github.com/bgoldovsky/shortener/internal/app/models"
	"github.com/bgoldovsky/shortener/internal/app/repo/urls/file"
	"github.com/bgoldovsky/shortener/internal/app/repo/urls/inmemory"
)

type Repo interface {
	Add(id, url string) error
	Get(id string) (string, error)
	GetList() ([]models.URL, error)
}

// Factory Инициализирует новый репозиторий
func Factory(filePath string) Repo {
	if filePath != "" {
		r, err := file.NewRepo(filePath)
		if err != nil {
			panic(err)
		}

		return r
	}

	return inmemory.NewRepo()
}
