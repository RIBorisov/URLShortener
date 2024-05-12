package routes

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"shortener/internal/storage"
)

func TestSaveHandler(t *testing.T) {
	db := storage.GetStorage()
	type want struct {
		statusCode int
	}
	cases := []struct {
		name   string
		route  string
		method string
		body   string
		want   want
	}{
		{
			name:   "Positive #1",
			route:  "/",
			method: http.MethodPost,
			body:   "https://dzen.ru",
			want: want{
				statusCode: http.StatusCreated,
			},
		},
		{
			name:   "Positive #2",
			route:  "/",
			method: http.MethodPost,
			body:   "https://example.org",
			want: want{
				statusCode: http.StatusCreated,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			handler := SaveHandler(db)
			reqBody := strings.NewReader(tt.body)
			r := httptest.NewRequest(tt.method, tt.route, reqBody)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, r)
			res := w.Result()
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				log.Printf("error when reading response body")
			}
			if err = res.Body.Close(); err != nil {
				log.Printf("error when closing response body")
			}
			require.NoError(t, err)
			assert.NotEmpty(t, resBody)
			assert.Equal(t, res.StatusCode, tt.want.statusCode)
		})
	}
}
