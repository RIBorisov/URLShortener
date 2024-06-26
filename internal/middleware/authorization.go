package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"shortener/internal/logger"
	"shortener/internal/models"
	"shortener/internal/service"
)

const (
	tokenExp     = time.Hour * 720
	unauthorized = "Access denied"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

type BaseAuth struct {
	Service *service.Service
	User    *models.User
}

func Auth(svc *service.Service, user *models.User) *BaseAuth {
	return &BaseAuth{Service: svc, User: user}
}

func (ba *BaseAuth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("token")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				newToken, err := buildJWTString(ba.Service.SecretKey)
				if err != nil {
					ba.Service.Log.Err("failed build JWTString: ", err)
					http.Error(w, "", http.StatusInternalServerError)
					return
				}
				w.Header().Set("Authorization", newToken)
				http.SetCookie(w, &http.Cookie{Name: "token", Value: newToken})

				userID := getUserID(newToken, ba.Service.SecretKey, ba.Service.Log)
				if userID == "" {
					ba.Service.Log.Err(unauthorized, "no userID")
					http.Error(w, unauthorized, http.StatusUnauthorized)
					return
				}
				ba.User.ID = userID
				next.ServeHTTP(w, r)
				return
			}
			ba.Service.Log.Err("failed get cookie: ", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		userID := getUserID(token.Value, ba.Service.SecretKey, ba.Service.Log)
		if userID == "" {
			ba.Service.Log.Err(unauthorized, "no userID")
			http.Error(w, unauthorized, http.StatusUnauthorized)
			return
		}
		ba.User.ID = userID
		next.ServeHTTP(w, r)
	})
}

type BaseCheck struct {
	Log  *logger.Log
	User *models.User
}

func CheckAuth(log *logger.Log, user *models.User) *BaseCheck {
	return &BaseCheck{Log: log, User: user}
}
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

func buildJWTString(secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		UserID: uuid.NewString(),
	})
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to create token string: %w", err)
	}
	return tokenString, nil
}

func getUserID(tokenString, secretKey string, log *logger.Log) string {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		log.Err("failed parse with claims tokenString: ", err)
		return ""
	}
	if !token.Valid {
		log.Err("Token is not valid: ", token)
		return ""
	}

	return claims.UserID
}
