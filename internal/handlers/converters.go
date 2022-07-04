package handlers

import "github.com/bgoldovsky/shortener/internal/app/models"

func toGetUrlsReply(model []models.URL) []GetUrlsReply {
	reply := make([]GetUrlsReply, len(model))

	for idx, m := range model {
		reply[idx] = GetUrlsReply{
			ShortURL:    m.ShortURL,
			OriginalURL: m.OriginalURL,
		}
	}

	return reply
}
