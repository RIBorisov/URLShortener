package handlers

import (
	"encoding/json"
	"net/http"

	"shortener/internal/logger"
	"shortener/internal/models"
	"shortener/internal/service"
)

func BatchHandler(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req []models.BatchRequest
		// обрабатываем вход
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			logger.Err("failed to decode request body", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer func() {
			if err := r.Body.Close(); err != nil {
				logger.Err("failed to close request body", err)
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
		}()

		saved, err := svc.SaveURLs(svc.Ctx, req)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		enc := json.NewEncoder(w)
		err = enc.Encode(saved)
		if err != nil {
			logger.Err("failed to encode response", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}

/*
Все записи о коротких URL сохраняйте в базе данных.
Не забудьте добавить реализацию для сохранения в файл и в память.
Стоит помнить, что:
	нужно соблюдать обратную совместимость;
	отправлять пустые батчи не нужно;
	вы умеете сжимать контент по алгоритму gzip;
	изменение в базе можно выполнять в рамках одной транзакции или одного запроса;
	необходимо избегать формирования условий для возникновения состояния гонки (race condition).
*/
