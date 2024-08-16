package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

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
