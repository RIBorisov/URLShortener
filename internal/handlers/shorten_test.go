package handlers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"shortener/internal/logger"
	"shortener/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"shortener/internal/config"
	"shortener/internal/service"
	"shortener/internal/storage"
)

func TestShortenHandler(t *testing.T) {
	const (
		ct    = "application/json"
		route = "/api/shorten"
	)
	ctx := context.Background()
	cfg := config.LoadConfig()
	log := &logger.Log{}
	userID := "123"
	log.Initialize("INFO")
	s, err := storage.LoadStorage(ctx, cfg, log)
	assert.NoError(t, err)

	svc := &service.Service{Storage: s, BaseURL: cfg.App.BaseURL, Log: log}
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
			oldCtx := r.Context()
			rWithCtx := r.WithContext(context.WithValue(oldCtx, models.CtxUserIDKey, userID))
			router.ServeHTTP(w, rWithCtx)
			res := w.Result()
			resBody, err := io.ReadAll(res.Body)
			assert.NoError(t, err)
			if err = res.Body.Close(); err != nil {
				log.Err("failed to close response body: ", err)
				return
			}
			assert.NotEmpty(t, resBody)
			assert.Equal(t, tc.want.statusCode, res.StatusCode)
			assert.Equal(t, tc.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
