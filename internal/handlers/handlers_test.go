package handlers

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bgoldovsky/shortener/internal/app/models"
	"github.com/bgoldovsky/shortener/internal/app/services/urls"
	mockHandlers "github.com/bgoldovsky/shortener/internal/handlers/mocks"
)

const defaultUserID = "user123"

func TestHandler_ShortenV1_Success(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		shortcut    string
	}

	tests := []struct {
		name     string
		request  string
		url      string
		shortcut string
		want     want
	}{
		{
			name:     "success",
			url:      "https://avito.ru",
			shortcut: "http://localhost:8080/xyz",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  201,
				shortcut:    "http://localhost:8080/xyz",
			},
			request: "/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			urlsSrvMock := mockHandlers.NewMockurlsService(ctrl)
			urlsSrvMock.EXPECT().Shorten(ctx, tt.url, defaultUserID).Return(tt.shortcut, nil)

			authMock := mockHandlers.NewMockauth(ctrl)
			authMock.EXPECT().UserID(gomock.Any()).Return(defaultUserID)

			httpHandler := New(urlsSrvMock, authMock, nil, nil)

			buffer := new(bytes.Buffer)
			buffer.WriteString(tt.url)
			request := httptest.NewRequest(http.MethodPost, tt.request, buffer)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(httpHandler.ShortenV1)
			h.ServeHTTP(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			userResult, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			require.NoError(t, err)

			assert.Equal(t, tt.want.shortcut, string(userResult))
		})
	}
}

func TestHandler_ShortenV2_Success(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		response    string
	}
	tests := []struct {
		name     string
		request  string
		url      string
		body     string
		shortcut string
		want     want
	}{
		{
			name:     "success",
			url:      "https://avito.ru",
			body:     "{\"url\":\"https://avito.ru\"}",
			shortcut: "http://localhost:8080/xyz",
			want: want{
				contentType: "application/json",
				statusCode:  201,
				response:    "{\"result\":\"http://localhost:8080/xyz\"}",
			},
			request: "/api/shorten",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			urlSrvMock := mockHandlers.NewMockurlsService(ctrl)
			urlSrvMock.EXPECT().Shorten(ctx, tt.url, defaultUserID).Return(tt.shortcut, nil)

			authMock := mockHandlers.NewMockauth(ctrl)
			authMock.EXPECT().UserID(gomock.Any()).Return(defaultUserID)

			httpHandler := New(urlSrvMock, authMock, nil, nil)

			buffer := new(bytes.Buffer)
			buffer.WriteString(tt.body)
			request := httptest.NewRequest(http.MethodPost, tt.request, buffer)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(httpHandler.ShortenV2)
			h.ServeHTTP(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			userResult, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			require.NoError(t, err)

			assert.Equal(t, tt.want.response, string(userResult))
		})
	}
}

func TestHandler_ShortenV2_BadRequest(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		response    string
	}
	tests := []struct {
		name     string
		request  string
		url      string
		body     string
		shortcut string
		want     want
	}{
		{
			name:     "bad-request",
			url:      "https://avito.ru",
			body:     "{\"url\":\"\"}",
			shortcut: "http://localhost:8080/xyz",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
				response:    "request in not valid\n",
			},
			request: "/api/shorten",
		},
		{
			name:     "bad-request",
			url:      "https://avito.ru",
			body:     "{\"url\":\"qwerty\"}",
			shortcut: "http://localhost:8080/xyz",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
				response:    "request in not valid\n",
			},
			request: "/api/shorten",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			httpHandler := New(nil, nil, nil, nil)

			buffer := new(bytes.Buffer)
			buffer.WriteString(tt.body)
			request := httptest.NewRequest(http.MethodPost, tt.request, buffer)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(httpHandler.ShortenV2)
			h.ServeHTTP(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			userResult, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			require.NoError(t, err)

			assert.Equal(t, tt.want.response, string(userResult))
		})
	}
}

func TestHandler_ShortenV2_Conflict(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		response    string
	}
	tests := []struct {
		name     string
		request  string
		url      string
		body     string
		shortcut string
		err      error
		want     want
	}{
		{
			name:     "success",
			url:      "https://avito.ru",
			body:     "{\"url\":\"https://avito.ru\"}",
			shortcut: "http://localhost:8080/xyz",
			err:      urls.ErrNotUniqueURL,
			want: want{
				contentType: "application/json",
				statusCode:  409,
				response:    "{\"result\":\"http://localhost:8080/xyz\"}",
			},
			request: "/api/shorten",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			urlSrvMock := mockHandlers.NewMockurlsService(ctrl)
			urlSrvMock.EXPECT().Shorten(ctx, tt.url, defaultUserID).Return(tt.shortcut, tt.err)

			authMock := mockHandlers.NewMockauth(ctrl)
			authMock.EXPECT().UserID(gomock.Any()).Return(defaultUserID)

			httpHandler := New(urlSrvMock, authMock, nil, nil)

			buffer := new(bytes.Buffer)
			buffer.WriteString(tt.body)
			request := httptest.NewRequest(http.MethodPost, tt.request, buffer)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(httpHandler.ShortenV2)
			h.ServeHTTP(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			userResult, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			require.NoError(t, err)

			assert.Equal(t, tt.want.response, string(userResult))
		})
	}
}

func TestHandler_Expand_Success(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		response    string
		location    string
	}
	tests := []struct {
		name     string
		request  string
		url      string
		urlID    string
		shortcut string
		err      error
		want     want
	}{
		{
			name:     "success",
			url:      "https://avito.ru",
			urlID:    "xyz",
			shortcut: "http://localhost:8080/xyz",
			err:      nil,
			want: want{
				contentType: "",
				statusCode:  307,
				response:    "",
			},
			request: "/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			urlsSrvMock := mockHandlers.NewMockurlsService(ctrl)
			urlsSrvMock.EXPECT().Expand(gomock.Any(), tt.urlID).Return(tt.url, tt.err)

			httpHandler := New(urlsSrvMock, nil, nil, nil)

			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.urlID)

			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()
			h := http.HandlerFunc(httpHandler.Expand)

			h.ServeHTTP(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			userResult, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			require.NoError(t, err)

			assert.Equal(t, tt.want.response, string(userResult))
		})
	}
}

func TestHandler_GetUrls_Success(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		response    string
	}
	tests := []struct {
		name    string
		request string
		urls    []models.URL
		err     error
		want    want
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
			want: want{
				contentType: "application/json",
				statusCode:  200,
				response:    "[{\"short_url\":\"http://localhost:8080/xyz\",\"original_url\":\"https://avito.ru\"},{\"short_url\":\"http://localhost:8080/qwerty\",\"original_url\":\"https://yandex.ru\"}]",
			},
			request: "/api/user/urls",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			urlsSrvMock := mockHandlers.NewMockurlsService(ctrl)
			urlsSrvMock.EXPECT().GetUrls(ctx, defaultUserID).Return(tt.urls, tt.err)

			authMock := mockHandlers.NewMockauth(ctrl)
			authMock.EXPECT().UserID(gomock.Any()).Return(defaultUserID)

			httpHandler := New(urlsSrvMock, authMock, nil, nil)

			request := httptest.NewRequest(http.MethodGet, tt.request, nil)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(httpHandler.GetUrls)
			h.ServeHTTP(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			userResult, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			require.NoError(t, err)

			assert.Equal(t, tt.want.response, string(userResult))
		})
	}
}

func TestHandler_Ping(t *testing.T) {
	type want struct {
		statusCode int
	}
	tests := []struct {
		name    string
		request string
		success bool
		want    want
	}{
		{
			name:    "success",
			success: true,
			want: want{
				statusCode: 200,
			},
			request: "/ping",
		},
		{
			name:    "fail",
			success: false,
			want: want{
				statusCode: 500,
			},
			request: "/ping",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			infraMock := mockHandlers.NewMockinfra(ctrl)
			infraMock.EXPECT().Ping(ctx).Return(tt.success)

			httpHandler := New(nil, nil, infraMock, nil)

			request := httptest.NewRequest(http.MethodGet, tt.request, nil)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(httpHandler.Ping)
			h.ServeHTTP(w, request)

			result := w.Result()
			err := result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}

func TestHandler_ShortenBatch_Success(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		response    string
	}
	tests := []struct {
		name         string
		request      string
		originalURLs []models.OriginalURL
		urls         []models.URL
		err          error
		body         string
		want         want
	}{
		{
			name: "success",
			originalURLs: []models.OriginalURL{
				{
					CorrelationID: "qwerty",
					URL:           "https://avito.ru",
				},
			},
			urls: []models.URL{
				{
					ShortURL:    "http://localhost:8080/xyz",
					OriginalURL: "https://avito.ru",
				},
			},
			body: "[{\"correlation_id\":\"qwerty\",\"original_url\":\"https://avito.ru\"}]",
			want: want{
				contentType: "application/json",
				statusCode:  201,
				response:    "[{\"correlation_id\":\"\",\"short_url\":\"http://localhost:8080/xyz\"}]",
			},
			request: "/api/shorten/batch",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			authMock := mockHandlers.NewMockauth(ctrl)
			authMock.EXPECT().UserID(gomock.Any()).Return(defaultUserID)

			urlsSrvMock := mockHandlers.NewMockurlsService(ctrl)
			urlsSrvMock.EXPECT().ShortenBatch(ctx, tt.originalURLs, defaultUserID).Return(tt.urls, tt.err)

			httpHandler := New(urlsSrvMock, authMock, nil, nil)

			buffer := new(bytes.Buffer)
			buffer.WriteString(tt.body)
			request := httptest.NewRequest(http.MethodPost, tt.request, buffer)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(httpHandler.ShortenBatch)
			h.ServeHTTP(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			userResult, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			require.NoError(t, err)

			assert.Equal(t, tt.want.response, string(userResult))
		})
	}
}

func TestHandler_ShortenBatch_BadRequest(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		response    string
	}
	tests := []struct {
		name         string
		request      string
		originalURLs []models.OriginalURL
		urls         []models.URL
		err          error
		body         string
		want         want
	}{
		{
			name: "bad request",
			originalURLs: []models.OriginalURL{
				{
					CorrelationID: "qwerty",
					URL:           "https://avito.ru",
				},
			},
			urls: []models.URL{
				{
					ShortURL:    "http://localhost:8080/xyz",
					OriginalURL: "https://avito.ru",
				},
			},
			body: "[\"original_url\":\"https://avito.ru\"}]",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
				response:    "request in not valid\n",
			},
			request: "/api/shorten/batch",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			httpHandler := New(nil, nil, nil, nil)

			buffer := new(bytes.Buffer)
			buffer.WriteString(tt.body)
			request := httptest.NewRequest(http.MethodPost, tt.request, buffer)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(httpHandler.ShortenBatch)
			h.ServeHTTP(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			userResult, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			require.NoError(t, err)

			assert.Equal(t, tt.want.response, string(userResult))
		})
	}
}
