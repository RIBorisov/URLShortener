package handlers

import (
	"bytes"
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

func TestDeleteURLsHandler(t *testing.T) {
	const (
		DELETE = http.MethodDelete
		route  = "/api/user/urls"
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
		body      string
		want      want
	}{
		{
			name:      "Positive DELETE #1",
			callTimes: 1,
			method:    DELETE,
			body:      `["short1", "short2", "short3"]`,
			want: want{
				respErr: nil,
				status:  http.StatusAccepted,
			},
		},
		{
			name:      "Negative DELETE #1",
			callTimes: 1,
			method:    DELETE,
			body:      `["1", "x", "3"]`,
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
			mockStore.EXPECT().DeleteURLs(ctx, gomock.Any()).Times(tt.callTimes).Return(tt.want.respErr)

			svc := &service.Service{Storage: mockStore, BaseURL: cfg.Service.BaseURL, Log: log}
			handler := DeleteURLsHandler(svc)

			req, err := http.NewRequest(DELETE, route, bytes.NewBufferString(tt.body))

			assert.NoError(t, err)

			w := httptest.NewRecorder()

			handler(w, req)
			resp := w.Result()
			assert.NoError(t, resp.Body.Close())

			assert.Equal(t, tt.want.status, resp.StatusCode)
		})
	}
}
