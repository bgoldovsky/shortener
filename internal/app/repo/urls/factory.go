package urls

import (
	file "github.com/bgoldovsky/shortener/internal/app/repo/urls/file"
	inmemory "github.com/bgoldovsky/shortener/internal/app/repo/urls/inmemory"
)

type repo interface {
	Add(id, url string) error
	Get(id string) (string, error)
}

// Factory Инициализирует новый репозиторий
func Factory(filePath string) repo {
	if filePath != "" {
		r, err := file.NewRepo(filePath)
		if err != nil {
			panic(err)
		}

		return r
	}

	return inmemory.NewRepo()
}
