//go:generate mockgen -source=handlers.go -destination=mocks/mocks.go
package handlers

import (
	"io"
	"net/http"

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

// Shorten Сокращает URL
func (h *handler) Shorten(w http.ResponseWriter, r *http.Request) {
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
