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
		log.Fatal(err)
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
	defer r.Body.Close()
	// if r.Header.Get("Content-Type") != "text/plain" {
	// 	http.Error(w, "Url value should be sent as 'text/plain'", http.StatusBadRequest)
	// }

	longURL, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error when reading body value", http.StatusBadRequest)
	}
	shortURL := "EwHXdJfB"
	URLMap[shortURL] = string(longURL)

	setHeaders(w)
	responseValue := generateURL(r, shortURL)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(responseValue))
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

	w.Write([]byte(longURL))
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
