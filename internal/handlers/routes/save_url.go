package routes

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"shortener/internal/handlers"
)

type urlStorage interface {
	Get(shortLink string) (string, bool)
	Save(shortLink, longLink string)
}

func SaveHandler(db urlStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const (
			op = "handlers.routes.saveUrl.SaveHandler"
			// TODO(SSH): можно удалить
			errMsg = "Internal server error" // TODO: насколько корректно таким образом константу использовать?
		)
		long, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError) // TODO: Тут.
			log.Printf("%s: %+v", op, err)
		}
		short := handlers.GenerateUniqueShortLink(db)
		db.Save(short, string(long))
		// TODO(SSH)
		url.JoinPath(cfg.BasePath, short)
		// resp := handlers.GetOriginalURL(short)
		if resp == "" {
			http.Error(w, "", http.StatusInternalServerError) // TODO: Тут.
			return
		}
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(resp))
		if err != nil {
			log.Printf("%s: %+v", op, err)
			http.Error(w, "", http.StatusInternalServerError) // TODO: Тут.
			return
		}
	}
}
