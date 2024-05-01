package handlers

import "net/http"

func GenerateURL(r *http.Request, shortURL string) string {
	return "http://" + r.Host + "/" + shortURL
}
