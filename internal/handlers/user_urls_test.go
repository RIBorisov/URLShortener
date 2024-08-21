package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"shortener/internal/config"
	"shortener/internal/logger"
	"shortener/internal/models"
	"shortener/internal/service"
	"shortener/internal/storage"
)

func BenchmarkGetURLsHandler(b *testing.B) {
	cfg := config.LoadConfig()
	ctx := context.Background()
	userID := "100500"
	log := &logger.Log{}
	log.Initialize("INFO")
	s, _ := storage.LoadStorage(ctx, cfg, log)
	svc := &service.Service{Storage: s, BaseURL: cfg.Service.BaseURL}

	router := chi.NewRouter()
	router.Get("/api/user/urls", GetURLsHandler(svc))
	r := httptest.NewRequest(http.MethodGet, "/api/user/urls", http.NoBody)
	oldCtx := r.Context()
	rWithCtx := r.WithContext(context.WithValue(oldCtx, models.CtxUserIDKey, userID))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, rWithCtx)
	}
}

func TestGetURLsHandler(t *testing.T) {
	const route = "/api/user/urls"
	ctx := context.Background()
	cfg := config.LoadConfig()
	log := &logger.Log{}
	log.Initialize("INFO")
	s, err := storage.LoadStorage(ctx, cfg, log)
	assert.NoError(t, err)

	svc := &service.Service{Storage: s, BaseURL: cfg.Service.BaseURL}

	tests := []struct {
		name         string
		userID       string
		expectedCode int
		expectedBody models.UserURLs
	}{
		{
			name:         "POSITIVE #1",
			userID:       "abc123qwe987",
			expectedCode: http.StatusNoContent,
		},
		{
			name:         "POSITIVE #2",
			userID:       "abc123qwe986",
			expectedCode: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := GetURLsHandler(svc)
			r := httptest.NewRequest(http.MethodGet, route, nil)
			w := httptest.NewRecorder()
			oldCtx := r.Context()
			rWithCtx := r.WithContext(context.WithValue(oldCtx, models.CtxUserIDKey, tt.userID))
			handler.ServeHTTP(w, rWithCtx)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}
