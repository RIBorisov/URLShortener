// Package middleware using for processing requests with additional business logic.
package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"shortener/internal/logger"
	"shortener/internal/models"
	"shortener/internal/service"
)

const (
	tokenExp     = time.Hour * 720
	unauthorized = "Access denied"
)

// Claims represents the claims for a JWT token.
type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

// BaseAuth represents the base authentication middleware.
type BaseAuth struct {
	Service *service.Service
}

// Auth creates a new instance of the BaseAuth middleware.
func Auth(svc *service.Service) *BaseAuth {
	return &BaseAuth{Service: svc}
}

// Middleware returns an HTTP handler that checks for the presence of a JWT token in the request.
func (ba *BaseAuth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("token")
		rCtx := r.Context()
		if err != nil && errors.Is(err, http.ErrNoCookie) {
			newToken, err := ba.Service.BuildJWTString()
			if err != nil {
				ba.Service.Log.Err("failed build JWTString: ", err)
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Authorization", newToken)
			http.SetCookie(w, &http.Cookie{Name: "token", Value: newToken})
			token = &http.Cookie{Name: "token", Value: newToken}
		} else if err != nil {
			ba.Service.Log.Err("failed get cookie: ", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		userID := ba.Service.GetUserID(token.Value, ba.Service.SecretKey, ba.Service.Log)
		if userID == "" {
			ba.Service.Log.Err(unauthorized, "no userID")
			http.Error(w, unauthorized, http.StatusUnauthorized)
			return
		}
		newCtx := context.WithValue(rCtx, models.CtxUserIDKey, userID)
		rWithCtx := r.WithContext(newCtx)
		next.ServeHTTP(w, rWithCtx)
	})
}

// BaseCheck represents the base check middleware.
type BaseCheck struct {
	Log *logger.Log
}

// CheckAuth creates a new instance of the BaseCheck middleware.
func CheckAuth(log *logger.Log) *BaseCheck {
	return &BaseCheck{Log: log}
}

// Middleware returns an HTTP handler that checks for the presence of an Authorization header and a token cookie.
func (bc *BaseCheck) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, unauthorized, http.StatusUnauthorized)
			return
		}
		token, err := r.Cookie("token")
		if err != nil {
			bc.Log.Err("failed get token from cookies: ", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		if auth != token.Value {
			http.Error(w, unauthorized, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
