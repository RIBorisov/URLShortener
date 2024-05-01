package main

import (
	"io"
	"log"
	"net/http"
	"strings"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", mainHandler)

	log.Println("Server started on port 8080")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}

var URLMap = map[string]string{}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getOriginalURL(w, r)
	case http.MethodPost:
		shortURL(w, r)
	default:
		errorMethodHandler(w)
	}
}

func errorMethodHandler(w http.ResponseWriter) {
	http.Error(w, "Only GET or POST methods available", http.StatusMethodNotAllowed)
}

func shortURL(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(r.Body)

	longURL, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error when reading body value", http.StatusBadRequest)
	}
	shortURL := "EwHXdJfB"
	URLMap[shortURL] = string(longURL)

	setHeaders(w)
	responseValue := generateURL(r, shortURL)
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(responseValue))
	if err != nil {
		return
	}
}

func getOriginalURL(w http.ResponseWriter, r *http.Request) {
	shortURL := strings.TrimPrefix(r.URL.Path, "/")
	setHeaders(w)

	longURL := URLMap[shortURL]
	if longURL == "" {
		w.WriteHeader(http.StatusBadRequest)
	}
	originalURL := generateURL(r, shortURL)
	redirectToURL(w, r, longURL, originalURL)

	_, err := w.Write([]byte(longURL))
	if err != nil {
		return
	}
}

func setHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
}

func redirectToURL(w http.ResponseWriter, r *http.Request, redirectTo string, originalURL string) {
	w.Header().Set("Location", originalURL)
	http.Redirect(w, r, redirectTo, http.StatusTemporaryRedirect)
}

func generateURL(r *http.Request, shortURL string) string {
	return "http://" + r.Host + "/" + shortURL
}
