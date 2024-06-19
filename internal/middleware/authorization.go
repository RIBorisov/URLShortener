package middleware

import (
	"net/http"
	"shortener/internal/logger"
)

type BaseAuth struct {
	Log *logger.Log
}

func Auth(log *logger.Log) *BaseAuth {
	return &BaseAuth{Log: log}
}

func (ba *BaseAuth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//next.ServeHTTP(ow, r)
	})
}
