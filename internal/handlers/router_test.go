package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"shortener/internal/config"
	"shortener/internal/logger"
	"shortener/internal/service"
	"shortener/internal/service/mocks"
)

func TestNewRouter(t *testing.T) {
	cfg := config.LoadConfig()
	log := &logger.Log{}
	log.Initialize("INFO")

	t.Run("Positive #1", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStore := mocks.NewMockURLStorage(ctrl)
		svc := &service.Service{Storage: mockStore, BaseURL: cfg.Service.BaseURL, Log: log}

		r := NewRouter(svc)
		assert.NotEmpty(t, r)
	})
}
