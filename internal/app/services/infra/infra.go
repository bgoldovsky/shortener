//go:generate mockgen -source=infra.go -destination=mocks/mocks.go
package infra

import (
	"context"

	"github.com/sirupsen/logrus"
)

type urlsRepository interface {
	Ping(ctx context.Context) error
}

type service struct {
	urlsRepo urlsRepository
}

func NewService(urlsRepo urlsRepository) *service {
	return &service{
		urlsRepo: urlsRepo,
	}
}

// Ping Проверяет доступность базы данных
func (s *service) Ping(ctx context.Context) bool {
	err := s.urlsRepo.Ping(ctx)
	if err != nil {
		logrus.WithError(err).Error("ping database error")
		return false
	}

	return true
}
