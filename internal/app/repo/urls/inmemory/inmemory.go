package inmemory

import (
	"errors"
	"sync"

	"github.com/bgoldovsky/shortener/internal/app/models"
)

type inmemoryRepo struct {
	store map[string]string
	ma    sync.RWMutex
}

func NewRepo() *inmemoryRepo {
	return &inmemoryRepo{
		store: map[string]string{},
	}
}

// Add Сохраняет URL
func (r *inmemoryRepo) Add(id, url string) error {
	r.ma.Lock()
	defer r.ma.Unlock()

	r.store[id] = url
	return nil
}

// Get Возвращает URL
func (r *inmemoryRepo) Get(id string) (string, error) {
	r.ma.RLock()
	defer r.ma.RUnlock()

	url, ok := r.store[id]
	if !ok {
		return "", errors.New("url not found")
	}

	return url, nil
}

// GetList Возвращает список всех сокращенных URL
func (r *inmemoryRepo) GetList() ([]models.URL, error) {
	r.ma.RLock()
	defer r.ma.RUnlock()

	urls := make([]models.URL, 0, len(r.store))
	for shortURL, originalURL := range r.store {
		urls = append(urls, models.URL{
			ShortURL:    shortURL,
			OriginalURL: originalURL,
		})
	}

	return urls, nil
}
