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
