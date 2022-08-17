//go:generate mockgen -source=urls.go -destination=mocks/mocks.go
package urls

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/bgoldovsky/shortener/internal/app/models"
	internalErrors "github.com/bgoldovsky/shortener/internal/app/repositories/urls/errors"
)

const idLength int64 = 5

var (
	ErrURLNotFound  = errors.New("url not found error")
	ErrURLDeleted   = errors.New("url has been deleted error")
	ErrNotUniqueURL = errors.New("url not unique error")
)

type urlsRepository interface {
	Add(ctx context.Context, urlID, url, userID string) error
	AddBatch(ctx context.Context, urls []models.URL, userID string) error
	Get(ctx context.Context, urlID string) (string, error)
	GetList(ctx context.Context, userID string) ([]models.URL, error)
}

type generator interface {
	RandomString(n int64) (string, error)
}

type service struct {
	urlsRepo  urlsRepository
	generator generator
	host      string
}

func NewService(urlsRepo urlsRepository, generator generator, host string) *service {
	return &service{
		urlsRepo:  urlsRepo,
		generator: generator,
		host:      host,
	}
}

// Shorten Сокращает URL
func (s *service) Shorten(ctx context.Context, url, userID string) (string, error) {
	urlID, err := s.generator.RandomString(idLength)
	if err != nil {
		logrus.WithError(err).
			WithField("userID", userID).
			WithField("url", url).
			Error("generate urlID error")
		return "", err
	}

	if err = s.urlsRepo.Add(ctx, urlID, url, userID); err != nil {
		var uniqueErr *internalErrors.NotUniqueURLErr
		if errors.As(err, &uniqueErr) {
			return s.buildShortURL(uniqueErr.URLID), ErrNotUniqueURL
		}

		logrus.WithError(err).
			WithField("userID", userID).
			WithField("urlID", urlID).
			WithField("url", url).
			Error("add url error")
		return "", err
	}

	return s.buildShortURL(urlID), nil
}

// ShortenBatch Сокращает несколько URL
func (s *service) ShortenBatch(ctx context.Context, originalURLs []models.OriginalURL, userID string) ([]models.URL, error) {
	urls := make([]models.URL, len(originalURLs))
	for idx := range urls {
		urlID, err := s.generator.RandomString(idLength)
		if err != nil {
			logrus.WithError(err).
				WithField("userID", userID).
				WithField("originalURLs", originalURLs).
				Error("generate urlID error")
			return nil, err
		}

		urls[idx] = models.URL{
			CorrelationID: originalURLs[idx].CorrelationID,
			ShortURL:      urlID,
			OriginalURL:   originalURLs[idx].URL,
		}
	}

	err := s.urlsRepo.AddBatch(ctx, urls, userID)
	if err != nil {
		logrus.WithError(err).
			WithField("userID", userID).
			WithField("originalURLs", originalURLs).
			WithField("urls", urls).
			Error("add urls batch error")
		return nil, err
	}

	for idx := range urls {
		urls[idx].ShortURL = s.buildShortURL(urls[idx].ShortURL)
	}

	return urls, nil
}

// Expand Возвращает полный URL по идентификатору сокращенного
func (s *service) Expand(ctx context.Context, urlID string) (string, error) {
	url, err := s.urlsRepo.Get(ctx, urlID)
	if err != nil {
		if errors.Is(err, internalErrors.ErrURLNotFound) {
			return "", ErrURLNotFound
		}

		if errors.Is(err, internalErrors.ErrURLDeleted) {
			return "", ErrURLDeleted
		}

		logrus.WithError(err).WithField("urlID", urlID).Error("get url error")
		return "", err
	}

	return url, nil
}

// GetUrls Возвращает список всех сокращенных URL
func (s *service) GetUrls(ctx context.Context, userID string) ([]models.URL, error) {
	urls, err := s.urlsRepo.GetList(ctx, userID)
	if err != nil {
		logrus.WithError(err).WithField("urlID", userID).Error("get url list error")
		return nil, err
	}

	for idx := range urls {
		urls[idx].ShortURL = s.buildShortURL(urls[idx].ShortURL)
	}

	return urls, nil
}

func (s *service) buildShortURL(id string) string {
	return fmt.Sprintf("%s/%s", s.host, id)
}
