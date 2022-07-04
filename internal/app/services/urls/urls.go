//go:generate mockgen -source=urls.go -destination=mocks/mocks.go
package urls

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/bgoldovsky/shortener/internal/app/models"
)

type urlRepo interface {
	Add(id, url string) error
	Get(id string) (string, error)
	GetList() ([]models.URL, error)
}

type generator interface {
	ID() string
}

type service struct {
	repo      urlRepo
	generator generator
	host      string
}

func NewService(repo urlRepo, generator generator, host string) *service {
	return &service{
		repo:      repo,
		generator: generator,
		host:      host,
	}
}

// Shorten Сокращает URL
func (s *service) Shorten(url string) (string, error) {
	id := s.generator.ID()
	err := s.repo.Add(id, url)
	if err != nil {
		logrus.WithError(err).WithField("id", id).WithField("url", url).Error("add url error")
		return "", err
	}

	return fmt.Sprintf("%s/%s", s.host, id), nil
}

// Expand Возвращает полный URL по идентификатору сокращенного
func (s *service) Expand(id string) (string, error) {
	url, err := s.repo.Get(id)
	if err != nil {
		logrus.WithError(err).WithField("id", id).Error("get url error")
		return "", err
	}

	return url, nil
}

// GetUrls Возвращает список всех сокращенных URL
func (s *service) GetUrls() ([]models.URL, error) {
	urls, err := s.repo.GetList()
	if err != nil {
		logrus.WithError(err).Error("get url list error")
		return nil, err
	}

	return urls, nil
}
