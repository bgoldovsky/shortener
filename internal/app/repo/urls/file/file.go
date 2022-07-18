package file

import (
	"bufio"
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"os"
	"sync"

	"github.com/bgoldovsky/shortener/internal/app/models"
	internalErrors "github.com/bgoldovsky/shortener/internal/app/repo/urls/errors"
)

type fileRepo struct {
	store    map[string]map[string]string
	ma       sync.RWMutex
	filePath string
}

// NewRepo Инициализирует репозиторий данными из файла
func NewRepo(filePath string) (*fileRepo, error) {
	store, err := readLines(filePath)
	if err != nil {
		return nil, fmt.Errorf("read urls from file error: %w", err)
	}

	return &fileRepo{
		store:    store,
		filePath: filePath,
	}, nil
}

func readLines(filePath string) (map[string]map[string]string, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	scanner := bufio.NewScanner(file)
	res := make(map[string]map[string]string)

	if ok := scanner.Scan(); !ok {
		return res, nil
	}

	res, err = unmarshal(scanner.Bytes())
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Add Сохраняет URL
func (r *fileRepo) Add(_ context.Context, urlID, url, userID string) error {
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

	return r.save()
}

func (r *fileRepo) urlExist(url string) (string, bool) {
	for _, userStore := range r.store {
		for urlID, originalURL := range userStore {
			if url == originalURL {
				return urlID, true
			}
		}
	}

	return "", false
}

func (r *fileRepo) AddBatch(_ context.Context, urls []models.URL, userID string) error {
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

	return r.save()
}

func (r *fileRepo) save() error {
	file, err := os.OpenFile(r.filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("open file error: %w", err)
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	data, err := marshal(r.store)
	if err != nil {
		return fmt.Errorf("serialize url error: %w", err)
	}

	_, err = file.WriteString(string(data))
	if err != nil {
		return fmt.Errorf("write url to file error: %w", err)
	}

	return nil
}

// Get Возвращает URL
func (r *fileRepo) Get(_ context.Context, urlID string) (string, error) {
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
func (r *fileRepo) GetList(_ context.Context, userID string) ([]models.URL, error) {
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
func (r *fileRepo) Ping(_ context.Context) error {
	return nil
}

func (r *fileRepo) Close() error {
	return nil
}

func marshal(store map[string]map[string]string) ([]byte, error) {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)

	err := encoder.Encode(store)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

func unmarshal(data []byte) (map[string]map[string]string, error) {
	store := map[string]map[string]string{}

	buff := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buff)

	err := decoder.Decode(&store)
	if err != nil {
		return nil, err
	}

	return store, nil
}
