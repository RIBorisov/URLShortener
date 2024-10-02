package tasks

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"shortener/internal/logger"
	"shortener/internal/service/mocks"
)

func TestRunBackgroundCleanupDB(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	log := &logger.Log{}
	log.Initialize("DEBUG")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockURLStorage(ctrl)

	cases := []struct {
		name         string
		cleanupURLs  []string
		cleanupError error
	}{
		{
			name:         "Positive #1",
			cleanupURLs:  []string{"url1", "url2"},
			cleanupError: nil,
		},
		{
			name:         "Positive #2",
			cleanupURLs:  []string{},
			cleanupError: nil,
		},
	}
	interval := 50 * time.Millisecond

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			mockStore.EXPECT().Cleanup(ctx).AnyTimes().Return(tt.cleanupURLs, tt.cleanupError)

			go Run(ctx, mockStore, log, interval)

			time.Sleep(150 * time.Millisecond)

			assert.ErrorIs(t, ctx.Err(), context.DeadlineExceeded)
		})
	}
}
