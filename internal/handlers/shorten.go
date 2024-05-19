package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"shortener/internal/models"
	"shortener/internal/service"
)

func ShortenHandler(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.Request
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&req); err != nil {
			slog.Error("failed to decode request body", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		short := svc.SaveURL(req.URL)
		resultURL, err := url.JoinPath(svc.BaseURL, short)
		if err != nil {
			slog.Error("failed to join path to get result URL", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		resp := models.Response{
			Result: resultURL,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		enc := json.NewEncoder(w)
		if err = enc.Encode(resp); err != nil {
			slog.Error("failed to encode response", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
