//go:generate mockgen -source=handlers.go -destination=mocks/mocks.go
package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"

	"github.com/bgoldovsky/shortener/internal/app/models"
	urlsSrv "github.com/bgoldovsky/shortener/internal/app/services/urls"
)

type urlsService interface {
	Shorten(ctx context.Context, url, userID string) (string, error)
	ShortenBatch(ctx context.Context, originalURLs []models.OriginalURL, userID string) ([]models.URL, error)
	Expand(ctx context.Context, id string) (string, error)
	GetUrls(ctx context.Context, userID string) ([]models.URL, error)
}

type auth interface {
	UserID(ctx context.Context) string
}

type infra interface {
	Ping(ctx context.Context) bool
}

type cleaner interface {
	Queue(urls models.UserCollection)
}

type handler struct {
	urlsService urlsService
	auth        auth
	infra       infra
	cleaner     cleaner
}

func New(urlsService urlsService, auth auth, infra infra, cleaner cleaner) *handler {
	return &handler{
		urlsService: urlsService,
		auth:        auth,
		infra:       infra,
		cleaner:     cleaner,
	}
}

// ShortenV1 Сокращает URL, принимает и возвращает строку
func (h *handler) ShortenV1(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userID := h.auth.UserID(r.Context())
	url := string(b)
	statusCode := http.StatusCreated

	shortcut, err := h.urlsService.Shorten(r.Context(), url, userID)
	if err != nil {
		if !errors.Is(err, urlsSrv.ErrNotUniqueURL) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		statusCode = http.StatusConflict
	}

	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(statusCode)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	userID := h.auth.UserID(r.Context())
	statusCode := http.StatusCreated

	shortcut, err := h.urlsService.Shorten(r.Context(), req.URL, userID)
	if err != nil {
		if !errors.Is(err, urlsSrv.ErrNotUniqueURL) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		statusCode = http.StatusConflict
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(statusCode)

	resp := ShortenReply{ShortURL: shortcut}
	marshal, err := json.Marshal(&resp)
	if err != nil {
		logrus.WithError(err).WithField("resp", resp).Error("marshal response error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(marshal)
	if err != nil {
		logrus.WithError(err).WithField("shortcut", shortcut).Error("write response error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// ShortenBatch Сокращает несколько URL
func (h *handler) ShortenBatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var req []ShortenBatchRequest
	if err = json.Unmarshal(b, &req); err != nil {
		http.Error(w, "request in not valid", http.StatusBadRequest)
		return
	}

	if len(req) == 0 {
		http.Error(w, "url list not specified", http.StatusBadRequest)
		return
	}

	for idx := range req {
		if ok, err := govalidator.ValidateStruct(req[idx]); err != nil || !ok {
			http.Error(w, "element of url list not valid", http.StatusBadRequest)
			return
		}
	}

	userID := h.auth.UserID(r.Context())
	originalUrls := toShortenBatchRequest(req)

	urls, err := h.urlsService.ShortenBatch(r.Context(), originalUrls, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp := toShortenBatchReply(urls)
	marshal, err := json.Marshal(&resp)
	if err != nil {
		logrus.WithError(err).WithField("resp", resp).Error("marshal response error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(marshal)
	if err != nil {
		logrus.WithError(err).WithField("urls", urls).Error("write response error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Expand Возвращает полный URL по идентификатору сокращенного
func (h *handler) Expand(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "id parameter is empty", http.StatusBadRequest)
		return
	}

	url, err := h.urlsService.Expand(r.Context(), id)
	if err != nil {
		if errors.Is(err, urlsSrv.ErrURLNotFound) {
			http.Error(w, "url not found", http.StatusNoContent)
			return
		}

		if errors.Is(err, urlsSrv.ErrURLDeleted) {
			http.Error(w, "url has been deleted", http.StatusGone)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// GetUrls Возвращает список всех сокращенных URL пользователя
func (h *handler) GetUrls(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := h.auth.UserID(r.Context())

	urls, err := h.urlsService.GetUrls(r.Context(), userID)
	if err != nil {
		http.Error(w, "get urls error", http.StatusInternalServerError)
		return
	}
	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := toGetUrlsReply(urls)
	marshal, err := json.Marshal(&resp)
	if err != nil {
		logrus.WithError(err).WithField("resp", resp).Error("marshal response error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(marshal)
	if err != nil {
		logrus.WithError(err).WithField("resp", resp).Error("write response error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// DeleteUrls Удаляет список сокращенных URL пользователя
func (h *handler) DeleteUrls(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := h.auth.UserID(r.Context())

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var urlIDs []string
	if err := json.Unmarshal(b, &urlIDs); err != nil {
		http.Error(w, "request in not valid", http.StatusBadRequest)
		return
	}

	userURLs := models.UserCollection{
		UserID: userID,
		URLIDs: urlIDs,
	}

	h.cleaner.Queue(userURLs)

	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusAccepted)
}

// Ping Проверяет доступность базы данных
func (h *handler) Ping(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if success := h.infra.Ping(r.Context()); !success {
		http.Error(w, "ping database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
