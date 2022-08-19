//go:generate mockgen -source=cleaner.go -destination=mocks/mocks.go

package cleaner

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/bgoldovsky/shortener/internal/app/models"
)

type urlsRepository interface {
	Delete(ctx context.Context, urlsBatch []models.UserCollection) error
}

type service struct {
	urlsRepo urlsRepository
	deleteCh chan models.UserCollection
	doneCh   <-chan struct{}
}

func NewService(urlsRepo urlsRepository, deleteCh chan models.UserCollection, doneCh <-chan struct{}) *service {
	return &service{
		urlsRepo: urlsRepo,
		deleteCh: deleteCh,
		doneCh:   doneCh,
	}
}

func (s *service) Queue(urls models.UserCollection) {
	s.deleteCh <- urls
}

// Run Запускает асинхронное удаление
func (s *service) Run() {
	go func() {
		for {
			select {
			case collection := <-s.deleteCh:
				err := s.urlsRepo.Delete(context.Background(), []models.UserCollection{collection})
				if err != nil {
					logrus.WithError(err).
						WithField("userID", collection.UserID).
						WithField("URLIDs", collection.URLIDs).
						Error("delete urls error")
				}
			case <-s.doneCh:
				logrus.Info("worker done")
				return
			}
		}
	}()
}
