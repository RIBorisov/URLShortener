package handlers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"shortener/internal/logger"
	"shortener/internal/models"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"shortener/internal/config"
	"shortener/internal/service"
	"shortener/internal/storage"
)

func TestSaveHandler(t *testing.T) {
	cfg := config.LoadConfig()
	ctx := context.Background()
	log := &logger.Log{}
	log.Initialize("INFO")
	s, err := storage.LoadStorage(ctx, cfg, log)
	assert.NoError(t, err)
	svc := &service.Service{Storage: s, BaseURL: cfg.Service.BaseURL}
	user := &models.User{}
	type want struct {
		statusCode int
	}
	cases := []struct {
		name   string
		route  string
		method string
		body   string
		want   want
	}{
		{
			name:   "Positive #1",
			route:  "/",
			method: http.MethodPost,
			body:   "https://dzen.ru",
			want: want{
				statusCode: http.StatusCreated,
			},
		},
		{
			name:   "Positive #2",
			route:  "/",
			method: http.MethodPost,
			body:   "https://example.org",
			want: want{
				statusCode: http.StatusCreated,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			router := chi.NewRouter()
			router.Post("/", SaveHandler(svc, user))

			reqBody := strings.NewReader(tt.body)
			r := httptest.NewRequest(tt.method, tt.route, reqBody)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)
			res := w.Result()
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				log.Err("error when reading response body: ", err)
			}
			if err = res.Body.Close(); err != nil {
				log.Err("error when closing response body: ", err)
			}
			require.NoError(t, err)
			assert.NotEmpty(t, resBody)
			assert.Equal(t, res.StatusCode, tt.want.statusCode)
		})
	}
}
