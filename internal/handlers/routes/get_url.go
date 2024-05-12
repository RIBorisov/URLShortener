package routes

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"shortener/internal/handlers"
)

type urlGetter interface {
	Get(shortLink string) (string, bool)
}

func GetHandler(db urlGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const (
			op     = "handlers.routes.getUrl.GetHandler"
			errMsg = "Internal server error"
		)
		short := chi.URLParam(r, "id")
		long, ok := db.Get(short)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
		}
		origin := handlers.GetOriginalURL(short)
		if origin == "" {
			http.Error(w, errMsg, http.StatusInternalServerError)
			return
		}
		handlers.RedirectToURL(w, r, long, origin)
		_, err := w.Write([]byte(long))
		if err != nil {
			log.Printf("%s: %+v", op, err)
			http.Error(w, errMsg, http.StatusInternalServerError)
			return
		}
	}
}
