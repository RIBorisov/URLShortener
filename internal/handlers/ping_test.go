package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"shortener/internal/config"
	"shortener/internal/logger"
	"shortener/internal/service"
	"shortener/internal/service/mocks"
)

func TestPingHandler(t *testing.T) {
	const (
		GET   = http.MethodGet
		route = "/ping"
	)
	cfg := config.LoadConfig()
	log := &logger.Log{}
	log.Initialize("INFO")

	type want struct {
		respErr error
		status  int
	}
	tests := []struct {
		name      string
		callTimes int
		method    string
		want      want
	}{
		{
			name:      "Positive #1",
			callTimes: 1,
			method:    GET,
			want: want{
				respErr: nil,
				status:  http.StatusOK,
			}},
		{
			name:      "Negative #1",
			callTimes: 1,
			method:    GET,
			want: want{
				respErr: errors.New("unexpected error"),
				status:  http.StatusInternalServerError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStore := mocks.NewMockURLStorage(ctrl)
			mockStore.EXPECT().Ping(ctx).Times(tt.callTimes).Return(tt.want.respErr)

			svc := &service.Service{Storage: mockStore, BaseURL: cfg.App.BaseURL, Log: log}
			handler := PingHandler(svc)
			req, err := http.NewRequest(GET, route, http.NoBody)
			assert.NoError(t, err)

			w := httptest.NewRecorder()

			handler(w, req)
			resp := w.Result()
			assert.NoError(t, resp.Body.Close())
			assert.Equal(t, tt.want.status, resp.StatusCode)
		})
	}
}
