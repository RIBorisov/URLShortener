package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"shortener/internal/config"
	"shortener/internal/logger"
	"shortener/internal/models"
	"shortener/internal/service"
	"shortener/internal/storage"
)

func TestBatchHandler(t *testing.T) {
	const (
		POST  = http.MethodPost
		route = "/api/shorten/batch"
	)
	cfg := config.LoadConfig()
	log := &logger.Log{}
	log.Initialize("INFO")
	ctx := context.Background()
	s, err := storage.LoadStorage(ctx, cfg, log)
	assert.NoError(t, err)
	svc := &service.Service{Storage: s, BaseURL: cfg.App.BaseURL}

	tests := []struct {
		name       string
		userID     string
		body       []models.BatchRequest
		wantStatus int
	}{
		{
			name:   "Positive #1",
			userID: "100500",
			body: []models.BatchRequest{
				{CorrelationID: "id1", OriginalURL: "t.me"},
				{CorrelationID: "id2", OriginalURL: "t.him"},
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "Negative #1",
			userID:     "100500",
			body:       []models.BatchRequest{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:   "Negative #2",
			userID: "",
			body: []models.BatchRequest{
				{CorrelationID: "id13", OriginalURL: "t.me"},
				{CorrelationID: "id21", OriginalURL: "t.him"},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := chi.NewRouter()
			router.Post(route, BatchHandler(svc))
			b, err := json.Marshal(tt.body)
			assert.NoError(t, err)
			r, err := http.NewRequest(POST, route, strings.NewReader(string(b)))
			assert.NoError(t, err)

			if tt.userID != "" {
				r = r.WithContext(context.WithValue(r.Context(), models.CtxUserIDKey, tt.userID))
			}

			w := httptest.NewRecorder()

			router.ServeHTTP(w, r)
			resp := w.Result()
			assert.NoError(t, resp.Body.Close())

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}
