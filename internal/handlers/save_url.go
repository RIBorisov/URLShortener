package handlers

import (
	"io"
	"log"
	"net/http"
	"net/url"

	"shortener/internal/config"
	"shortener/internal/service"
)

type urlStorage interface {
	Get(shortLink string) (string, bool)
	Save(shortLink, longLink string)
}

func SaveHandler(db urlStorage, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.routes.saveUrl.SaveHandler"
		long, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("%s: %+v", op, err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		s := service.Service{DB: db, BaseURL: cfg.Server.BaseURL}
		short := s.GenerateUniqueShortLink(cfg.URL.Length)
		db.Save(short, string(long))

		resultURL, err := url.JoinPath(cfg.Server.BaseURL, short)
		if err != nil {
			log.Printf("%s: %+v", op, err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(resultURL))
		if err != nil {
			log.Printf("%s: %+v", op, err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
