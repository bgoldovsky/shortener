package file

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

type fileRepo struct {
	store    map[string]string
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

func readLines(filePath string) (map[string]string, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	scanner := bufio.NewScanner(file)
	res := make(map[string]string)

	if ok := scanner.Scan(); !ok {
		return res, nil
	}

	err = json.Unmarshal([]byte(scanner.Text()), &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Add Сохраняет URL
func (r *fileRepo) Add(id, url string) error {
	r.ma.Lock()
	defer r.ma.Unlock()

	r.store[id] = url

	file, err := os.OpenFile(r.filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("open file error: %w", err)
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	data, err := json.Marshal(r.store)
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
func (r *fileRepo) Get(id string) (string, error) {
	r.ma.RLock()
	defer r.ma.RUnlock()

	url, ok := r.store[id]
	if !ok {
		return "", errors.New("url not found")
	}

	return url, nil
}
