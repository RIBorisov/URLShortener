package routes

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSaveURLHandler(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		body        string
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
			assert.Equal(t, tt.want.statusCode, res.StatusCode)
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					fmt.Println("Ошибка при закрытии тела ответа")
					return
				}
			}(r.Body)
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.NotEmpty(t, resBody)
		})
	}
}
