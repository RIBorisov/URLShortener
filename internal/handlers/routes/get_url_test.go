package routes

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"shortener/internal/storage"
	"testing"
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
			URLMap := storage.Mapper
			URLMap.Set("BFG9000x", "www.yandex.ru")
			r := httptest.NewRequest(tt.method, tt.route, nil)
			w := httptest.NewRecorder()
			GetURLHandler(w, r)
			res := w.Result()
			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.Equal(t, tt.want.statusCode, res.StatusCode)
		})
	}
}
