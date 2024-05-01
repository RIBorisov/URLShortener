package handlers

import "net/http"

func redirectToURL(w http.ResponseWriter, r *http.Request, redirectTo string, originalURL string) {
	w.Header().Set("Location", originalURL)
	http.Redirect(w, r, redirectTo, http.StatusTemporaryRedirect)
}
