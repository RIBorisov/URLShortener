package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"shortener/internal/logger"
)

func Test_initApp(t *testing.T) {
	type args struct {
		log *logger.Log
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Positive #1",
			args: args{
				log: &logger.Log{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("SERVER_ADDRESS", "localhost")
			t.Setenv("SERVER_PORT", "8080")
			t.Setenv("ENABLE_HTTPS", "0")
			tt.args.log.Initialize("INFO")

			shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 500*time.Millisecond)
			defer shutdownRelease()

			go func() {
				err := initApp(shutdownCtx, tt.args.log)
				assert.NoError(t, err)
			}()
			<-shutdownCtx.Done()
		})
	}
}
