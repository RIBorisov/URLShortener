package routes

import (
	"io"
	"log"
	"net/http"
	"shortener/internal/handlers"
)

type urlStorage interface {
	Get(shortLink string) (string, bool)
	Save(shortLink, longLink string)
}

func SaveHandler(db urlStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const (
			op     = "handlers.routes.saveUrl.SaveHandler"
			errMsg = "Internal server error"
		)
		long, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, errMsg, http.StatusInternalServerError)
			log.Printf("%s: %+v", op, err)
		}
		short := handlers.GenerateUniqueShortLink(db)
		db.Save(short, string(long))
		resp := handlers.GetOriginalURL(short)
		if resp == "" {
			http.Error(w, errMsg, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(resp))
		if err != nil {
			log.Printf("%s: %+v", op, err)
			http.Error(w, errMsg, http.StatusInternalServerError)
			return
		}
	}
}
