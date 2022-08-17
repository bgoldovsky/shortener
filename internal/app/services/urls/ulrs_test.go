package urls

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/bgoldovsky/shortener/internal/app/models"
	internalErrors "github.com/bgoldovsky/shortener/internal/app/repositories/urls/errors"
	mockUrls "github.com/bgoldovsky/shortener/internal/app/services/urls/mocks"
)

const (
	host          = "http://localhost:8080"
	defaultUserID = "qwerty"
)

func TestService_Shorten(t *testing.T) {
	tests := []struct {
		name     string
		urlID    string
		shortcut string
		url      string
		err      error
	}{
		{
			name:     "success",
			urlID:    "qwerty",
			url:      "avito.ru",
			shortcut: "http://localhost:8080/qwerty",
		},
		{
			name:     "error",
			urlID:    "qwerty",
			url:      "avito.ru",
			shortcut: "",
			err:      errors.New("test err"),
		},
	}

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		genMock := mockUrls.NewMockgenerator(ctrl)
		genMock.EXPECT().RandomString(idLength).Return(tt.urlID, nil)

		repoMock := mockUrls.NewMockurlsRepository(ctrl)
		repoMock.EXPECT().Add(ctx, tt.urlID, tt.url, defaultUserID).Return(tt.err)

		s := NewService(repoMock, genMock, host)
		act, err := s.Shorten(ctx, tt.url, defaultUserID)

		assert.Equal(t, tt.err, err)
		assert.Equal(t, tt.shortcut, act)
	}
}

func TestService_Shorten_NotUniqueErr(t *testing.T) {
	tests := []struct {
		name   string
		urlID  string
		url    string
		err    error
		expErr error
		expURL string
	}{
		{
			name:   "unique error",
			urlID:  "qwerty",
			url:    "avito.ru",
			err:    internalErrors.NewNotUniqueURLErr("qwerty", "avito.ru", errors.New("test err")),
			expErr: ErrNotUniqueURL,
			expURL: "http://localhost:8080/qwerty",
		},
	}

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		genMock := mockUrls.NewMockgenerator(ctrl)
		genMock.EXPECT().RandomString(idLength).Return(tt.urlID, nil)

		repoMock := mockUrls.NewMockurlsRepository(ctrl)
		repoMock.EXPECT().Add(ctx, tt.urlID, tt.url, defaultUserID).Return(tt.err)

		s := NewService(repoMock, genMock, host)
		act, err := s.Shorten(ctx, tt.url, defaultUserID)

		assert.Equal(t, tt.expErr, err)
		assert.Equal(t, tt.expURL, act)
	}
}

func TestService_Expand(t *testing.T) {
	tests := []struct {
		name     string
		shortcut string
		url      string
		err      error
	}{
		{
			name:     "success",
			shortcut: "qwerty",
			url:      "avito.ru",
			err:      nil,
		},
		{
			name:     "repo err",
			shortcut: "qwerty",
			url:      "",
			err:      errors.New("test err"),
		},
	}

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		repoMock := mockUrls.NewMockurlsRepository(ctrl)
		repoMock.EXPECT().Get(ctx, tt.shortcut).Return(tt.url, tt.err)

		s := NewService(repoMock, nil, host)
		act, err := s.Expand(ctx, tt.shortcut)

		assert.Equal(t, tt.err, err)
		assert.Equal(t, tt.url, act)
	}
}

func TestService_GetUrls(t *testing.T) {
	tests := []struct {
		name string
		urls []models.URL
		err  error
	}{
		{
			name: "success",
			urls: []models.URL{
				{
					ShortURL:    "http://localhost:8080/xyz",
					OriginalURL: "https://avito.ru",
				},
				{
					ShortURL:    "http://localhost:8080/qwerty",
					OriginalURL: "https://yandex.ru",
				},
			},
			err: nil,
		},
		{
			name: "repo err",
			urls: nil,
			err:  errors.New("test err"),
		},
	}

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		repoMock := mockUrls.NewMockurlsRepository(ctrl)
		repoMock.EXPECT().GetList(ctx, defaultUserID).Return(tt.urls, tt.err)

		s := NewService(repoMock, nil, host)
		act, err := s.GetUrls(ctx, defaultUserID)

		assert.Equal(t, tt.err, err)
		assert.Equal(t, tt.urls, act)
	}
}

func TestService_ShortenBatch(t *testing.T) {
	tests := []struct {
		name         string
		originalURLs []models.OriginalURL
		urls         []models.URL
		err          error
		exp          []models.URL
	}{
		{
			name: "success",
			originalURLs: []models.OriginalURL{
				{
					CorrelationID: "1",
					URL:           "https://avito.ru",
				},
				{
					CorrelationID: "2",
					URL:           "https://yandex.ru",
				},
			},
			urls: []models.URL{
				{
					CorrelationID: "1",
					ShortURL:      "xyz",
					OriginalURL:   "https://avito.ru",
				},
				{
					CorrelationID: "2",
					ShortURL:      "qwerty",
					OriginalURL:   "https://yandex.ru",
				},
			},
			exp: []models.URL{
				{
					CorrelationID: "1",
					ShortURL:      "http://localhost:8080/xyz",
					OriginalURL:   "https://avito.ru",
				},
				{
					CorrelationID: "2",
					ShortURL:      "http://localhost:8080/qwerty",
					OriginalURL:   "https://yandex.ru",
				},
			},
			err: nil,
		},
		{
			name: "repo err",
			originalURLs: []models.OriginalURL{
				{
					CorrelationID: "1",
					URL:           "https://avito.ru",
				},
				{
					CorrelationID: "2",
					URL:           "https://yandex.ru",
				},
			},
			urls: []models.URL{
				{
					CorrelationID: "1",
					ShortURL:      "xyz",
					OriginalURL:   "https://avito.ru",
				},
				{
					CorrelationID: "2",
					ShortURL:      "qwerty",
					OriginalURL:   "https://yandex.ru",
				},
			},
			err: errors.New("test err"),
		},
	}

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		repoMock := mockUrls.NewMockurlsRepository(ctrl)
		repoMock.EXPECT().AddBatch(ctx, tt.urls, defaultUserID).Return(tt.err)

		genMock := mockUrls.NewMockgenerator(ctrl)
		for _, url := range tt.urls {
			genMock.EXPECT().RandomString(idLength).Return(url.ShortURL, nil)
		}

		s := NewService(repoMock, genMock, host)
		act, err := s.ShortenBatch(ctx, tt.originalURLs, defaultUserID)

		assert.Equal(t, tt.err, err)
		assert.Equal(t, tt.exp, act)
	}
}
