package handlers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"shortener/internal/models"
	"shortener/internal/storage"
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

func (m *MockDB) Get(_ context.Context, shortLink string) (string, error) {
	args := m.Called(shortLink)
	return args.String(0), args.Error(1)
}

func (m *MockDB) Save(_ context.Context, shortLink, longLink string, _ *models.User) error {
	m.Called(shortLink, longLink)
	return nil
}

func (m *MockDB) BatchSave(_ context.Context, _ models.BatchArray, _ *models.User) (models.BatchArray, error) {
	return nil, nil
}

func (m *MockDB) GetByUserID(_ context.Context, _ *models.User) ([]models.BaseRow, error) {
	return nil, nil
}

func (m *MockDB) Close() error {
	return nil
}

func (m *MockDB) Ping(_ context.Context) error {
	return nil
}

func TestGetHandler(t *testing.T) {
	cfg := config.LoadConfig()

	mockedDB := &MockDB{}
	svc := &service.Service{Storage: mockedDB, BaseURL: cfg.Service.BaseURL}

	type want struct {
		contentType string
		statusCode  int
		err         error
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
				err:         nil,
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
				err:         nil,
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
				err:         storage.ErrURLNotFound,
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			mockedDB.On("Save", tt.route, tt.longURL).Return()
			mockedDB.On("Get", tt.route).Return(tt.longURL, tt.want.err)

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
