package inmemory

import (
	"context"
	"sync"

	"github.com/bgoldovsky/shortener/internal/app/models"
	internalErrors "github.com/bgoldovsky/shortener/internal/app/repositories/urls/errors"
)

type inmemoryRepository struct {
	store map[string]map[string]string
	ma    sync.RWMutex
}

func NewRepository() *inmemoryRepository {
	return &inmemoryRepository{
		store: map[string]map[string]string{},
	}
}

// Add Сохраняет URL
func (r *inmemoryRepository) Add(_ context.Context, urlID, url, userID string) error {
	r.ma.Lock()
	defer r.ma.Unlock()

	// Проверяем не содержится ли в репозитории такой URL
	if lastURLID, exist := r.urlExist(url); exist {
		return internalErrors.NewNotUniqueURLErr(lastURLID, url, nil)
	}

	// Извлекаем коллекцию URL пользователя из хранилища, если нет, то создаем новую
	userStore, ok := r.store[userID]
	if !ok {
		userStore = map[string]string{}
	}

	// Сохраняем коллекцию URL пользователя в хранилище
	userStore[urlID] = url
	r.store[userID] = userStore

	return nil
}

func (r *inmemoryRepository) urlExist(url string) (string, bool) {
	for _, userStore := range r.store {
		for urlID, originalURL := range userStore {
			if url == originalURL {
				return urlID, true
			}
		}
	}

	return "", false
}

func (r *inmemoryRepository) AddBatch(_ context.Context, urls []models.URL, userID string) error {
	r.ma.Lock()
	defer r.ma.Unlock()

	// Извлекаем коллекцию URL пользователя из хранилища, если нет, то создаем новую
	userStore, ok := r.store[userID]
	if !ok {
		userStore = map[string]string{}
	}

	// Добавляем URL в коллекцию пользователя, избегая копирования
	for idx := range urls {
		userStore[urls[idx].ShortURL] = urls[idx].OriginalURL
	}

	// Сохраняем коллекцию URL пользователя в хранилище
	r.store[userID] = userStore

	return nil
}

// Get Возвращает URL
func (r *inmemoryRepository) Get(_ context.Context, urlID string) (string, error) {
	r.ma.RLock()
	defer r.ma.RUnlock()

	for _, userStore := range r.store {
		if url, ok := userStore[urlID]; ok {
			return url, nil
		}
	}

	return "", internalErrors.ErrURLNotFound
}

// GetList Возвращает список всех сокращенных URL
func (r *inmemoryRepository) GetList(_ context.Context, userID string) ([]models.URL, error) {
	r.ma.RLock()
	defer r.ma.RUnlock()

	urls := make([]models.URL, 0)

	userStore, ok := r.store[userID]
	if !ok {
		return urls, nil
	}

	for shortURL, originalURL := range userStore {
		urls = append(urls, models.URL{
			ShortURL:    shortURL,
			OriginalURL: originalURL,
		})
	}

	return urls, nil
}

// Ping Проверяет доступность базы данных
func (r *inmemoryRepository) Ping(_ context.Context) error {
	return nil
}

func (r *inmemoryRepository) Close() error {
	return nil
}
