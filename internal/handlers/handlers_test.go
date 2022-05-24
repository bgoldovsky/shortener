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

	"github.com/bgoldovsky/shortener/internal/app/services/urls"
	mockUrls "github.com/bgoldovsky/shortener/internal/app/services/urls/mocks"
)

const host = "http://localhost:8080"

func TestShortenHandler(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		shortcut    string
	}
	tests := []struct {
		name    string
		request string
		url     string
		id      string
		want    want
	}{
		{
			name: "success",
			url:  "https://avito.ru",
			id:   "xyz",
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

			genMock := mockUrls.NewMockgenerator(ctrl)
			genMock.EXPECT().ID().Return(tt.id)

			repoMock := mockUrls.NewMockrepo(ctrl)
			repoMock.EXPECT().Add(tt.id, tt.url)

			srv := urls.NewService(repoMock, genMock, host)

			httpHandler := New(srv)

			buffer := new(bytes.Buffer)
			buffer.WriteString(tt.url)
			request := httptest.NewRequest(http.MethodPost, tt.request, buffer)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(httpHandler.Shorten)
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
