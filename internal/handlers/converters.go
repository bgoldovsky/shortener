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

func toShortenBatchRequest(model []ShortenBatchRequest) []models.OriginalURL {
	reply := make([]models.OriginalURL, len(model))

	for idx, m := range model {
		reply[idx] = models.OriginalURL{
			CorrelationID: m.CorrelationID,
			URL:           m.OriginalURL,
		}
	}

	return reply
}

func toShortenBatchReply(model []models.URL) []ShortenBatchReply {
	reply := make([]ShortenBatchReply, len(model))

	for idx, m := range model {
		reply[idx] = ShortenBatchReply{
			CorrelationID: m.CorrelationID,
			ShortURL:      m.ShortURL,
		}
	}

	return reply
}
