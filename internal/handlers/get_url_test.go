package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"shortener/internal/config"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"shortener/internal/storage"
)

func TestGetHandler(t *testing.T) {
	cfg := config.LoadConfig()
	db := storage.LoadStorage()
	db.Save("BFG9000x", "https://example.org")
	db.Save("Xo0lK6n5", "https://dzen.ru")
	type want struct {
		contentType string
		statusCode  int
	}
	cases := []struct {
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
			name:   "Positive GET #2",
			route:  "/Xo0lK6n5",
			method: http.MethodGet,
			want: want{
				contentType: `"text/plain; charset=utf-8"`,
				statusCode:  http.StatusTemporaryRedirect,
			},
		},
		{
			name:   "Negative GET #1",
			route:  "/MissingRoute",
			method: http.MethodGet,
			want: want{
				contentType: `"text/plain; charset=utf-8"`,
				statusCode:  http.StatusBadRequest,
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Run(tt.name, func(t *testing.T) {
				router := chi.NewRouter()
				router.Get("/{id}", GetHandler(db, cfg))
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
				assert.NotEmpty(t, res.Body)
				assert.Equal(t, tt.want.statusCode, res.StatusCode)
			})
		})
	}
}
