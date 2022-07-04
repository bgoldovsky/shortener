package handlers

type ShortenRequest struct {
	URL string `json:"url" valid:"url,required"`
}

type ShortenReply struct {
	ShortenURL string `json:"result"`
}

type GetUrlsReply struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
