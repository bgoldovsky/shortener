package urls

import (
	"errors"
	"sync"
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
func (r *inmemoryRepo) Add(id, url string) {
	r.ma.Lock()
	defer r.ma.Unlock()

	r.store[id] = url
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
