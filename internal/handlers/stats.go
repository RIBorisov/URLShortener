package handlers

import (
	"encoding/json"
	"net/http"

	"shortener/internal/service"
)

// StatsHandler returns users and urls counter.
func StatsHandler(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		realIP := r.Header.Get("X-Real-IP")
		if svc.IsSubnetTrusted(realIP) {
			http.Error(w, "Untrusted subnet", http.StatusForbidden)
			return
		}

		stats, err := svc.GetStats(ctx)
		if err != nil {
			svc.Log.Err("failed to get stats", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")

		if err = json.NewEncoder(w).Encode(stats); err != nil {
			svc.Log.Err("failed to encode response: ", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
