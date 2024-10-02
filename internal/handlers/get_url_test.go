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
	"shortener/internal/storage"
)

func TestGetHandler(t *testing.T) {
	const GET = http.MethodGet
	cfg := config.LoadConfig()
	log := &logger.Log{}
	log.Initialize("INFO")

	type want struct {
		response any
		respErr  error
		status   int
	}
	tests := []struct {
		name      string
		route     string
		callTimes int
		longURL   string
		method    string
		want      want
	}{
		{
			name:      "Positive GET #1",
			route:     "/BFG9000x",
			callTimes: 1,
			longURL:   "https://example.org",
			method:    http.MethodGet,
			want: want{
				response: "https://example.org",
				respErr:  nil,
				status:   http.StatusTemporaryRedirect,
			},
		},
		{
			name:      "Positive GET #2",
			route:     "/Xo0lK6n5",
			callTimes: 1,
			longURL:   "https://dzen.ru",
			method:    http.MethodGet,
			want: want{
				response: "https://dzen.ru",
				respErr:  nil,
				status:   http.StatusTemporaryRedirect,
			},
		},
		{
			name:      "Negative GET #1",
			route:     "/QwsDqr1",
			callTimes: 1,
			longURL:   "Missing",
			method:    GET,
			want: want{
				response: "asd",
				respErr:  storage.ErrURLDeleted,
				status:   http.StatusGone,
			},
		},
		{
			name:      "Negative GET #2",
			route:     "/BadRequestShort",
			callTimes: 1,
			longURL:   "BadRequest",
			method:    GET,
			want: want{
				response: "",
				respErr:  errors.New("unexpected error leading to 400"),
				status:   http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStore := mocks.NewMockURLStorage(ctrl)
			mockStore.EXPECT().Get(ctx, gomock.Any()).Times(tt.callTimes).Return(tt.want.response, tt.want.respErr)

			svc := &service.Service{Storage: mockStore, BaseURL: cfg.App.BaseURL, Log: log}
			handler := GetHandler(svc)
			req, err := http.NewRequest(GET, tt.route, http.NoBody)
			assert.NoError(t, err)

			w := httptest.NewRecorder()

			handler(w, req)
			resp := w.Result()
			assert.NoError(t, resp.Body.Close())
			assert.Equal(t, tt.want.status, resp.StatusCode)
		})
	}
}
