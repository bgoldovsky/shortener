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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			srvMock := mockHandlers.NewMockurlService(ctrl)
			srvMock.EXPECT().Shorten(tt.url).Return(tt.shortcut)

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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			srvMock := mockHandlers.NewMockurlService(ctrl)
			srvMock.EXPECT().Shorten(tt.url).Return(tt.shortcut)

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
