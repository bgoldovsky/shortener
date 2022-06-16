//go:generate mockgen -source=handlers.go -destination=mocks/mocks.go
package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type urlService interface {
	Shorten(url string) string
	Expand(id string) (string, error)
}

type handler struct {
	service urlService
}

func New(service urlService) *handler {
	return &handler{
		service: service,
	}
}

// ShortenV1 Сокращает URL, принимает и возвращает строку
func (h *handler) ShortenV1(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	shortcut := h.service.Shorten(string(b))

	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(shortcut))
	if err != nil {
		logrus.WithError(err).WithField("shortcut", shortcut).Error("write response error")
		return
	}
}

// ShortenV2 Сокращает URL, принимает и возвращает JSON
func (h *handler) ShortenV2(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	req := ShortenRequest{}
	if err = json.Unmarshal(b, &req); err != nil {
		http.Error(w, "request in not valid", http.StatusBadRequest)
		return
	}

	ok, err := govalidator.ValidateStruct(req)
	if err != nil || !ok {
		http.Error(w, "request in not valid", http.StatusBadRequest)
		return
	}

	shortcut := h.service.Shorten(req.URL)

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp := ShortenReply{ShortenUrl: shortcut}
	marshal, err := json.Marshal(&resp)
	if err != nil {
		logrus.WithError(err).WithField("resp", resp).Error("marshal response error")
		http.Error(w, err.Error(), 500)
		return
	}

	_, err = w.Write(marshal)
	if err != nil {
		logrus.WithError(err).WithField("shortcut", shortcut).Error("write response error")
		return
	}
}

// Expand Возвращает полный URL по идентификатору сокращенного
func (h *handler) Expand(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "id parameter is empty", http.StatusBadRequest)
		return
	}

	url, err := h.service.Expand(id)
	if err != nil {
		http.Error(w, "url not found", http.StatusNoContent)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
