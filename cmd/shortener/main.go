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

var urls = map[string]string{}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getOriginalUrl(w, r)
	case http.MethodPost:
		shortUrl(w, r)
	default:
		errorMethodHandler(w)
	}
}

func errorMethodHandler(w http.ResponseWriter) {
	http.Error(w, "Only GET or POST methods available", http.StatusMethodNotAllowed)
}

func shortUrl(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Url value should be sent as 'text/plain'", http.StatusBadRequest)
	}

	longUrl, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error when reading body value", http.StatusBadRequest)
	}
	shortUrl := "EwHXdJfB"
	urls[shortUrl] = string(longUrl)

	setHeaders(w)
	responseValue := generateUrl(r, shortUrl)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(responseValue))
}

func getOriginalUrl(w http.ResponseWriter, r *http.Request) {
	shortUrl := strings.TrimPrefix(r.URL.Path, "/")
	setHeaders(w)

	longUrl := urls[shortUrl]
	if longUrl == "" {
		w.WriteHeader(http.StatusBadRequest)
	}
	originalUrl := generateUrl(r, shortUrl)
	redirectToUrl(w, r, longUrl, originalUrl)

	w.Write([]byte(longUrl))
}

func setHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
}

func redirectToUrl(w http.ResponseWriter, r *http.Request, redirectTo string, originalUrl string) {
	w.Header().Set("Location", originalUrl)
	http.Redirect(w, r, redirectTo, http.StatusTemporaryRedirect)
}

func generateUrl(r *http.Request, shortUrl string) string {
	return "http://" + r.Host + "/" + shortUrl
}
