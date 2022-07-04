package urls

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/bgoldovsky/shortener/internal/app/models"
	mockUrls "github.com/bgoldovsky/shortener/internal/app/services/urls/mocks"
)

const host = "http://localhost:8080"

func TestService_Shorten(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		shortcut string
		url      string
		err      error
	}{
		{
			name:     "success",
			id:       "qwerty",
			url:      "avito.ru",
			shortcut: "http://localhost:8080/qwerty",
		},
		{
			name:     "success",
			id:       "qwerty",
			url:      "avito.ru",
			shortcut: "",
			err:      errors.New("test err"),
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		genMock := mockUrls.NewMockgenerator(ctrl)
		genMock.EXPECT().ID().Return(tt.id)

		repoMock := mockUrls.NewMockurlRepo(ctrl)
		repoMock.EXPECT().Add(tt.id, tt.url).Return(tt.err)

		s := NewService(repoMock, genMock, host)
		act, err := s.Shorten(tt.url)

		assert.Equal(t, tt.err, err)
		assert.Equal(t, tt.shortcut, act)
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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		repoMock := mockUrls.NewMockurlRepo(ctrl)
		repoMock.EXPECT().Get(tt.shortcut).Return(tt.url, tt.err)

		s := NewService(repoMock, nil, host)
		act, err := s.Expand(tt.shortcut)

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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		repoMock := mockUrls.NewMockurlRepo(ctrl)
		repoMock.EXPECT().GetList().Return(tt.urls, tt.err)

		s := NewService(repoMock, nil, host)
		act, err := s.GetUrls()

		assert.Equal(t, tt.err, err)
		assert.Equal(t, tt.urls, act)
	}
}
