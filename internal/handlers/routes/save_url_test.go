package routes

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveURLHandler(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name   string
		route  string
		method string
		body   string
		want   want
	}{
		{
			name:   "Positive POST #1",
			route:  "/",
			method: http.MethodPost,
			body:   "https://yandex.ru",
			want: want{
				contentType: `"text/plain; charset=utf-8"`,
				statusCode:  http.StatusCreated,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyValue := strings.NewReader(tt.body)
			r := httptest.NewRequest(tt.method, tt.route, bodyValue)
			w := httptest.NewRecorder()
			SaveURLHandler(w, r)
			res := w.Result()
			defer func() {
				if err := res.Body.Close(); err != nil {
					fmt.Println("Ошибка при закрытии тела ответа!")
				}
			}()
			assert.Equal(t, tt.want.statusCode, res.StatusCode)
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.NotEmpty(t, resBody)
		})
	}
}
