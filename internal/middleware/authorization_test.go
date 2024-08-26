package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"shortener/internal/logger"
)

func BenchmarkGetUserID(b *testing.B) {
	tokenString := `"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.
eyJleHAiOjE3MjY0Mjc5NjksIlVzZXJJRCI6ImUwOTI2OGQ2LWVlMTQtNGU3Yi04MWZiLTUxOGU4M2JmMDM0NiJ9.
yTmWk0mALkC1Lb2j9Qcz70GqY-RA-BOUWX_0e_TbA0U"`
	secretKey := "!@#$%^YdBg0DS"
	log := &logger.Log{}
	log.Initialize("INFO")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getUserID(tokenString, secretKey, log)
	}
}

func TestBaseCheck_Middleware(t *testing.T) {
	log := &logger.Log{}
	log.Initialize("INFO")

	tests := []struct {
		name             string
		authHeader       string
		tokenCookieValue string
		wantStatus       int
	}{
		{
			name:             "Positive #1",
			authHeader:       "valid-token",
			tokenCookieValue: "valid-token",
			wantStatus:       http.StatusOK,
		},
		{
			name:             "Negative #1 (401)",
			authHeader:       "",
			tokenCookieValue: "valid-token",
			wantStatus:       http.StatusUnauthorized,
		},
		{
			name:             "Negative #2 (500)",
			authHeader:       "valid-token",
			tokenCookieValue: "",
			wantStatus:       http.StatusInternalServerError,
		},
		{
			name:             "Negative #3 (401)",
			authHeader:       "invalid-token",
			tokenCookieValue: "valid-token",
			wantStatus:       http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			if tt.tokenCookieValue != "" {
				req.AddCookie(&http.Cookie{Name: "token", Value: tt.tokenCookieValue})
			}

			rw := httptest.NewRecorder()
			handler := CheckAuth(log).Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			handler.ServeHTTP(rw, req)

			assert.Equal(t, tt.wantStatus, rw.Code)
		})
	}
}
