package routes

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"shortener/internal/storage"
)

func TestGetURLHandler(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name   string
		route  string
		method string
		want   want
	}{
		{
			name:   "Positive GET #1",
			route:  "/BFG9000x",
			method: http.MethodGet,
			want: want{
				contentType: `"text/plain; charset=utf-8"`,
				statusCode:  http.StatusTemporaryRedirect,
			},
		},
		{
			name:   "Negative GET #1",
			route:  "/MisSing",
			method: "GET",
			want: want{
				contentType: `"text/plain; charset=utf-8"`,
				statusCode:  http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := chi.NewRouter()
			router.Get("/{id}", GetURLHandler)
			URLMap := storage.Mapper
			URLMap.Set("BFG9000x", "www.yandex.ru")
			r := httptest.NewRequest(tt.method, tt.route, http.NoBody)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)
			res := w.Result()
			_, err := io.ReadAll(res.Body)
			if err != nil {
				return
			}
			err = res.Body.Close() // так требует golangci, defer с безымянной функцией не хочет
			if err != nil {
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want.statusCode, res.StatusCode)
		})
	}
}
