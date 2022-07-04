package handlers

import (
	"github.com/bgoldovsky/shortener/internal/app/models"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestToGetUrlsReply(t *testing.T) {
	tests := []struct {
		model []models.URL
		exp   []GetUrlsReply
	}{
		{
			model: []models.URL{
				{
					ShortURL:    "http://localhost:8080/xyz",
					OriginalURL: "https://avito.ru",
				},
				{
					ShortURL:    "http://localhost:8080/qwerty",
					OriginalURL: "https://yandex.ru",
				},
			},
			exp: []GetUrlsReply{
				{
					ShortURL:    "http://localhost:8080/xyz",
					OriginalURL: "https://avito.ru",
				},
				{
					ShortURL:    "http://localhost:8080/qwerty",
					OriginalURL: "https://yandex.ru",
				},
			},
		},
		{
			model: []models.URL{},
			exp:   []GetUrlsReply{},
		},
		{
			model: nil,
			exp:   []GetUrlsReply{},
		},
	}

	for _, tt := range tests {
		act := toGetUrlsReply(tt.model)

		assert.Equal(t, tt.exp, act)
	}
}

func TestToShortenBatchRequest(t *testing.T) {
	tests := []struct {
		model []ShortenBatchRequest
		exp   []models.OriginalURL
	}{
		{
			model: []ShortenBatchRequest{
				{
					CorrelationID: "1",
					OriginalURL:   "https://avito.ru",
				},
				{
					CorrelationID: "2",
					OriginalURL:   "https://yandex.ru",
				},
			},
			exp: []models.OriginalURL{
				{
					CorrelationID: "1",
					URL:           "https://avito.ru",
				},
				{
					CorrelationID: "2",
					URL:           "https://yandex.ru",
				},
			},
		},
		{
			model: []ShortenBatchRequest{},
			exp:   []models.OriginalURL{},
		},
		{
			model: nil,
			exp:   []models.OriginalURL{},
		},
	}

	for _, tt := range tests {
		act := toShortenBatchRequest(tt.model)

		assert.Equal(t, tt.exp, act)
	}
}

func TestToShortenBatchReply(t *testing.T) {
	tests := []struct {
		model []models.URL
		exp   []ShortenBatchReply
	}{
		{
			model: []models.URL{
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
			exp: []ShortenBatchReply{
				{
					CorrelationID: "1",
					ShortURL:      "http://localhost:8080/xyz",
				},
				{
					CorrelationID: "2",
					ShortURL:      "http://localhost:8080/qwerty",
				},
			},
		},
		{
			model: []models.URL{},
			exp:   []ShortenBatchReply{},
		},
		{
			model: nil,
			exp:   []ShortenBatchReply{},
		},
	}

	for _, tt := range tests {
		act := toShortenBatchReply(tt.model)

		assert.Equal(t, tt.exp, act)
	}
}
