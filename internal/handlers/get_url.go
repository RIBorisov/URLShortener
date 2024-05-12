package handlers

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"net/url"
	"shortener/internal/config"
)

func GetHandler(db urlStorage, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.routes.getUrl.GetHandler"

		short := chi.URLParam(r, "id")
		long, ok := db.Get(short)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
		}
		origin, err := url.JoinPath(cfg.Server.BaseURL, short)

		if origin == "" {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, long, http.StatusTemporaryRedirect)
		w.Header().Set("Location", origin)
		_, err = w.Write([]byte(long))
		if err != nil {
			log.Printf("%s: %+v", op, err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
