package handlers

import (
	"net/http"
)

func GenerateURL(r *http.Request, shortURL string) string {
	//result = fmt.Sprintf("%s") oute
	//r.HandleFunc( "/update/{widgetType}/{name}/{value}", widgetController.Update, ).Methods(http.MethodPost)
	return "http://" + r.Host + "/" + shortURL
}
