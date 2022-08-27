//go:generate mockgen -source=cleaner.go -destination=mocks/mocks.go

package cleaner

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/bgoldovsky/shortener/internal/app/models"
)

const queueSize = 100

type urlsRepository interface {
	Delete(ctx context.Context, urlsBatch []models.UserCollection) error
}

type service struct {
	urlsRepo urlsRepository
	deleteCh chan models.UserCollection
	bufferCh chan models.UserCollection
	doneCh   <-chan struct{}
}

func NewService(urlsRepo urlsRepository, deleteCh chan models.UserCollection, doneCh <-chan struct{}) *service {
	return &service{
		urlsRepo: urlsRepo,
		deleteCh: deleteCh,
		doneCh:   doneCh,
		bufferCh: make(chan models.UserCollection, queueSize),
	}
}

func (s *service) Queue(urls models.UserCollection) {
	s.deleteCh <- urls
}

// Run Запускает асинхронное удаление
func (s *service) Run() {
	go func() {
		ticker := time.NewTicker(time.Millisecond * 100)

		for {
			select {
			case collection := <-s.deleteCh:
				s.bufferCh <- collection
			case <-ticker.C:
				if len(s.bufferCh) == 0 {
					continue
				}

				var batch []models.UserCollection
				for i := 0; i < len(s.bufferCh); i++ {
					batch = append(batch, <-s.bufferCh)
				}

				err := s.urlsRepo.Delete(context.Background(), batch)
				if err != nil {
					logrus.WithError(err).WithField("batch", batch).Error("delete urls error")
				}
			case <-s.doneCh:
				logrus.Info("worker done")
				return
			}
		}
	}()
}
