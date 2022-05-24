//go:generate mockgen -source=urls.go -destination=mocks/mocks.go
package urls

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type repo interface {
	Add(id, url string)
	Get(id string) (string, error)
}

type generator interface {
	ID() string
}

type service struct {
	repo      repo
	generator generator
	host      string
}

func NewService(repo repo, generator generator, host string) *service {
	return &service{
		repo:      repo,
		generator: generator,
		host:      host,
	}
}

// Shorten Сокращает URL
func (s *service) Shorten(url string) string {
	id := s.generator.ID()
	s.repo.Add(id, url)

	return fmt.Sprintf("%s/%s", s.host, id)
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
