package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"shortener/internal/config"
	"shortener/internal/logger"
	"shortener/internal/service"
	"shortener/internal/service/mocks"
)

func TestBatchHandler(t *testing.T) {
	const (
		POST  = http.MethodPost
		route = "/api/shorten/batch"
	)
	cfg := config.LoadConfig()
	log := &logger.Log{}
	log.Initialize("INFO")

	type want struct {
		respErr  error
		status   int
		response any
	}
	tests := []struct {
		name         string
		callTimes    int
		callGetTimes int
		method       string
		body         string
		want         want
	}{
		{
			name:         "Negative POST #1",
			callTimes:    0,
			callGetTimes: 0,
			method:       POST,
			body: `correlation_id": "id1","original_url": "https://ya.ru"},
{"correlation_id": "id2","original_url": "https://t.me"}]`,
			want: want{
				respErr:  errors.New("failed to decode request into model"),
				status:   http.StatusInternalServerError,
				response: nil,
			},
		},
		{
			name:      "Negative POST #2",
			callTimes: 0,
			method:    POST,
			body:      `[]`,
			want: want{
				respErr:  errors.New("empty request batch"),
				status:   http.StatusBadRequest,
				response: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStore := mocks.NewMockURLStorage(ctrl)

			mockStore.EXPECT().Get(ctx, gomock.Any()).Times(tt.callGetTimes).Return("", nil)
			mockStore.EXPECT().BatchSave(ctx, tt.body).Times(tt.callTimes).Return(tt.want.response, tt.want.respErr)

			svc := &service.Service{Storage: mockStore, BaseURL: cfg.Service.BaseURL, Log: log}
			handler := BatchHandler(svc)

			req, err := http.NewRequest(POST, route, strings.NewReader(tt.body))

			assert.NoError(t, err)

			w := httptest.NewRecorder()

			handler(w, req)
			resp := w.Result()
			assert.NoError(t, resp.Body.Close())

			assert.Equal(t, tt.want.status, resp.StatusCode)
		})
	}
}
