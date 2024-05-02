package handlers

import "net/http"

func RedirectToURL(w http.ResponseWriter, r *http.Request, redirectTo string, originalURL string) {
	http.Redirect(w, r, redirectTo, http.StatusTemporaryRedirect)
	w.Header().Set("Location", originalURL)
}
