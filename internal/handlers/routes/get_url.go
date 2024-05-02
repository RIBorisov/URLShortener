package routes

import (
	"fmt"
	"net/http"
	"shortener/internal/handlers"
	"shortener/internal/storage"
	"strings"
)

func GetURLHandler(w http.ResponseWriter, r *http.Request) {
	shortLink := strings.TrimPrefix(r.URL.Path, "/")
	handlers.SetHeadersHandler(w)

	mapper := storage.Mapper
	fmt.Println("Mapper objects: ", mapper.Count())
	longLink, ok := mapper.Get(shortLink)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
	}

	originalURL := handlers.GenerateURL(r, shortLink)
	handlers.RedirectToURL(w, r, longLink, originalURL)

	_, err := w.Write([]byte(longLink))
	if err != nil {
		return
	}
}
