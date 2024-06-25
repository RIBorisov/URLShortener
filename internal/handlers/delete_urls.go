package handlers

import (
	"fmt"
	"net/http"
	"shortener/internal/models"
	"shortener/internal/service"
)

func DeleteURLsHandler(svc *service.Service, user *models.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		fmt.Println(ctx)
		return
	}
}

/*
Добавьте в сервис новый хендлер DELETE /api/user/urls,
который в теле запроса принимает список идентификаторов сокращённых URL
для асинхронного удаления. Запрос может быть таким:
```
DELETE http://localhost:8080/api/user/urls
Content-Type: application/json

["6qxTVvsy", "RTfd56hn", "Jlfd67ds"]
```

Успешно удалить URL может пользователь, его создавший.
При запросе удалённого URL с помощью хендлера GET /{id} нужно вернуть статус 410 Gone.
*/
