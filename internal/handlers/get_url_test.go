package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"shortener/internal/config"
	"shortener/internal/service"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Get(shortLink string) (string, bool) {
	args := m.Called(shortLink)
	return args.String(0), args.Bool(1)
}

func (m *MockDB) Save(shortLink, longLink string) {
	m.Called(shortLink, longLink)
}

func TestGetHandler(t *testing.T) {
	cfg := config.LoadConfig()

	mockedDB := &MockDB{}
	svc := &service.Service{DB: mockedDB, BaseURL: cfg.Server.BaseURL}

	type want struct {
		contentType string
		statusCode  int
		success     bool
	}
	cases := []struct {
		name    string
		route   string
		longURL string
		method  string
		want    want
	}{
		{
			name:    "Positive GET #1",
			route:   "BFG9000x",
			longURL: "https://example.org",
			method:  http.MethodGet,
			want: want{
				contentType: `"text/plain; charset=utf-8"`,
				statusCode:  http.StatusTemporaryRedirect,
				success:     true,
			},
		},
		{
			name:    "Positive GET #2",
			route:   "Xo0lK6n5",
			longURL: "https://dzen.ru",
			method:  http.MethodGet,
			want: want{
				contentType: `"text/plain; charset=utf-8"`,
				statusCode:  http.StatusTemporaryRedirect,
				success:     true,
			},
		},
		{
			name:    "Negative GET #1",
			route:   "MissingRoute",
			longURL: "",
			method:  http.MethodGet,
			want: want{
				contentType: `"text/plain; charset=utf-8"`,
				statusCode:  http.StatusBadRequest,
				success:     false,
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			mockedDB.On("Save", tt.route, tt.longURL).Return()
			mockedDB.On("Get", tt.route).Return(tt.longURL, tt.want.success)

			router := chi.NewRouter()
			router.Get("/{id}", GetHandler(svc))
			r := httptest.NewRequest(tt.method, "/"+tt.route, http.NoBody)
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
	}
}
