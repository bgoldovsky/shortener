package handlers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bgoldovsky/shortener/internal/app/models"
	mockHandlers "github.com/bgoldovsky/shortener/internal/handlers/mocks"
)

func TestShortenV1Handler(t *testing.T) {
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
		err      error
		want     want
	}{
		{
			name:     "success",
			url:      "https://avito.ru",
			shortcut: "http://localhost:8080/xyz",
			err:      nil,
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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			srvMock := mockHandlers.NewMockurlService(ctrl)
			srvMock.EXPECT().Shorten(tt.url).Return(tt.shortcut, nil)

			httpHandler := New(srvMock)

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

func TestShortenV2Handler_Success(t *testing.T) {
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
			err:      nil,
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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			srvMock := mockHandlers.NewMockurlService(ctrl)
			srvMock.EXPECT().Shorten(tt.url).Return(tt.shortcut, tt.err)

			httpHandler := New(srvMock)

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

func TestShortenV2Handler_BadRequest(t *testing.T) {
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

			httpHandler := New(nil)

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

func TestGetUrlsHandler_Success(t *testing.T) {
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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			srvMock := mockHandlers.NewMockurlService(ctrl)
			srvMock.EXPECT().GetUrls().Return(tt.urls, tt.err)

			httpHandler := New(srvMock)

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
