package handlers

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"shortener/internal/models"

	"shortener/internal/service"
	"shortener/internal/storage"
)

func SaveHandler(svc *service.Service, user *models.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		long, err := io.ReadAll(r.Body)
		if err != nil {
			svc.Log.Err("failed to read body: ", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		short, err := svc.SaveURL(ctx, string(long), user)
		if err != nil {
			var duplicateErr *storage.DuplicateRecordError
			if errors.As(err, &duplicateErr) {
				w.WriteHeader(http.StatusConflict)
				short = duplicateErr.Message
			} else {
				svc.Log.Err("failed to save url: ", err)
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
		} else {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusCreated)
		}

		resultURL, err := url.JoinPath(svc.BaseURL, short)
		if err != nil {
			svc.Log.Err("failed to join path to get result URL: ", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		_, err = w.Write([]byte(resultURL))
		if err != nil {
			svc.Log.Err("failed to write the full URL response to client: ", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
