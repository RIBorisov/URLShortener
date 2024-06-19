package handlers

import (
	"fmt"
	"net/http"
	"shortener/internal/models"
	"shortener/internal/service"
)

func GetURLsHandler(svc *service.Service, user *models.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		urls, err := svc.GetUserURLs(ctx, user)
		if err != nil {
			svc.Log.Err("failed get user urls: ", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		fmt.Println(urls)
		fmt.Println(err)
	}
}
