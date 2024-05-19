package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"shortener/internal/config"
	"shortener/internal/service"
	"shortener/internal/storage"
	"strings"
	"testing"
)

func TestShortenHandler(t *testing.T) {
	const (
		ct    = "application/json"
		route = "/api/shorten"
	)

	cfg := config.LoadConfig()
	db := storage.LoadStorage()
	svc := &service.Service{DB: db, BaseURL: cfg.Server.BaseURL}
	type want struct {
		statusCode  int
		contentType string
	}
	cases := []struct {
		name   string
		method string
		body   string
		want   want
	}{
		{
			name:   "Positive #1",
			method: http.MethodPost,
			body:   `{"request": {"type": "SimpleRequest", "url": "https://www.kinopoisk.ru/"}}`,
			want: want{
				statusCode:  http.StatusCreated,
				contentType: ct,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			router := chi.NewRouter()
			router.Post(route, ShortenHandler(svc))

			reqBody := strings.NewReader(tc.body)
			r := httptest.NewRequest(tc.method, route, reqBody)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)

			res := w.Result()
			resBody, err := io.ReadAll(res.Body)
			assert.NoError(t, err)
			assert.NotEmpty(t, resBody)
			assert.Equal(t, tc.want.statusCode, res.StatusCode)
			assert.Equal(t, tc.want.contentType, res.Header.Get("Content-Type"))

		})
	}
}
